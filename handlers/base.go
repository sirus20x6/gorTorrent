// handlers/base.go
package handlers

import (
    "fmt"
    "net/http"
    "html/template"
    "path/filepath"

    "ruTorrent-web/services"
)

// Handler holds dependencies for all handlers
type Handler struct {
    templates     *template.Template
    torrentSvc    *services.TorrentService
    templateCache map[string]*template.Template
}

// Config holds all configuration for the handlers
type Config struct {
    TemplatesDir  string
    RTorrentURL   string
}

// New creates a new handler instance
func New(cfg Config) (*Handler, error) {
    // Initialize template cache
    templates, err := parseTemplates(cfg.TemplatesDir)
    if err != nil {
        return nil, fmt.Errorf("failed to parse templates: %w", err)
    }

    // Initialize services
    torrentSvc := services.NewTorrentService(cfg.RTorrentURL)

    return &Handler{
        templates:  templates,
        torrentSvc: torrentSvc,
        templateCache: make(map[string]*template.Template),
    }, nil
}

// parseTemplates loads and parses all templates
func parseTemplates(templatesDir string) (*template.Template, error) {
    // Define template functions
    funcMap := template.FuncMap{
        "formatBytes": formatBytes,
        "formatSpeed": formatSpeed,
        "formatDate": formatDate,
    }

    // Parse all templates
    pattern := filepath.Join(templatesDir, "**", "*.html")
    tmpl, err := template.New("").Funcs(funcMap).ParseGlob(pattern)
    if err != nil {
        return nil, err
    }

    return tmpl, nil
}

// render handles template rendering with error handling
func (h *Handler) render(w http.ResponseWriter, name string, data interface{}) {
    // Set default content type
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    // Check if template exists
    tmpl := h.templates.Lookup(name)
    if tmpl == nil {
        http.Error(w, fmt.Sprintf("Template %s not found", name), http.StatusInternalServerError)
        return
    }

    // Execute template
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// renderPartial handles HTMX partial updates
func (h *Handler) renderPartial(w http.ResponseWriter, name string, data interface{}) {
    // For HTMX requests, we don't need the full base template
    tmpl := h.templates.Lookup(name)
    if tmpl == nil {
        http.Error(w, fmt.Sprintf("Template %s not found", name), http.StatusInternalServerError)
        return
    }

    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// isHXRequest checks if the request is from HTMX
func (h *Handler) isHXRequest(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}

// Helper functions for templates
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatSpeed(bytesPerSec int64) string {
    return formatBytes(bytesPerSec) + "/s"
}

func formatDate(timestamp int64) string {
    return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// RequestError represents an error that occurred during request processing
type RequestError struct {
    Status  int
    Message string
}

func (e RequestError) Error() string {
    return e.Message
}

// handleError handles errors in a consistent way
func (h *Handler) handleError(w http.ResponseWriter, err error) {
    // Check if it's a known request error
    if reqErr, ok := err.(RequestError); ok {
        http.Error(w, reqErr.Message, reqErr.Status)
        return
    }

    // Unknown error - return 500
    http.Error(w, "Internal server error", http.StatusInternalServerError)
}