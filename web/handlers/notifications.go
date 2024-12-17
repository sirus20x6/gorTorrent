// handlers/notifications.go
package handlers

import (
    "net/http"
)

type ToastData struct {
    Message string
    Type    string // "info", "success", "error", "warning"
}

func ShowNotification(w http.ResponseWriter, r *http.Request) {
    data := ToastData{
        Message: r.URL.Query().Get("message"),
        Type:    r.URL.Query().Get("type"),
    }

    // Render the toast partial
    err := templates.ExecuteTemplate(w, "toast", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Add HTMX specific headers
    w.Header().Set("HX-Trigger", "showToast")
}