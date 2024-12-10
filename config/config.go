// config/config.go
package config

import (
    "os"
    "fmt"
    "path/filepath"
    "encoding/json"
)

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

var cfg *Config

// Load loads the configuration from the specified file
func Load(configPath string) (*Config, error) {
    if cfg != nil {
        return cfg, nil
    }

    cfg = &Config{}
    *cfg = defaultConfig

    // If config file exists, load it
    if configPath != "" {
        data, err := os.ReadFile(configPath)
        if err != nil {
            if !os.IsNotExist(err) {
                return nil, fmt.Errorf("error reading config file: %w", err)
            }
            // File doesn't exist, use defaults
        } else {
            if err := json.Unmarshal(data, cfg); err != nil {
                return nil, fmt.Errorf("error parsing config file: %w", err)
            }
        }
    }

    // Create necessary directories
    dirs := []string{cfg.Server.TempDir, cfg.Server.DownloadDir}
    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return nil, fmt.Errorf("error creating directory %s: %w", dir, err)
        }
    }

    return cfg, nil
}

// Save saves the current configuration to the specified file
func Save(configPath string) error {
    if cfg == nil {
        return fmt.Errorf("no configuration loaded")
    }

    data, err := json.MarshalIndent(cfg, "", "    ")
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

// Get returns the current configuration
func Get() *Config {
    if cfg == nil {
        cfg = &Config{}
        *cfg = defaultConfig
    }
    return cfg
}

// Update updates the configuration with new values
func Update(newCfg Config) error {
    if cfg == nil {
        cfg = &Config{}
    }
    *cfg = newCfg
    return nil
}

// DefaultConfigPath returns the default path for the config file
func DefaultConfigPath() string {
    configDir, err := os.UserConfigDir()
    if err != nil {
        configDir = "."
    }
    return filepath.Join(configDir, "rutorrent", "config.json")
}

// Example config.json:
/*
{
    "rtorrent": {
        "endpoint": "http://localhost/RPC2",
        "username": "",
        "password": ""
    },
    "server": {
        "port": 3000,
        "host": "127.0.0.1",
        "base_url": "/",
        "temp_dir": "/tmp/rutorrent",
        "download_dir": "/downloads"
    },
    "auth": {
        "enabled": false,
        "type": "none",
        "username": "",
        "password": ""
    }
}
*/
