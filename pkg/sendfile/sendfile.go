// util/sendfile/sendfile.go
package sendfile

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "your-project/util/fileutil"
)

// SendFile sends a file to the client with proper headers and range support
func SendFile(w http.ResponseWriter, filename string, opts ...Option) error {
    options := defaultOptions()
    for _, opt := range opts {
        opt(options)
    }

    // Check if file exists and is readable
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed to stat file: %w", err)
    }

    // Validate file is regular and readable
    if !stat.Mode().IsRegular() {
        return errors.New("not a regular file")
    }

    // Get the filename to send
    nameToSend := options.contentName
    if nameToSend == "" {
        nameToSend = filepath.Base(filename)
    }

    // Handle user agent specific filename encoding
    if isIE := strings.Contains(strings.ToLower(options.userAgent), "msie"); isIE {
        nameToSend = url.QueryEscape(nameToSend)
    }

    // Set content type
    contentType := options.contentType
    if contentType == "" {
        contentType = "application/octet-stream"
    }
    w.Header().Set("Content-Type", contentType)

    // Set disposition header
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf(`attachment; filename="%s"`, nameToSend))

    // Handle caching headers
    etag := fmt.Sprintf(`"%x-%x-%x"`, stat.Sys().(*syscall.Stat_t).Ino, 
        stat.Size(), stat.ModTime().UnixNano())

    w.Header().Set("ETag", etag)
    w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

    // Check if-none-match
    if match := options.request.Header.Get("If-None-Match"); match != "" {
        if match == etag {
            w.WriteHeader(http.StatusNotModified)
            return nil
        }
    }

    // Check if-modified-since
    if t, err := time.Parse(http.TimeFormat, 
        options.request.Header.Get("If-Modified-Since")); err == nil {
        if stat.ModTime().Unix() <= t.Unix() {
            w.WriteHeader(http.StatusNotModified)
            return nil
        }
    }

    // Handle range requests
    rangeHeader := options.request.Header.Get("Range")
    if rangeHeader != "" {
        ranges, err := parseRange(rangeHeader, stat.Size())
        if err != nil {
            w.Header().Set("Content-Range", 
                fmt.Sprintf("bytes */%d", stat.Size()))
            return &httpError{
                StatusCode: http.StatusRequestedRangeNotSatisfiable,
                Message:    "invalid range",
            }
        }

        if len(ranges) == 1 {
            r := ranges[0]
            if _, err := file.Seek(r.start, 0); err != nil {
                return fmt.Errorf("failed to seek file: %w", err)
            }
            w.Header().Set("Content-Range", 
                fmt.Sprintf("bytes %d-%d/%d", r.start, r.end, stat.Size()))
            w.Header().Set("Content-Length", strconv.FormatInt(r.length, 10))
            w.WriteHeader(http.StatusPartialContent)
            
            if options.request.Method != http.MethodHead {
                if _, err := io.CopyN(w, file, r.length); err != nil {
                    return fmt.Errorf("failed to copy range: %w", err)
                }
            }
            return nil
        }
    }

    // Full file response
    w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
    w.Header().Set("Accept-Ranges", "bytes")

    if options.request.Method != http.MethodHead {
        if _, err := io.Copy(w, file); err != nil {
            return fmt.Errorf("failed to copy file: %w", err)
        }
    }

    return nil
}

// SendCachedImage sends an image with caching headers
func SendCachedImage(w http.ResponseWriter, location string, imageType string, duration time.Duration) error {
    w.Header().Set("Content-Type", imageType)
    w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int(duration.Seconds())))
    http.ServeFile(w, &http.Request{Method: "GET"}, location)
    return nil
}

// Option type for configuring SendFile behavior
type Option func(*options)

type options struct {
    contentType string
    contentName string
    userAgent   string
    request     *http.Request
}

func defaultOptions() *options {
    return &options{
        request: &http.Request{Method: "GET"},
    }
}

// WithContentType sets the content type
func WithContentType(contentType string) Option {
    return func(o *options) {
        o.contentType = contentType
    }
}

// WithContentName sets the filename to send
func WithContentName(name string) Option {
    return func(o *options) {
        o.contentName = name
    }
}

// WithRequest sets the original request
func WithRequest(r *http.Request) Option {
    return func(o *options) {
        o.request = r
        o.userAgent = r.UserAgent()
    }
}

// httpError represents an HTTP error response
type httpError struct {
    StatusCode int
    Message    string
}

func (e *httpError) Error() string {
    return e.Message
}

// byteRange specifies the byte range to be sent.
type byteRange struct {
    start, end int64
    length     int64
}

// parseRange parses a Range header string as per RFC 7233
func parseRange(s string, size int64) ([]byteRange, error) {
    if !strings.HasPrefix(s, "bytes=") {
        return nil, fmt.Errorf("invalid range format")
    }

    var ranges []byteRange
    noOverlap := false
    for _, ra := range strings.Split(s[6:], ",") {
        ra = strings.TrimSpace(ra)
        if ra == "" {
            continue
        }
        i := strings.Index(ra, "-")
        if i < 0 {
            return nil, fmt.Errorf("invalid range format")
        }
        start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
        r := byteRange{start: 0, end: size - 1}
        if start == "" {
            // If no start is specified, end specifies the
            // range start relative to the end of the file,
            // and we are dealing with <suffix-length>
            // which has to be a non-negative integer as per
            // RFC 7233 Section 2.1 "Byte-Ranges".
            if end == "" || end[0] == '-' {
                return nil, fmt.Errorf("invalid range format")
            }
            i, err := strconv.ParseInt(end, 10, 64)
            if i < 0 || err != nil {
                return nil, fmt.Errorf("invalid range format")
            }
            if i > size {
                i = size
            }
            r.start = size - i
            r.end = size - 1
        } else {
            i, err := strconv.ParseInt(start, 10, 64)
            if err != nil || i < 0 {
                return nil, fmt.Errorf("invalid range format")
            }
            if i >= size {
                // If the range begins after the size of the content,
                // then it does not overlap.
                noOverlap = true
                continue
            }
            r.start = i
            if end == "" {
                // If no end is specified, range extends to end of the file.
                r.end = size - 1
            } else {
                i, err := strconv.ParseInt(end, 10, 64)
                if err != nil || r.start > i {
                    return nil, fmt.Errorf("invalid range format")
                }
                if i >= size {
                    i = size - 1
                }
                r.end = i
            }
        }
        r.length = r.end - r.start + 1
        ranges = append(ranges, r)
    }
    if noOverlap && len(ranges) == 0 {
        // The specified ranges did not overlap with the content.
        return nil, fmt.Errorf("invalid range: no overlap")
    }
    return ranges, nil
}