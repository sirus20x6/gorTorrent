// handlers/speed.go

type SpeedData struct {
    Download int64 `json:"download"`
    Upload   int64 `json:"upload"`
}

func (h *WebUIHandler) GetSpeed(w http.ResponseWriter, r *http.Request) {
    // Get current speeds from rTorrent
    speeds, err := h.torrentService.GetTotalSpeeds()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return speed data as JSON with HTMX trigger
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("HX-Trigger", "speedUpdate")
    json.NewEncoder(w).Encode(SpeedData{
        Download: speeds.Download,
        Upload:   speeds.Upload,
    })
}