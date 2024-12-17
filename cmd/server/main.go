// cmd/server/main.go

func main() {
    // ... other setup ...

    th := handlers.NewTorrentHandler("http://localhost:5000")

    r := chi.NewRouter()
    r.Get("/torrents", th.GetTorrents)
    r.Post("/torrents/{hash}/start", th.StartTorrent)
    r.Post("/torrents/{hash}/pause", th.PauseTorrent)
    r.Delete("/torrents/{hash}", th.DeleteTorrent)

    // ... start server ...
}