// handlers/torrents.go
package handlers

import (
    "net/http"
    "your-project/internal/rtorrent"
)

type TorrentListHeaderData struct {
    Search        string
    CurrentFilter string
    Filters       []Filter
    SortOptions   []SortOption
    Stats         Stats
    Torrents      map[string]Torrent
}

type Filter struct {
    ID    string
    Label string
}

type SortOption struct {
    Field string
    Label string
}

type Stats struct {
    DownloadSpeed  string
    UploadSpeed    string
    ActiveTorrents int
}

type Torrent struct {
    Name         string
    Size         int64
    Downloaded   int64
    UploadRate   int64
    DownloadRate int64
    State        int
    SeedersTotal int
    PeersTotal   int
    Label        string
    Progress     float64
}

type TorrentHandler struct {
    client *rtorrent.Client
}

func NewTorrentHandler(endpoint string) *TorrentHandler {
    return &TorrentHandler{
        client: rtorrent.New(endpoint),
    }
}

func (h *TorrentHandler) HandleTorrentList(w http.ResponseWriter, r *http.Request) {
    // Get torrent list from rTorrent
    torrents, err := h.getTorrents()
    if err != nil {
        http.Error(w, "Failed to fetch torrents: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Calculate stats
    var stats Stats
    var activeTorrents int
    var totalDownRate, totalUpRate int64
    for _, t := range torrents {
        if t.DownloadRate > 0 || t.UploadRate > 0 {
            activeTorrents++
        }
        totalDownRate += t.DownloadRate
        totalUpRate += t.UploadRate
    }

    data := TorrentListHeaderData{
        Search:        r.URL.Query().Get("search"),
        CurrentFilter: r.URL.Query().Get("filter"),
        Filters: []Filter{
            {ID: "all", Label: "All"},
            {ID: "downloading", Label: "Downloading"},
            {ID: "seeding", Label: "Seeding"},
            {ID: "completed", Label: "Completed"},
            {ID: "active", Label: "Active"},
            {ID: "inactive", Label: "Inactive"},
            {ID: "error", Label: "Error"},
        },
        SortOptions: []SortOption{
            {Field: "name", Label: "Name"},
            {Field: "size", Label: "Size"},
            {Field: "progress", Label: "Progress"},
            {Field: "speed", Label: "Speed"},
            {Field: "added", Label: "Date Added"},
        },
        Stats: Stats{
            DownloadSpeed:  formatSpeed(totalDownRate),
            UploadSpeed:    formatSpeed(totalUpRate),
            ActiveTorrents: activeTorrents,
        },
        Torrents: torrents,
    }

    // Filter torrents based on search and filter
    data.Torrents = h.filterTorrents(torrents, data.Search, data.CurrentFilter)

    // Sort torrents
    sortField := r.URL.Query().Get("sort")
    if sortField != "" {
        h.sortTorrents(data.Torrents, sortField)
    }

    // If this is an HTMX request, render just the torrent list
    if r.Header.Get("HX-Request") == "true" {
        renderPartial(w, "torrents/_list", data)
        return
    }

    // Otherwise render the full page
    renderTemplate(w, "torrents/index", data)
}

func (h *TorrentHandler) getTorrents() (map[string]Torrent, error) {
    // Get list of torrent hashes
    hashes, err := h.client.GetDownloadList()
    if err != nil {
        return nil, err
    }

    torrents := make(map[string]Torrent)
    for _, hash := range hashes {
        info, err := h.client.GetTorrentInfo(hash)
        if err != nil {
            continue
        }

        // Calculate progress percentage
        progress := float64(0)
        if info["size"].(int64) > 0 {
            progress = float64(info["completed"].(int64)) / float64(info["size"].(int64)) * 100
        }

        torrents[hash] = Torrent{
            Name:         info["name"].(string),
            Size:         info["size"].(int64),
            Downloaded:   info["completed"].(int64),
            UploadRate:   info["upload_rate"].(int64),
            DownloadRate: info["download_rate"].(int64),
            State:        info["state"].(int),
            SeedersTotal: info["seeders"].(int),
            PeersTotal:   info["peers"].(int),
            Progress:     progress,
        }
    }

    return torrents, nil
}

func (h *TorrentHandler) filterTorrents(torrents map[string]Torrent, search string, filter string) map[string]Torrent {
    if search == "" && (filter == "" || filter == "all") {
        return torrents
    }

    filtered := make(map[string]Torrent)
    for hash, torrent := range torrents {
        // Apply search filter
        if search != "" && !strings.Contains(strings.ToLower(torrent.Name), strings.ToLower(search)) {
            continue
        }

        // Apply status filter
        if filter != "" && filter != "all" {
            switch filter {
            case "downloading":
                if torrent.DownloadRate == 0 {
                    continue
                }
            case "seeding":
                if torrent.Progress < 100 || torrent.UploadRate == 0 {
                    continue
                }
            case "completed":
                if torrent.Progress < 100 {
                    continue
                }
            case "active":
                if torrent.DownloadRate == 0 && torrent.UploadRate == 0 {
                    continue
                }
            case "inactive":
                if torrent.DownloadRate > 0 || torrent.UploadRate > 0 {
                    continue
                }
            case "error":
                if torrent.State&0x10 == 0 { // Assuming error bit is 0x10
                    continue
                }
            }
        }

        filtered[hash] = torrent
    }

    return filtered
}

func (h *TorrentHandler) sortTorrents(torrents map[string]Torrent, field string) {
    // Create a slice of hashes for sorting
    hashes := make([]string, 0, len(torrents))
    for hash := range torrents {
        hashes = append(hashes, hash)
    }

    // Sort the hashes based on the field
    sort.Slice(hashes, func(i, j int) bool {
        t1, t2 := torrents[hashes[i]], torrents[hashes[j]]
        switch field {
        case "name":
            return t1.Name < t2.Name
        case "size":
            return t1.Size < t2.Size
        case "progress":
            return t1.Progress < t2.Progress
        case "speed":
            return (t1.DownloadRate + t1.UploadRate) < (t2.DownloadRate + t2.UploadRate)
        default:
            return false
        }
    })

    // Create a new sorted map
    sorted := make(map[string]Torrent)
    for _, hash := range hashes {
        sorted[hash] = torrents[hash]
    }

    // Replace the original map
    torrents = sorted
}

// Helper functions
func formatSpeed(bytesPerSec int64) string {
    const unit = 1024
    if bytesPerSec < unit {
        return fmt.Sprintf("%d B/s", bytesPerSec)
    }
    div, exp := int64(unit), 0
    for n := bytesPerSec / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB/s", float64(bytesPerSec)/float64(div), "KMGTPE"[exp])
}