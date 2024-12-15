// middleware/requests/requests.go
package requests

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "sync"
)

var (
    allowedMethods = map[string]bool{
        http.MethodGet:  true,
        http.MethodPost: true,
    }

    // Origins that are allowed for CORS/CSRF checking
    enabledOrigins = make(map[string]bool)
    originsLock   sync.RWMutex
)

// Config holds configuration for request validation
type Config struct {
    EnableCSRFCheck bool
    AllowedOrigins  []string
}

// Validator provides request validation
type Validator struct {
    config Config
}

// NewValidator creates a new request validator
func NewValidator(config Config) *Validator {
    // Initialize allowed origins
    originsLock.Lock()
    for _, origin := range config.AllowedOrigins {
        enabledOrigins[origin] = true
    }
    originsLock.Unlock()

    return &Validator{
        config: config,
    }
}

// AddAllowedOrigin adds an origin to the allowed list
func (v *Validator) AddAllowedOrigin(origin string) {
    originsLock.Lock()
    defer originsLock.Unlock()
    enabledOrigins[origin] = true
}

// ValidateMethod ensures the request method is allowed
func (v *Validator) ValidateMethod(r *http.Request) error {
    if !allowedMethods[r.Method] {
        return fmt.Errorf("method %s not allowed", r.Method)
    }
    return nil
}

// ValidateCSRF performs CSRF validation
func (v *Validator) ValidateCSRF(r *http.Request) error {
    if !v.config.EnableCSRFCheck {
        return nil
    }

    // Don't check GET requests
    if r.Method == http.MethodGet {
        return nil
    }

    // Check Origin header
    if origin := r.Header.Get("Origin"); origin != "" {
        if origin == "null" { // privacy-sensitive context
            return nil
        }

        originHost, err := extractHost(origin)
        if err != nil {
            return fmt.Errorf("invalid origin header: %w", err)
        }

        if isOriginAllowed(originHost) {
            return nil
        }
    }

    // Check Referer header
    if referer := r.Header.Get("Referer"); referer != "" {
        refererHost, err := extractHost(referer)
        if err != nil {
            return fmt.Errorf("invalid referer header: %w", err)
        }

        if isOriginAllowed(refererHost) {
            return nil
        }
    }

    // Check X-Forwarded-Host
    if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
        host := strings.SplitN(forwardedHost, ":", 2)[0]
        if isOriginAllowed(host) {
            return nil
        }
    }

    // Check Host header
    if host := r.Host; host != "" {
        host = strings.SplitN(host, ":", 2)[0]
        if isOriginAllowed(host) {
            return nil
        }
    }

    return fmt.Errorf("request origin not allowed")
}

// Middleware returns an http.Handler middleware that performs all validations
func (v *Validator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request method
        if err := v.ValidateMethod(r); err != nil {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        // Validate CSRF
        if err := v.ValidateCSRF(r); err != nil {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Helper functions

// extractHost extracts the host from a URL string
func extractHost(urlStr string) (string, error) {
    parsed, err := url.Parse(urlStr)
    if err != nil {
        return "", err
    }
    return parsed.Host, nil
}

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin string) bool {
    originsLock.RLock()
    defer originsLock.RUnlock()
    return enabledOrigins[origin]
}