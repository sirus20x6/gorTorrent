// services/uisettings.go
package services

import (
    "encoding/json"
    "fmt"
    "sync"
)

// UISettings manages the web UI settings
type UISettings struct {
    hash     string
    modified bool
    data     string
    mu       sync.RWMutex
    cache    *Cache
}

// NewUISettings creates a new UI settings instance
func NewUISettings() (*UISettings, error) {
    cache, err := NewCache("")
    if err != nil {
        return nil, fmt.Errorf("failed to create cache: %w", err)
    }

    settings := &UISettings{
        hash:     "WebUISettings.dat",
        modified: false,
        data:     "{}",
        cache:    cache,
    }

    // Load existing settings if they exist
    if err := settings.load(); err != nil {
        return nil, fmt.Errorf("failed to load settings: %w", err)
    }

    return settings, nil
}

// load retrieves settings from the cache
func (s *UISettings) load() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    type cachedSettings struct {
        Hash     string `json:"hash"`
        Modified bool   `json:"modified"`
        Data     string `json:"jsonString"`
    }

    var settings cachedSettings
    if err := s.cache.Get(s.hash, &settings); err != nil {
        return fmt.Errorf("failed to get settings from cache: %w", err)
    }

    // Only update if we got data from cache
    if settings.Data != "" {
        s.data = settings.Data
        s.modified = settings.Modified
    }

    return nil
}

// Store saves the current settings to cache
func (s *UISettings) Store() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    return s.cache.Set(s.hash, map[string]interface{}{
        "hash":       s.hash,
        "modified":   s.modified,
        "jsonString": s.data,
    })
}

// Get returns the current settings JSON string
func (s *UISettings) Get() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

// Set updates the settings with new JSON data
func (s *UISettings) Set(jsonData string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Validate JSON
    var tmp interface{}
    if err := json.Unmarshal([]byte(jsonData), &tmp); err != nil {
        return fmt.Errorf("invalid JSON data: %w", err)
    }

    if s.data != jsonData {
        s.data = jsonData
        s.modified = true
        return s.Store()
    }

    return nil
}

// GetValue retrieves a specific setting value
func (s *UISettings) GetValue(key string) (interface{}, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var settings map[string]interface{}
    if err := json.Unmarshal([]byte(s.data), &settings); err != nil {
        return nil, fmt.Errorf("failed to parse settings JSON: %w", err)
    }

    if value, ok := settings[key]; ok {
        return value, nil
    }
    return nil, nil
}

// SetValue updates a specific setting value
func (s *UISettings) SetValue(key string, value interface{}) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    var settings map[string]interface{}
    if err := json.Unmarshal([]byte(s.data), &settings); err != nil {
        settings = make(map[string]interface{})
    }

    settings[key] = value
    
    newData, err := json.Marshal(settings)
    if err != nil {
        return fmt.Errorf("failed to marshal settings: %w", err)
    }

    s.data = string(newData)
    s.modified = true
    return s.Store()
}

// LoadOrCreate loads existing settings or creates default ones
func LoadOrCreate() (*UISettings, error) {
    settings, err := NewUISettings()
    if err != nil {
        return nil, fmt.Errorf("failed to create settings: %w", err)
    }
    return settings, nil
}