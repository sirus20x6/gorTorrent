// handlers/torrent.go
package handlers

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func (h *Handler) HandleTorrentList(w http.ResponseWriter, r *http.Request) {
    // Add HTMX specific headers
    w.Header().Set("HX-Trigger", "updateTorrentList")
    
    torrents, err := h.torrentService.GetTorrents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Only return the inner content for HTMX requests
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "torrent-list.html", torrents)
        return
    }

    // Return full page for regular requests
    tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
        "Page": "torrents",
        "Torrents": torrents,
    })
}
