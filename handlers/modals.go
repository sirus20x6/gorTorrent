// handlers/modals.go
package handlers

import (
   "net/http"
   "path/filepath"
   "github.com/go-chi/chi/v5"
)

// Modal handler routes
func (h *Handler) ModalRoutes() chi.Router {
   r := chi.NewRouter()
   
   r.Get("/add-torrent", h.handleAddTorrentModal)
   r.Get("/torrent-details/{hash}", h.handleTorrentDetailsModal)
   r.Get("/settings", h.handleSettingsModal)
   r.Get("/close-modal", h.handleCloseModal)
   
   return r
}

// Handle add torrent modal display
func (h *Handler) handleAddTorrentModal(w http.ResponseWriter, r *http.Request) {
   // Get available download directories
   dirs, err := h.torrentSvc.GetDownloadDirectories()
   if err != nil {
       h.handleError(w, err)
       return
   }

   // Get label list if available
   labels, err := h.torrentSvc.GetLabels()
   if err != nil {
       labels = []string{} // Default to empty if not available
   }

   data := map[string]interface{}{
       "Directories": dirs,
       "Labels":      labels,
       "DefaultDir":  h.torrentSvc.GetDefaultDirectory(),
   }

   h.renderPartial(w, "modals/add-torrent.html", data)
}

// Handle torrent details modal
func (h *Handler) handleTorrentDetailsModal(w http.ResponseWriter, r *http.Request) {
   hash := chi.URLParam(r, "hash")
   
   // Get torrent details
   details, err := h.torrentSvc.GetTorrentDetails(hash)
   if err != nil {
       h.handleError(w, err)
       return
   }

   // Get additional information in parallel
   type result struct {
       files    []services.TorrentFile
       peers    []services.Peer
       trackers []services.Tracker
       err      error
   }
   
   ch := make(chan result)
   go func() {
       files, err := h.torrentSvc.GetTorrentFiles(hash)
       ch <- result{files: files, err: err}
   }()

   go func() {
       peers, err := h.torrentSvc.GetPeers(hash)
       ch <- result{peers: peers, err: err}
   }()

   go func() {
       trackers, err := h.torrentSvc.GetTrackers(hash)
       ch <- result{trackers: trackers, err: err}
   }()

   // Collect results
   var files []services.TorrentFile
   var peers []services.Peer
   var trackers []services.Tracker
   
   for i := 0; i < 3; i++ {
       res := <-ch
       if res.err != nil {
           h.handleError(w, res.err)
           return
       }
       if res.files != nil {
           files = res.files
       }
       if res.peers != nil {
           peers = res.peers
       }
       if res.trackers != nil {
           trackers = res.trackers
       }
   }

   data := map[string]interface{}{
       "Torrent":  details,
       "Files":    files,
       "Peers":    peers,
       "Trackers": trackers,
   }

   h.renderPartial(w, "modals/torrent-details.html", data)
}

// Handle settings modal
func (h *Handler) handleSettingsModal(w http.ResponseWriter, r *http.Request) {
   // Get current settings
   settings, err := h.torrentSvc.GetSettings()
   if err != nil {
       h.handleError(w, err)
       return
   }

   data := map[string]interface{}{
       "Settings": settings,
   }

   h.renderPartial(w, "modals/settings.html", data)
}

// Handle modal close
func (h *Handler) handleCloseModal(w http.ResponseWriter, r *http.Request) {
   // For HTMX, return empty response to close modal
   w.Header().Set("HX-Trigger", "closeModal")
   w.WriteHeader(http.StatusOK)
}

// Handle file browser modal
func (h *Handler) handleFileBrowserModal(w http.ResponseWriter, r *http.Request) {
   path := r.URL.Query().Get("path")
   if path == "" {
       path = h.torrentSvc.GetDefaultDirectory()
   }

   // Get directory contents
   entries, err := h.torrentSvc.ListDirectory(path)
   if err != nil {
       h.handleError(w, err)
       return
   }

   data := map[string]interface{}{
       "CurrentPath": path,
       "ParentPath":  filepath.Dir(path),
       "Entries":     entries,
   }

   h.renderPartial(w, "modals/file-browser.html", data)
}

// Helper struct for directory entries
type DirectoryEntry struct {
   Name      string
   Path      string
   IsDir     bool
   Size      int64
   Modified  time.Time
}

// Helper struct for settings
type Settings struct {
   DownloadDir      string
   DownloadRate     int
   UploadRate       int
   MaxConnections   int
   MaxUploads       int
   DHT              bool
   PEX              bool
   LSD              bool
   UTorrentPeerID   bool
}

// Helper types for service layer
type (
   TorrentFile struct {
       Path     string
       Size     int64
       Progress float64
       Priority int
   }

   Peer struct {
       Address     string
       Client      string
       Flags       string
       Progress    float64
       DownSpeed   int64
       UpSpeed    int64
   }

   Tracker struct {
       URL          string
       Status       string
       Seeds        int
       Peers        int
       LastUpdated  time.Time
   }
)