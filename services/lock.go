// services/lock.go
package services

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// Lock implements a file-based locking mechanism
type Lock struct {
    file      string
    lockTime  time.Duration
    locked    bool
    mu        sync.Mutex
}

const (
    // MaxLockTime is the default maximum time a lock can be held
    MaxLockTime = 30 * time.Minute
)

// NewLock creates a new lock instance
func NewLock(name string, maxLockTime time.Duration) (*Lock, error) {
    if maxLockTime == 0 {
        maxLockTime = MaxLockTime
    }

    // Get settings path from config or environment
    settingsPath := os.Getenv("RUTORRENT_SETTINGS_PATH")
    if settingsPath == "" {
        settingsPath = "/var/lib/rutorrent/settings"
    }

    // Create lock file path
    lockFile := filepath.Join(settingsPath, name+".lock")

    return &Lock{
        file:     lockFile,
        lockTime: maxLockTime,
    }, nil
}

// Lock attempts to acquire the lock
func (l *Lock) Lock() bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    // Check if lock file exists
    info, err := os.Stat(l.file)
    if err == nil {
        // Lock file exists, check if it's stale
        if time.Since(info.ModTime()) > l.lockTime {
            // Lock is stale, remove it
            if err := os.Remove(l.file); err != nil {
                return false
            }
        } else {
            // Lock is still valid
            return false
        }
    } else if !os.IsNotExist(err) {
        // Error checking file
        return false
    }

    // Create lock file
    file, err := os.OpenFile(l.file, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
    if err != nil {
        return false
    }
    file.Close()

    l.locked = true
    return true
}

// Release releases the lock if held
func (l *Lock) Release() {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.locked {
        os.Remove(l.file)
        l.locked = false
    }
}

// IsLocked checks if the lock is currently held
func (l *Lock) IsLocked() bool {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.locked
}

// Obtain is a convenience function that creates and acquires a lock
func Obtain(name string, maxLockTime time.Duration) (*Lock, error) {
    lock, err := NewLock(name, maxLockTime)
    if err != nil {
        return nil, fmt.Errorf("failed to create lock: %w", err)
    }

    if !lock.Lock() {
        return nil, fmt.Errorf("failed to acquire lock")
    }

    return lock, nil
}

// WithLock executes a function with a lock
func WithLock(name string, maxLockTime time.Duration, f func() error) error {
    lock, err := Obtain(name, maxLockTime)
    if err != nil {
        return err
    }
    defer lock.Release()

    return f()
}