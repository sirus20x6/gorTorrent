// middleware/cachedresponse/cached.go
package cachedresponse

import (
    "bytes"
    "compress/gzip"
    "crypto/sha256"
    "fmt"
    "hash/crc32"
    "net/http"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

// Config holds configuration for cached responses
type Config struct {
    UseGzip         bool
    GzipLevel       int
    MinGzipSize     int64
    GzipCommand     string
    DefaultHeaders  map[string]string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
    return Config{
        UseGzip:     true,
        GzipLevel:   6,
        MinGzipSize: 2048,
        GzipCommand: "gzip",
        DefaultHeaders: map[string]string{
            "Cache-Control": "no-cache",
            "Pragma":       "no-cache",
        },
    }
}

// CachedResponse handles cached content serving
type CachedResponse struct {
    config Config
}

// New creates a new CachedResponse instance
func New(config Config) *CachedResponse {
    return &CachedResponse{
        config: config,
    }
}

// Send sends content with proper caching headers
func (c *CachedResponse) Send(w http.ResponseWriter, content []byte, contentType string, cacheable bool) error {
    // Set server timestamp
    w.Header().Set("X-Server-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))

    // Handle caching
    if cacheable && w.Header().Get("Cache-Control") == "" {
        etag := fmt.Sprintf("\"%X\"", crc32.ChecksumIEEE(content))
        w.Header().Set("ETag", etag)

        if match := w.Header().Get("If-None-Match"); match == etag {
            w.WriteHeader(http.StatusNotModified)
            return nil
        }
    }

    // Set content type
    if contentType != "" {
        w.Header().Set("Content-Type", contentType+"; charset=UTF-8")
    }

    // Handle compression
    contentLength := len(content)
    if c.shouldCompress(r, contentLength) {
        compressed, err := c.compress(content)
        if err != nil {
            return fmt.Errorf("compression failed: %w", err)
        }
        content = compressed
        w.Header().Set("Content-Encoding", "gzip")
    }

    // Set content length
    w.Header().Set("Content-Length", strconv.Itoa(len(content)))

    // Write content
    if _, err := w.Write(content); err != nil {
        return fmt.Errorf("failed to write response: %w", err)
    }

    return nil
}

// shouldCompress determines if content should be compressed
func (c *CachedResponse) shouldCompress(r *http.Request, contentLength int) bool {
    if !c.config.UseGzip || contentLength < int(c.config.MinGzipSize) {
        return false
    }

    acceptEncoding := r.Header.Get("Accept-Encoding")
    return strings.Contains(acceptEncoding, "gzip")
}

// compress compresses content using gzip
func (c *CachedResponse) compress(content []byte) ([]byte, error) {
    // Try native gzip first
    var buf bytes.Buffer
    gz, err := gzip.NewWriterLevel(&buf, c.config.GzipLevel)
    if err != nil {
        return nil, err
    }

    if _, err := gz.Write(content); err != nil {
        return nil, err
    }

    if err := gz.Close(); err != nil {
        return nil, err
    }

    // If native compression succeeded, return result
    if buf.Len() > 0 {
        return buf.Bytes(), nil
    }

    // Fallback to external gzip command
    cmd := exec.Command(c.config.GzipCommand, 
        fmt.Sprintf("-%d", c.config.GzipLevel))
    cmd.Stdin = bytes.NewReader(content)

    compressed, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("external gzip failed: %w", err)
    }

    return compressed, nil
}

// Middleware returns middleware that adds default cache headers
func (c *CachedResponse) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add default headers
        for key, value := range c.config.DefaultHeaders {
            w.Header().Set(key, value)
        }
        next.ServeHTTP(w, r)
    })
}

// ResponseWriter wraps http.ResponseWriter to capture content
type ResponseWriter struct {
    http.ResponseWriter
    content    bytes.Buffer
    statusCode int
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
    return w.content.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}

// CacheMiddleware creates middleware that caches responses
func (c *CachedResponse) CacheMiddleware(duration time.Duration) func(http.Handler) http.Handler {
    cache := make(map[string]cachedItem)
    var mu sync.RWMutex

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only cache GET requests
            if r.Method != http.MethodGet {
                next.ServeHTTP(w, r)
                return
            }

            // Generate cache key
            key := fmt.Sprintf("%s-%s", r.URL.Path, r.URL.RawQuery)
            mu.RLock()
            if item, ok := cache[key]; ok && !item.isExpired() {
                mu.RUnlock()
                c.Send(w, item.content, item.contentType, true)
                return
            }
            mu.RUnlock()

            // Capture response
            rw := &ResponseWriter{ResponseWriter: w}
            next.ServeHTTP(rw, r)

            // Cache response if successful
            if rw.statusCode == 0 || rw.statusCode == http.StatusOK {
                mu.Lock()
                cache[key] = cachedItem{
                    content:     rw.content.Bytes(),
                    contentType: w.Header().Get("Content-Type"),
                    expiry:     time.Now().Add(duration),
                }
                mu.Unlock()
            }

            // Send response
            c.Send(w, rw.content.Bytes(), w.Header().Get("Content-Type"), true)
        })
    }
}

type cachedItem struct {
    content     []byte
    contentType string
    expiry      time.Time
}

func (i cachedItem) isExpired() bool {
    return time.Now().After(i.expiry)
}