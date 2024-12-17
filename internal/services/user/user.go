// internal/services/user/user.go
package user

import (
    "fmt"
    "net"
    "os"
    "regexp"
    "strings"
    "sync"
)

var (
    userLoginInstance string
    localModeInstance *bool
    userMu           sync.RWMutex
    modeMu           sync.RWMutex
)

// Config holds user configuration
type Config struct {
    LocalHosts       []string
    ForbidUserSettings bool
    DefaultSCGIPort    int
    DefaultSCGIHost    string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
    return Config{
        LocalHosts: []string{
            "127.0.0.1",
            "localhost",
        },
        ForbidUserSettings: false,
        DefaultSCGIPort:    0,
        DefaultSCGIHost:    "localhost",
    }
}

// Service provides user management functionality
type Service struct {
    config Config
}

// New creates a new user service
func New(config Config) *Service {
    return &Service{
        config: config,
    }
}

// GetLogin returns the current user login
func (s *Service) GetLogin() string {
    userMu.RLock()
    if userLoginInstance != "" {
        defer userMu.RUnlock()
        return userLoginInstance
    }
    userMu.RUnlock()

    userMu.Lock()
    defer userMu.Unlock()

    // Check REMOTE_USER environment variable first
    userLoginInstance = os.Getenv("REMOTE_USER")

    // Check PHP_AUTH_USER if REMOTE_USER not set
    if userLoginInstance == "" {
        userLoginInstance = os.Getenv("PHP_AUTH_USER")
    }

    // Check REDIRECT_REMOTE_USER if still not set
    if userLoginInstance == "" {
        userLoginInstance = os.Getenv("REDIRECT_REMOTE_USER")
    }

    // Clean and validate username
    if userLoginInstance != "" {
        userLoginInstance = cleanUsername(userLoginInstance)
    }

    return userLoginInstance
}

// GetUser returns the current user for settings
func (s *Service) GetUser() string {
    if s.config.ForbidUserSettings {
        return ""
    }
    return s.GetLogin()
}

// IsLocalMode checks if running in local mode
func (s *Service) IsLocalMode(host string, port int) bool {
    modeMu.RLock()
    if localModeInstance != nil {
        defer modeMu.RUnlock()
        return *localModeInstance
    }
    modeMu.RUnlock()

    modeMu.Lock()
    defer modeMu.Unlock()

    if host == "" {
        host = s.config.DefaultSCGIHost
    }
    if port == 0 {
        port = s.config.DefaultSCGIPort
    }

    // Check if it's a local connection
    isLocal := false
    if port == 0 {
        isLocal = true
    } else {
        for _, localhost := range s.config.LocalHosts {
            if host == localhost {
                isLocal = true
                break
            }
        }
    }

    localModeInstance = &isLocal
    return isLocal
}

// Validate checks if a user is valid
func (s *Service) Validate(username string) (bool, error) {
    // Don't allow empty usernames
    if username == "" {
        return false, fmt.Errorf("empty username not allowed")
    }

    // Clean and check username
    cleaned := cleanUsername(username)
    if cleaned != username {
        return false, fmt.Errorf("invalid username characters")
    }

    return true, nil
}

// Helper functions

// cleanUsername removes invalid characters from username
func cleanUsername(username string) string {
    // Only allow alphanumeric, dash and underscore
    reg := regexp.MustCompile("[^a-z0-9\\-_]")
    return strings.ToLower(reg.ReplaceAllString(username, "_"))
}

// ValidateHost checks if a host is valid
func ValidateHost(host string) bool {
    // Check if it's a valid IP
    if ip := net.ParseIP(host); ip != nil {
        return true
    }

    // Check if it's a valid hostname
    if matched, err := regexp.MatchString(
        `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`,
        host); err == nil && matched {
        return true
    }

    return false
}

// Middleware creates user authentication middleware
func (s *Service) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get user from auth header
        username, _, ok := r.BasicAuth()
        if !ok {
            username = r.Header.Get("X-Remote-User")
        }

        // Validate user
        if username != "" {
            if valid, _ := s.Validate(username); !valid {
                http.Error(w, "Invalid user", http.StatusForbidden)
                return
            }
        }

        // Store user in context
        ctx := context.WithValue(r.Context(), "user", username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// FromContext extracts user from request context
func FromContext(ctx context.Context) string {
    if user, ok := ctx.Value("user").(string); ok {
        return user
    }
    return ""
}

// Package level convenience functions
var defaultService = New(DefaultConfig())

func GetLogin() string { return defaultService.GetLogin() }
func GetUser() string  { return defaultService.GetUser() }
func IsLocalMode(host string, port int) bool { 
    return defaultService.IsLocalMode(host, port) 
}