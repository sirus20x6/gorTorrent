// handlers/settings.go
package handlers

import (
    "encoding/json"
    "io"
    "net/http"

    "your-project-path/services"
)

type SettingsHandler struct {
    settings *services.UISettings
}

func NewSettingsHandler() (*SettingsHandler, error) {
    settings, err := services.LoadOrCreate()
    if err != nil {
        return nil, err
    }

    return &SettingsHandler{
        settings: settings,
    }, nil
}

// HandleGetSettings handles GET requests for settings
func (h *SettingsHandler) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(h.settings.Get()))
}

// HandleSetSettings handles POST requests to update settings
func (h *SettingsHandler) HandleSetSettings(w http.ResponseWriter, r *http.Request) {
    // Read request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }

    // Validate JSON
    if !json.Valid(body) {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    // Update settings
    if err := h.settings.Set(string(body)); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// RegisterRoutes sets up the settings routes
func (h *SettingsHandler) RegisterRoutes(router *Router) {
    router.Get("/api/settings", h.HandleGetSettings)
    router.Post("/api/settings", h.HandleSetSettings)
}