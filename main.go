// cmd/main.go
package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/anacrolix/torrent"
)

type Application struct {
    torrentClient *torrent.Client
    templates     *template.Template
}

func main() {
    // Initialize torrent client
    cfg := torrent.NewDefaultClientConfig()
    cfg.DataDir = "./downloads"
    client, err := torrent.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Initialize templates
    tmpl, err := template.ParseGlob("templates/**/*.html")
    if err != nil {
        log.Fatal(err)
    }

    app := &Application{
        torrentClient: client,
        templates:     tmpl,
    }

    // Setup router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RealIP)

    // Static file server
    fileServer(r)

    // Routes
    r.Get("/", app.handleHome)
    r.Get("/torrents", app.handleTorrentList)
    r.Get("/torrents/{hash}", app.handleTorrentDetails)
    r.Post("/torrents/add", app.handleAddTorrent)
    r.Post("/torrents/{hash}/start", app.handleStartTorrent)
    r.Post("/torrents/{hash}/pause", app.handlePauseTorrent)
    r.Delete("/torrents/{hash}", app.handleDeleteTorrent)

    // HTMX Routes
    r.Get("/modals/add-torrent", app.handleAddTorrentModal)
    r.Get("/modals/settings", app.handleSettingsModal)
    r.Get("/speed", app.handleSpeedUpdate)

    // Start server
    port := 3000
    log.Printf("Starting server on http://localhost:%d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func fileServer(r chi.Router) {
    workDir, _ := filepath.Abs(".")
    filesDir := http.Dir(filepath.Join(workDir, "static"))
    r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))
}

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "ruTorrent Web",
    }
    app.render(w, "base", data)
}

func (app *Application) handleTorrentList(w http.ResponseWriter, r *http.Request) {
    torrents := app.getTorrents()
    if r.Header.Get("HX-Request") == "true" {
        app.render(w, "torrent-list", torrents)
        return
    }
    app.render(w, "base", torrents)
}

func (app *Application) handleAddTorrentModal(w http.ResponseWriter, r *http.Request) {
    app.render(w, "add-torrent", nil)
}

func (app *Application) handleSettingsModal(w http.ResponseWriter, r *http.Request) {
    app.render(w, "settings", nil)
}

func (app *Application) handleSpeedUpdate(w http.ResponseWriter, r *http.Request) {
    speeds := map[string]interface{}{
        "download": 0,
        "upload": 0,
    }
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("HX-Trigger", "speedUpdate")
    json.NewEncoder(w).Encode(speeds)
}

func (app *Application) render(w http.ResponseWriter, name string, data interface{}) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := app.templates.ExecuteTemplate(w, name+".html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (app *Application) getTorrents() map[string]interface{} {
    // Mock data for now
    return map[string]interface{}{
        "Torrents": []map[string]interface{}{
            {
                "Hash":      "abc123",
                "Name":      "Ubuntu 22.04 LTS",
                "Size":      "3.2 GB",
                "Progress":  67.5,
                "Status":    "downloading",
                "Seeds":     10,
                "Peers":     25,
                "DownSpeed": "2.1 MB/s",
                "UpSpeed":   "500 KB/s",
                "AddedDate": "2024-01-15 14:30",
            },
        },
    }
}