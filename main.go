// main.go
package main

import (
    "log"
    "net/http"
    "html/template"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // Static files
    fileServer := http.FileServer(http.Dir("static"))
    r.Handle("/static/*", http.StripPrefix("/static", fileServer))

    // Load templates
    tmpl := template.Must(template.ParseGlob("templates/**/*"))

    // Routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]interface{}{
            "Torrents": []map[string]interface{}{
                {
                    "Hash":      "abc123",
                    "Name":      "Ubuntu 22.04 LTS",
                    "Size":      "3.2 GB",
                    "Progress":  67.5,
                    "Status":    "downloading",
                    "Seeds":     10,
                    "Peers":     25,
                    "DownSpeed": "2.1 MB",
                    "UpSpeed":   "500 KB",
                    "AddedDate": "2024-01-15 14:30",
                },
                // Add more sample data
            },
        }
        tmpl.ExecuteTemplate(w, "base.html", data)
    })

    // Modal routes
    r.Get("/modals/add-torrent", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "add-torrent.html", nil)
    })

    // HTMX partial updates
    r.Get("/torrents", func(w http.ResponseWriter, r *http.Request) {
        // Same sample data as above
        data := map[string]interface{}{
            "Torrents": []map[string]interface{}{
                {
                    "Hash":      "abc123",
                    "Name":      "Ubuntu 22.04 LTS",
                    "Size":      "3.2 GB",
                    "Progress":  67.5,
                    "Status":    "downloading",
                    "Seeds":     10,
                    "Peers":     25,
                    "DownSpeed": "2.1 MB",
                    "UpSpeed":   "500 KB",
                    "AddedDate": "2024-01-15 14:30",
                },
            },
        }
        tmpl.ExecuteTemplate(w, "torrent-list.html", data)
    })

    log.Println("Server starting on http://localhost:3000")
    http.ListenAndServe(":3000", r)
}
