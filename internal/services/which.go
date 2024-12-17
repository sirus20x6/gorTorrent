// internal/services/which.go
package services

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "sync"
)

// WhichCache stores cached executable paths
type WhichCache struct {
    filePaths map[string]string
    changed   bool
    mu        sync.RWMutex
}

// NewWhichCache creates a new which cache instance
func NewWhichCache() (*WhichCache, error) {
    cache := &WhichCache{
        filePaths: make(map[string]string),
        changed:   true,
    }

    // Load cached paths from disk
    if err := cache.load(); err != nil {
        return nil, fmt.Errorf("failed to load which cache: %w", err)
    }

    return cache, nil
}

// GetFilePath returns the full path for an executable
func (w *WhichCache) GetFilePath(exe string) (string, error) {
    w.mu.RLock()
    path, exists := w.filePaths[exe]
    w.mu.RUnlock()

    if exists && isExecutable(path) {
        return path, nil
    }

    // Try to find the executable
    path, err := w.findExecutable(exe)
    if err != nil {
        return "", err
    }

    // Cache the result
    w.mu.Lock()
    w.filePaths[exe] = path
    w.changed = true
    w.mu.Unlock()

    return path, nil
}

// findExecutable searches for an executable in PATH
func (w *WhichCache) findExecutable(exe string) (string, error) {
    // First try exec.LookPath
    path, err := exec.LookPath(exe)
    if err == nil {
        return path, nil
    }

    // Fallback to manual PATH search
    pathEnv := os.Getenv("PATH")
    for _, dir := range filepath.SplitList(pathEnv) {
        path := filepath.Join(dir, exe)
        if isExecutable(path) {
            return path, nil
        }
    }

    return "", fmt.Errorf("executable %s not found in PATH", exe)
}

// isExecutable checks if a file exists and is executable
func isExecutable(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }

    // Check if it's a regular file and has execute permission
    return !info.IsDir() && info.Mode()&0111 != 0
}

// load loads the cache from disk
func (w *WhichCache) load() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    cache, err := NewCache("")
    if err != nil {
        return err
    }

    type cacheData struct {
        FilePaths map[string]string `json:"filePaths"`
        Changed   bool             `json:"changed"`
    }

    var data cacheData
    if err := cache.Get("which.dat", &data); err != nil {
        return err
    }

    if data.FilePaths != nil {
        w.filePaths = data.FilePaths
        w.changed = false

        // Validate cached paths
        for exe, path := range w.filePaths {
            if !isExecutable(path) {
                delete(w.filePaths, exe)
                w.changed = true
            }
        }
    }

    return nil
}

// Save persists the cache to disk
func (w *WhichCache) Save() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    if !w.changed {
        return nil
    }

    cache, err := NewCache("")
    if err != nil {
        return err
    }

    data := struct {
        FilePaths map[string]string `json:"filePaths"`
        Changed   bool             `json:"changed"`
    }{
        FilePaths: w.filePaths,
        Changed:   false,
    }

    if err := cache.Set("which.dat", data); err != nil {
        return err
    }

    w.changed = false
    return nil
}

// PruneCache removes invalid entries from the cache
func (w *WhichCache) PruneCache() {
    w.mu.Lock()
    defer w.mu.Unlock()

    for exe, path := range w.filePaths {
        if !isExecutable(path) {
            delete(w.filePaths, exe)
            w.changed = true
        }
    }
}

// singleton instance
var (
    instance *WhichCache
    once     sync.Once
)

// GetInstance returns the singleton WhichCache instance
func GetInstance() (*WhichCache, error) {
    var err error
    once.Do(func() {
        instance, err = NewWhichCache()
    })
    return instance, err
}

// FindExe is a convenience function to find an executable
func FindExe(exe string) (string, error) {
    cache, err := GetInstance()
    if err != nil {
        return "", err
    }
    return cache.GetFilePath(exe)
}

// WithExternals helps manage paths to external executables
type WithExternals struct {
    paths map[string]string
}

// SetExternalPath sets a custom path for an executable
func (w *WithExternals) SetExternalPath(exe, path string) {
    if w.paths == nil {
        w.paths = make(map[string]string)
    }
    w.paths[exe] = path
}

// FindExe finds an executable, checking custom paths first
func (w *WithExternals) FindExe(exe string) (string, error) {
    // Check custom paths first
    if path, ok := w.paths[exe]; ok {
        if isExecutable(path) {
            return path, nil
        }
    }

    // Fall back to regular path search
    return FindExe(exe)
}