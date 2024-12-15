// handlers/addtorrent/addtorrent.go
package addtorrent

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "path/filepath"

    "your-project/services"
    "your-project/util/fileutil"
)

// Response represents the add torrent response
type Response struct {
    Result string `json:"result"`
    Name   string `json:"name,omitempty"`
}

// Handler handles torrent addition requests
type Handler struct {
    torrentService *services.TorrentService
    maxFileSize    int64
    tempDir       string
}

// Config holds handler configuration
type Config struct {
    MaxFileSize    int64
    TempDir       string
    TorrentService *services.TorrentService
}

// New creates a new add torrent handler
func New(config Config) *Handler {
    return &Handler{
        torrentService: config.TorrentService,
        maxFileSize:    config.MaxFileSize,
        tempDir:       config.TempDir,
    }
}

// ServeHTTP handles the add torrent request
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var responses []Response

    // Check content type for multipart form
    if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data") {
        h.handleNonMultipartRequest(w, r)
        return
    }

    // Parse multipart form
    if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
        h.sendError(w, r, "Failed to parse form", err)
        return
    }
    defer r.MultipartForm.RemoveAll()

    // Get form values
    startImmediately := r.FormValue("start_immediately") == "on"
    skipChecking := r.FormValue("skip_checking") == "on"
    label := r.FormValue("label")
    dirEdit := r.FormValue("dir_edit")

    // Process directory
    if dirEdit != "" {
        if err := h.torrentService.ValidateDirectory(dirEdit); err != nil {
            responses = append(responses, Response{
                Result: "FailedDirectory",
            })
            h.sendResponses(w, r, responses)
            return
        }
    }

    // Handle file uploads
    if files := r.MultipartForm.File["torrent_file"]; len(files) > 0 {
        responses = h.handleFileUploads(files, startImmediately, skipChecking, dirEdit, label)
    }

    // Handle magnet URLs
    if url := r.FormValue("url"); url != "" {
        if response := h.handleMagnet(url, startImmediately, dirEdit, label); response != nil {
            responses = append(responses, *response)
        }
    }

    // Send responses
    if len(responses) == 0 {
        responses = append(responses, Response{Result: "Failed"})
    }
    h.sendResponses(w, r, responses)
}

// handleFileUploads processes uploaded torrent files
func (h *Handler) handleFileUploads(files []*multipart.FileHeader, start, skip bool, dir, label string) []Response {
    var responses []Response

    for _, fileHeader := range files {
        response := Response{
            Name: fileHeader.Filename,
        }

        // Open uploaded file
        file, err := fileHeader.Open()
        if err != nil {
            response.Result = "Failed"
            responses = append(responses, response)
            continue
        }
        defer file.Close()

        // Create temporary file
        tempName := h.generateTempName(fileHeader.Filename)
        tempFile, err := os.OpenFile(tempName, os.O_WRONLY|os.O_CREATE, 0644)
        if err != nil {
            response.Result = "Failed"
            responses = append(responses, response)
            continue
        }
        defer tempFile.Close()
        defer os.Remove(tempName)

        // Copy file content
        if _, err := io.Copy(tempFile, file); err != nil {
            response.Result = "Failed"
            responses = append(responses, response)
            continue
        }
        tempFile.Close()

        // Validate torrent file
        if !h.validateTorrentFile(tempName) {
            response.Result = "FailedFile"
            responses = append(responses, response)
            continue
        }

        // Add torrent
        if err := h.torrentService.AddTorrentFile(tempName, &services.AddTorrentOptions{
            Start:     start,
            SkipCheck: skip,
            Directory: dir,
            Label:     label,
        }); err != nil {
            response.Result = "Failed"
            responses = append(responses, response)
            continue
        }

        response.Result = "Success"
        responses = append(responses, response)
    }

    return responses
}

// handleMagnet processes magnet URLs
func (h *Handler) handleMagnet(url string, start bool, dir, label string) *Response {
    response := &Response{
        Name: url,
    }

    if !strings.HasPrefix(url, "magnet:") {
        // Try downloading URL
        tempName := h.generateTempName("downloaded.torrent")
        if err := h.downloadURL(url, tempName); err != nil {
            response.Result = "FailedURL"
            return response
        }
        defer os.Remove(tempName)

        // Add downloaded torrent
        if err := h.torrentService.AddTorrentFile(tempName, &services.AddTorrentOptions{
            Start:     start,
            Directory: dir,
            Label:     label,
        }); err != nil {
            response.Result = "Failed"
            return response
        }
    } else {
        // Add magnet directly
        if err := h.torrentService.AddMagnet(url, &services.AddTorrentOptions{
            Start:     start,
            Directory: dir,
            Label:     label,
        }); err != nil {
            response.Result = "Failed"
            return response
        }
    }

    response.Result = "Success"
    return response
}

// Helper functions

func (h *Handler) generateTempName(original string) string {
    ext := filepath.Ext(original)
    if ext != ".torrent" {
        ext = ".torrent"
    }
    return fileutil.GetUniqueFilename(filepath.Join(h.tempDir, fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)))
}

func (h *Handler) validateTorrentFile(path string) bool {
    // Add torrent file validation
    return true // TODO: Implement actual validation
}

func (h *Handler) downloadURL(url, dest string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("bad status: %d", resp.StatusCode)
    }

    out, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}

func (h *Handler) sendResponses(w http.ResponseWriter, r *http.Request, responses []Response) {
    if r.FormValue("json") != "" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(responses[0])
        return
    }

    w.Header().Set("Content-Type", "text/html")
    var js strings.Builder
    for _, resp := range responses {
        name := ""
        if resp.Name != "" {
            name = filepath.Base(resp.Name) + " - "
        }
        js.WriteString(fmt.Sprintf(
            `noty("%s%s","%s");`,
            name,
            "addTorrent"+resp.Result,
            map[string]string{
                "Success": "success",
                "Failed":  "error",
            }[resp.Result],
        ))
    }
    w.Write([]byte(js.String()))
}

func (h *Handler) sendError(w http.ResponseWriter, r *http.Request, msg string, err error) {
    if r.FormValue("json") != "" {
        http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusBadRequest)
        return
    }
    http.Error(w, msg, http.StatusBadRequest)
}

// NewHandler returns an http.Handler for registration
func NewHandler(config Config) http.Handler {
    return New(config)
}