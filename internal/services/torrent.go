// internal/services/torrent.go
package services

import (
    "fmt"
    "time"
)

type TorrentService struct {
    client *RTorrentClient
}

type Torrent struct {
    Hash       string    `json:"hash"`
    Name       string    `json:"name"`
    Size       int64     `json:"size"`
    Progress   float64   `json:"progress"`
    Status     string    `json:"status"`
    Seeds      int       `json:"seeds"`
    Peers      int       `json:"peers"`
    DownSpeed  int64     `json:"down_speed"`
    UpSpeed    int64     `json:"up_speed"`
    AddedDate  time.Time `json:"added_date"`
    SavePath   string    `json:"save_path"`
    Comment    string    `json:"comment"`
    IsPrivate  bool      `json:"is_private"`
}

func NewTorrentService(endpoint string) *TorrentService {
    return &TorrentService{
        client: NewRTorrentClient(endpoint),
    }
}

func (s *TorrentService) GetTorrents() ([]Torrent, error) {
    // Get list of torrent hashes
    hashes, err := s.client.GetDownloadList()
    if err != nil {
        return nil, fmt.Errorf("error getting torrent list: %w", err)
    }

    torrents := make([]Torrent, 0, len(hashes))
    for _, hash := range hashes {
        info, err := s.client.GetTorrentInfo(hash)
        if err != nil {
            continue // Skip torrents that fail to load
        }

        // Calculate progress
        size := info["size"].(int64)
        completed := info["completed"].(int64)
        var progress float64
        if size > 0 {
            progress = float64(completed) / float64(size) * 100
        }

        // Convert state to status string
        status := "unknown"
        switch state := info["state"].(int); state {
        case 0:
            status = "stopped"
        case 1:
            if progress == 100 {
                status = "seeding"
            } else {
                status = "downloading"
            }
        }

        torrent := Torrent{
            Hash:      hash,
            Name:      info["name"].(string),
            Size:      size,
            Progress:  progress,
            Status:    status,
            Seeds:     info["seeders"].(int),
            Peers:     info["peers"].(int),
            DownSpeed: info["download_rate"].(int64),
            UpSpeed:   info["upload_rate"].(int64),
        }

        torrents = append(torrents, torrent)
    }

    return torrents, nil
}

func (s *TorrentService) AddTorrent(data []byte, start bool) error {
    _, err := s.client.Call("load.raw_start", string(data))
    if err != nil {
        return fmt.Errorf("error adding torrent: %w", err)
    }
    return nil
}

func (s *TorrentService) AddMagnet(uri string, start bool) error {
    method := "load.start"
    if !start {
        method = "load.normal"
    }
    
    _, err := s.client.Call(method, uri)
    if err != nil {
        return fmt.Errorf("error adding magnet: %w", err)
    }
    return nil
}

func (s *TorrentService) StartTorrent(hash string) error {
    _, err := s.client.Call("d.start", hash)
    return err
}

func (s *TorrentService) StopTorrent(hash string) error {
    _, err := s.client.Call("d.stop", hash)
    return err
}

func (s *TorrentService) DeleteTorrent(hash string) error {
    _, err := s.client.Call("d.erase", hash)
    return err
}

func (s *TorrentService) GetTorrentDetails(hash string) (*Torrent, error) {
    info, err := s.client.GetTorrentInfo(hash)
    if err != nil {
        return nil, err
    }

    size := info["size"].(int64)
    completed := info["completed"].(int64)
    var progress float64
    if size > 0 {
        progress = float64(completed) / float64(size) * 100
    }

    torrent := &Torrent{
        Hash:      hash,
        Name:      info["name"].(string),
        Size:      size,
        Progress:  progress,
        Seeds:     info["seeders"].(int),
        Peers:     info["peers"].(int),
        DownSpeed: info["download_rate"].(int64),
        UpSpeed:   info["upload_rate"].(int64),
    }

    return torrent, nil
}