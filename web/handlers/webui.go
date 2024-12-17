// handlers/webui.go

package handlers

import (
    "net/http"
    "your-project/internal/webui"
    "encoding/json"
)

type WebUIHandler struct {
    webui *webui.WebUI
}

func NewWebUIHandler(rtorrentURL string) *WebUIHandler {
    return &WebUIHandler{
        webui: webui.New(rtorrentURL),
    }
}

// MainPage handles the main page rendering
func (h *WebUIHandler) MainPage(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Version  string
        Settings map[string]interface{}
    }{
        Version:  h.webui.Version,
        Settings: h.webui.Settings,
    }
    
    render(w, "index", data)
}

// UpdateSettings handles settings updates
func (h *WebUIHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for key := range r.Form {
        value := r.Form.Get(key)
        h.webui.Settings[key] = value
    }

    // Save settings to DB/file
    if err := h.webui.SaveSettings(); err != nil {
        http.Error(w, "Failed to save settings", http.StatusInternalServerError)
        return
    }

    // Return updated settings via HTMX trigger
    w.Header().Set("HX-Trigger", `{"settingsUpdated": true}`)
}