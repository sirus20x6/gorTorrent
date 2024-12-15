// handlers/torrents.go
package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "path/filepath"
    "strconv"
)

// Routes returns all torrent-related routes
func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()

    // List/filter torrents
    r.Get("/", h.handleListTorrents)
    r.Get("/torrents", h.handleListTorrents)
    r.Get("/torrents/filter", h.handleFilterTorrents)
    r.Get("/torrents/sort", h.handleSortTorrents)

    // Add torrents
    r.Post("/torrents/add", h.handleAddTorrent)
    r.Post("/torrents/add-magnet", h.handleAddMagnet)

    // Torrent actions
    r.Post("/torrents/{hash}/start", h.handleStartTorrent)
    r.Post("/torrents/{hash}/stop", h.handleStopTorrent)
    r.Post("/torrents/{hash}/pause", h.handlePauseTorrent)
    r.Delete("/torrents/{hash}", h.handleDeleteTorrent)

    // Torrent details
    r.Get("/torrents/{hash}", h.handleTorrentDetails)
    r.Get("/torrents/{hash}/files", h.handleTorrentFiles)
    r.Get("/torrents/{hash}/peers", h.handleTorrentPeers)
    r.Get("/torrents/{hash}/trackers", h.handleTorrentTrackers)

    return r
}

// List torrents with optional filtering and sorting
func (h *Handler) handleListTorrents(w http.ResponseWriter, r *http.Request) {
    // Get filter parameters
    filter := r.URL.Query().Get("filter")
    search := r.URL.Query().Get("search")
    sort := r.URL.Query().Get("sort")
    
    // Get torrents from service
    torrents, err := h.torrentSvc.GetTorrents()
    if err != nil {
        h.handleError(w, err)
        return
    }

    // Apply filters if present
    if filter != "" {
        torrents = filterTorrents(torrents, filter)
    }
    if search != "" {
        torrents = searchTorrents(torrents, search)
    }
    if sort != "" {
        sortTorrents(torrents, sort)
    }

    // Prepare template data
    data := map[string]interface{}{
        "Torrents": torrents,
        "Filter":   filter,
        "Search":   search,
        "Sort":     sort,
    }

    // Check if it's an HTMX request
    if h.isHXRequest(r) {
        h.renderPartial(w, "components/torrent-list.html", data)
        return
    }

    // Regular request - render full page
    h.render(w, "base.html", data)
}

// Add a new torrent file
func (h *Handler) handleAddTorrent(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        h.handleError(w, RequestError{
            Status:  http.StatusBadRequest,
            Message: "Failed to parse form: " + err.Error(),
        })
        return
    }

    // Get form values
    startImmediately := r.FormValue("start_immediately") == "on"
    skipChecking := r.FormValue("skip_checking") == "on"
    directory := r.FormValue("directory")

    // Get uploaded file
    file, header, err := r.FormFile("torrent")
    if err != nil {
        h.handleError(w, RequestError{
            Status:  http.StatusBadRequest,
            Message: "No torrent file provided",
        })
        return
    }
    defer file.Close()

    // Read file content
    content, err := io.ReadAll(file)
    if err != nil {
        h.handleError(w, err)
        return
    }

    // Add torrent via service
    err = h.torrentSvc.AddTorrent(content, startImmediately, directory)
    if err != nil {
        h.handleError(w, err)
        return
    }

    // For HTMX requests, return updated list
    if h.isHXRequest(r) {
        h.handleListTorrents(w, r)
        return
    }

    // Regular form submission - redirect to list
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Add a magnet link
func (h *Handler) handleAddMagnet(w http.ResponseWriter, r *http.Request) {
    magnetURL := r.FormValue("magnet")
    if magnetURL == "" {
        h.handleError(w, RequestError{
            Status:  http.StatusBadRequest,
            Message: "No magnet URL provided",
        })
        return
    }

    startImmediately := r.FormValue("start_immediately") == "on"
    directory := r.FormValue("directory")

    err := h.torrentSvc.AddMagnet(magnetURL, startImmediately, directory)
    if err != nil {
        h.handleError(w, err)
        return
    }

    if h.isHXRequest(r) {
        h.handleListTorrents(w, r)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Start a torrent
func (h *Handler) handleStartTorrent(w http.ResponseWriter, r *http.Request) {
    hash := chi.URLParam(r, "hash")
    err := h.torrentSvc.StartTorrent(hash)
    if err != nil {
        h.handleError(w, err)
        return
    }

    if h.isHXRequest(r) {
        h.handleListTorrents(w, r)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Stop a torrent
func (h *Handler) handleStopTorrent(w http.ResponseWriter, r *http.Request) {
    hash := chi.URLParam(r, "hash")
    err := h.torrentSvc.StopTorrent(hash)
    if err != nil {
        h.handleError(w, err)
        return
    }

    if h.isHXRequest(r) {
        h.handleListTorrents(w, r)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete a torrent
func (h *Handler) handleDeleteTorrent(w http.ResponseWriter, r *http.Request) {
    hash := chi.URLParam(r, "hash")
    deleteFiles := r.URL.Query().Get("delete_files") == "true"

    err := h.torrentSvc.DeleteTorrent(hash, deleteFiles)
    if err != nil {
        h.handleError(w, err)
        return
    }

    if h.isHXRequest(r) {
        h.handleListTorrents(w, r)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Get torrent details
func (h *Handler) handleTorrentDetails(w http.ResponseWriter, r *http.Request) {
    hash := chi.URLParam(r, "hash")
    
    details, err := h.torrentSvc.GetTorrentDetails(hash)
    if err != nil {
        h.handleError(w, err)
        return
    }

    data := map[string]interface{}{
        "Torrent": details,
    }

    if h.isHXRequest(r) {
        h.renderPartial(w, "components/torrent-details.html", data)
        return
    }

    h.render(w, "torrent-details.html", data)
}

// Helper functions

func filterTorrents(torrents []services.Torrent, filter string) []services.Torrent {
    if filter == "" || filter == "all" {
        return torrents
    }

    filtered := make([]services.Torrent, 0)
    for _, t := range torrents {
        switch filter {
        case "downloading":
            if t.Status == "downloading" {
                filtered = append(filtered, t)
            }
        case "seeding":
            if t.Status == "seeding" {
                filtered = append(filtered, t)
            }
        case "completed":
            if t.Progress == 100 {
                filtered = append(filtered, t)
            }
        case "active":
            if t.DownSpeed > 0 || t.UpSpeed > 0 {
                filtered = append(filtered, t)
            }
        case "error":
            if t.Status == "error" {
                filtered = append(filtered, t)
            }
        }
    }
    return filtered
}

func searchTorrents(torrents []services.Torrent, query string) []services.Torrent {
    if query == "" {
        return torrents
    }

    query = strings.ToLower(query)
    filtered := make([]services.Torrent, 0)
    for _, t := range torrents {
        if strings.Contains(strings.ToLower(t.Name), query) {
            filtered = append(filtered, t)
        }
    }
    return filtered
}

func sortTorrents(torrents []services.Torrent, field string) {
    sort.Slice(torrents, func(i, j int) bool {
        switch field {
        case "name":
            return torrents[i].Name < torrents[j].Name
        case "size":
            return torrents[i].Size < torrents[j].Size
        case "progress":
            return torrents[i].Progress < torrents[j].Progress
        case "speed":
            return torrents[i].DownSpeed+torrents[i].UpSpeed < torrents[j].DownSpeed+torrents[j].UpSpeed
        case "added":
            return torrents[i].AddedDate.Before(torrents[j].AddedDate)
        default:
            return torrents[i].Name < torrents[j].Name
        }
    })
}