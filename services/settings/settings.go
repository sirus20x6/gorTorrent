// services/settings/settings.go
package settings

import (
    "encoding/json"
    "fmt"
    "sync"
)

// Settings manages application settings and plugins
type Settings struct {
    config     *Config
    plugins    map[string]interface{}
    aliases    map[string]CommandAlias
    hooks      map[string][]Hook
    mu         sync.RWMutex
}

// CommandAlias represents a command alias configuration
type CommandAlias struct {
    Name       string
    Command    string
    ParamCount int
}

// Hook represents an event hook configuration
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
            plugins: make(map[string]interface{}),
            aliases: make(map[string]CommandAlias),
            hooks:   make(map[string][]Hook),
        }
    })
    return instance
}

// Load loads settings from storage
func (s *Settings) Load(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Load config
    cfg, err := LoadConfig(configPath)
    if err != nil {
        return err
    }
    s.config = cfg

    // Create necessary directories
    if err := s.createRequiredDirs(); err != nil {
        return err
    }

    return nil
}

// Save persists settings to storage
func (s *Settings) Save(configPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    return SaveConfig(configPath, s.config)
}

// RegisterPlugin registers a plugin with optional data
func (s *Settings) RegisterPlugin(name string, data interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.plugins[name] = data
}

// UnregisterPlugin removes a plugin registration
func (s *Settings) UnregisterPlugin(name string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.plugins, name)
}

// GetPluginData returns plugin data if registered
func (s *Settings) GetPluginData(name string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.plugins[name]
}

// IsPluginRegistered checks if a plugin is registered
func (s *Settings) IsPluginRegistered(name string) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    _, exists := s.plugins[name]
    return exists
}

// RegisterAlias registers a command alias
func (s *Settings) RegisterAlias(name, cmd string, paramCount int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.aliases[name] = CommandAlias{
        Name:       name,
        Command:    cmd,
        ParamCount: paramCount,
    }
}

// GetCommand returns the real command for an alias
func (s *Settings) GetCommand(name string) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if alias, ok := s.aliases[name]; ok {
        return alias.Command
    }
    return name
}

// GetCommandParamCount returns the parameter count for a command
func (s *Settings) GetCommandParamCount(name string) int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if alias, ok := s.aliases[name]; ok {
        return alias.ParamCount
    }
    return 0
}

// RegisterEventHook registers an event hook
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

// GetEventHooks returns all hooks for an event
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

// createRequiredDirs creates necessary application directories
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