// internal/config/config.go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

// Config holds all application configuration
type Config struct {
    RTorrent struct {
        Endpoint string `json:"endpoint"`
        Username string `json:"username,omitempty"`
        Password string `json:"password,omitempty"`
    } `json:"rtorrent"`

    Server struct {
        Port         int    `json:"port"`
        Host         string `json:"host"` 
        BaseURL      string `json:"base_url"`
        TempDir      string `json:"temp_dir"`
        DownloadDir  string `json:"download_dir"`
    } `json:"server"`

    Auth struct {
        Enabled  bool   `json:"enabled"`
        Type     string `json:"type"` // "basic", "none"
        Username string `json:"username,omitempty"`
        Password string `json:"password,omitempty"`
    } `json:"auth"`
}

var defaultConfig = Config{
    RTorrent: struct {
        Endpoint string `json:"endpoint"`
        Username string `json:"username,omitempty"`
        Password string `json:"password,omitempty"`
    }{
        Endpoint: "http://localhost/RPC2",
    },
    Server: struct {
        Port         int    `json:"port"`
        Host         string `json:"host"`
        BaseURL      string `json:"base_url"`
        TempDir      string `json:"temp_dir"`
        DownloadDir  string `json:"download_dir"`
    }{
        Port:        3000,
        Host:        "127.0.0.1", 
        BaseURL:     "/",
        TempDir:     "/tmp/rutorrent",
        DownloadDir: "/downloads",
    },
    Auth: struct {
        Enabled  bool   `json:"enabled"`
        Type     string `json:"type"`
        Username string `json:"username,omitempty"`
        Password string `json:"password,omitempty"`
    }{
        Enabled: false,
        Type:    "none",
    },
}

// Settings manages application settings and plugins
type Settings struct {
    config     *Config
    aliases    map[string]CommandAlias
    plugins    map[string]interface{}
    hooks      map[string][]Hook
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
            config:  &defaultConfig,
            aliases: make(map[string]CommandAlias),
            plugins: make(map[string]interface{}),
            hooks:   make(map[string][]Hook),
        }
    })
    return instance
}

// Load loads settings from the specified file
func (s *Settings) Load(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // If config file exists, load it
    if configPath != "" {
        data, err := os.ReadFile(configPath)
        if err != nil {
            if !os.IsNotExist(err) {
                return fmt.Errorf("error reading config file: %w", err)
            }
            // File doesn't exist, use defaults
        } else {
            if err := json.Unmarshal(data, s.config); err != nil {
                return fmt.Errorf("error parsing config file: %w", err)
            }
        }
    }

    // Create necessary directories
    dirs := []string{s.config.Server.TempDir, s.config.Server.DownloadDir}
    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("error creating directory %s: %w", dir, err)
        }
    }

    return nil 
}

// Save saves the current settings to the specified file
func (s *Settings) Save(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    data, err := json.MarshalIndent(s.config, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling config: %w", err)
    }

    dir := filepath.Dir(configPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("error creating config directory: %w", err)
    }

    if err := os.WriteFile(configPath, data, 0644); err != nil {
        return fmt.Errorf("error writing config file: %w", err)
    }

    return nil
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

func (s *Settings) IsPluginRegistered(name string) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    _, exists := s.plugins[name]
    return exists
}

// Alias management
func (s *Settings) RegisterAlias(name, cmd string, paramCount int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.aliases[name] = CommandAlias{
        Name:       name,
        Command:    cmd,
        ParamCount: paramCount,
    }
}

func (s *Settings) GetCommand(name string) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if alias, ok := s.aliases[name]; ok {
        return alias.Command
    }
    return name
}

func (s *Settings) GetCommandParamCount(name string) int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if alias, ok := s.aliases[name]; ok {
        return alias.ParamCount
    }
    return 0
}

// Hook management
func (s *Settings) RegisterEventHook(event string, hook Hook) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.hooks[event] = append(s.hooks[event], hook)

    // Sort hooks by level
    hooks := s.hooks[event]
    for i := 0; i < len(hooks)-1; i++ {
        for j := i + 1; j < len(hooks); j++ {
            if hooks[i].Level > hooks[j].Level {
                hooks[i], hooks[j] = hooks[j], hooks[i]
            }
        }
    }
}

func (s *Settings) GetEventHooks(event string) []Hook {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.hooks[event]
}

// Config returns the current configuration
func (s *Settings) Config() *Config {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.config
}

// DefaultConfigPath returns the default path for the config file
func DefaultConfigPath() string {
    configDir, err := os.UserConfigDir()
    if err != nil {
        configDir = "."
    }
    return filepath.Join(configDir, "rutorrent", "config.json")
}