// internal/services/torrent_service.go

package services

import (
    "your-project/internal/rtorrent"
    "sync"
    "time"
)

type TorrentService struct {
    client      *rtorrent.Client
    updateChan  chan struct{}
    torrents    map[string]*Torrent
    mu          sync.RWMutex
}

type Torrent struct {
    Hash        string
    Name        string
    Size        int64
    Downloaded  int64
    UploadRate  int64
    DownloadRate int64
    State       int
    Progress    float64
    // Add more fields as needed
}

func NewTorrentService(client *rtorrent.Client) *TorrentService {
    ts := &TorrentService{
        client:     client,
        updateChan: make(chan struct{}),
        torrents:   make(map[string]*Torrent),
    }
    
    go ts.backgroundUpdater()
    return ts
}

func (s *TorrentService) backgroundUpdater() {
    ticker := time.NewTicker(2500 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.updateTorrents()
        case <-s.updateChan:
            s.updateTorrents()
        }
    }
}

func (s *TorrentService) updateTorrents() {
    torrents, err := s.client.GetTorrents()
    if err != nil {
        // Handle error, maybe log it
        return
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    // Update torrents map
    s.torrents = torrents
}

// two versions need to be combined

// internal/services/torrent_service.go

type Speeds struct {
    Download int64
    Upload   int64
}

func (s *TorrentService) GetTotalSpeeds() (Speeds, error) {
    speeds := Speeds{}
    
    // Calculate total speeds from all torrents
    s.mu.RLock()
    for _, torrent := range s.torrents {
        speeds.Download += torrent.DownloadRate
        speeds.Upload += torrent.UploadRate
    }
    s.mu.RUnlock()
    
    return speeds, nil
}