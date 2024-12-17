// internal/services/settings.go
package services

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
)

// Settings manages application settings, plugins, and UI preferences
type Settings struct {
    config     *Config      // Core application config
    plugins    map[string]interface{}
    aliases    map[string]CommandAlias
    hooks      map[string][]Hook
    uiData     string       // UI-specific settings as JSON string
    cache      *Cache
    modified   bool
    mu         sync.RWMutex
}

type CommandAlias struct {
    Name       string
    Command    string
    ParamCount int
}

type Hook struct {
    Name  string  `json:"name"`
    Level float64 `json:"level"`
}

var (
    instance *Settings
    once     sync.Once
)

// Get returns the singleton settings instance
func Get() *Settings {
    once.Do(func() {
        instance = &Settings{
            plugins:  make(map[string]interface{}),
            aliases:  make(map[string]CommandAlias),
            hooks:    make(map[string][]Hook),
            uiData:   "{}",
            modified: false,
        }
        
        // Initialize cache
        if cache, err := NewCache(""); err == nil {
            instance.cache = cache
            instance.loadUISettings()
        }
    })
    return instance
}

// Load loads core settings from config file
func (s *Settings) Load(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    cfg, err := LoadConfig(configPath) 
    if err != nil {
        return err
    }
    s.config = cfg

    if err := s.createRequiredDirs(); err != nil {
        return err
    }

    return nil
}

// Save persists all settings
func (s *Settings) Save(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Save core config
    if err := SaveConfig(configPath, s.config); err != nil {
        return err
    }

    // Save UI settings if modified
    if s.modified {
        if err := s.saveUISettings(); err != nil {
            return err
        }
    }

    return nil
}

// GetUISettings returns the UI settings JSON
func (s *Settings) GetUISettings() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.uiData
}

// SetUISettings updates UI settings
func (s *Settings) SetUISettings(jsonData string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Validate JSON
    var tmp interface{}
    if err := json.Unmarshal([]byte(jsonData), &tmp); err != nil {
        return fmt.Errorf("invalid JSON data: %w", err) 
    }

    if s.uiData != jsonData {
        s.uiData = jsonData
        s.modified = true
        return s.saveUISettings()
    }

    return nil
}

// GetUIValue gets a specific UI setting
func (s *Settings) GetUIValue(key string) (interface{}, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var settings map[string]interface{}
    if err := json.Unmarshal([]byte(s.uiData), &settings); err != nil {
        return nil, fmt.Errorf("failed to parse settings: %w", err)
    }

    return settings[key], nil
}

// SetUIValue sets a specific UI setting
func (s *Settings) SetUIValue(key string, value interface{}) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    var settings map[string]interface{}
    if err := json.Unmarshal([]byte(s.uiData), &settings); err != nil {
        settings = make(map[string]interface{})
    }

    settings[key] = value
    
    data, err := json.Marshal(settings)
    if err != nil {
        return fmt.Errorf("failed to marshal settings: %w", err)
    }

    s.uiData = string(data)
    s.modified = true
    return s.saveUISettings()
}

// Plugin management
func (s *Settings) RegisterPlugin(name string, data interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.plugins[name] = data
}

func (s *Settings) UnregisterPlugin(name string) {
    s.mu.Lock() 
    defer s.mu.Unlock()
    delete(s.plugins, name)
}

func (s *Settings) GetPluginData(name string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.plugins[name]
}

// Internal helper methods
func (s *Settings) loadUISettings() error {
    type cachedSettings struct {
        Modified bool   `json:"modified"`
        Data     string `json:"data"`
    }

    var settings cachedSettings
    if err := s.cache.Get("ui_settings", &settings); err != nil {
        return err
    }

    if settings.Data != "" {
        s.uiData = settings.Data
        s.modified = settings.Modified
    }

    return nil
}

func (s *Settings) saveUISettings() error {
    if s.cache == nil {
        return nil
    }

    return s.cache.Set("ui_settings", map[string]interface{}{
        "modified": s.modified,
        "data":     s.uiData,
    })
}

func (s *Settings) createRequiredDirs() error {
    dirs := []string{
        s.config.Server.TempDir,
        s.config.Server.DownloadDir,
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("error creating directory %s: %w", dir, err)
        }
    }

    return nil
}