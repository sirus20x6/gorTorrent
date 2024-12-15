// util/utility/utility.go
package utility

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

// GetExternal returns the path to an external executable
var externalPaths = make(map[string]string)

// SetExternalPath sets the path for an external executable
func SetExternalPath(exe, path string) {
    if path != "" {
        externalPaths[exe] = path
    }
}

// GetExternal returns the path to an external executable
func GetExternal(exe string) string {
    if path, ok := externalPaths[exe]; ok && path != "" {
        return path
    }
    return exe
}

// GetPHP returns the path to the PHP executable
func GetPHP() string {
    return GetExternal("php")
}

// QuoteAndDeslashEachItem quotes and escapes each item in a slice
func QuoteAndDeslashEachItem(items []string) []string {
    quoted := make([]string, len(items))
    for i, item := range items {
        quoted[i] = fmt.Sprintf("\"%s\"", 
            strings.NewReplacer(
                "\\", "\\\\",
                "\"", "\\\"",
                "\n", "\\n",
                "\r", "\\r",
                "\t", "\\t",
            ).Replace(item))
    }
    return quoted
}

// SortByTime sorts items by their time field
type TimeItem struct {
    Time int64
    Data interface{}
}

func SortByTime(items []TimeItem) {
    sort.Slice(items, func(i, j int) bool {
        return items[i].Time < items[j].Time
    })
}

// StringStartsWith checks if a string starts with another
func StringStartsWith(str, prefix string) bool {
    return strings.HasPrefix(str, prefix)
}

// StringEndsWith checks if a string ends with another
func StringEndsWith(str, suffix string) bool {
    return strings.HasSuffix(str, suffix)
}

// ExeExists checks if an executable exists in PATH
func ExeExists(exe string) bool {
    // Check custom path first
    if path := GetExternal(exe); path != exe {
        if _, err := os.Stat(path); err == nil {
            if info, err := os.Stat(path); err == nil && info.Mode()&0111 != 0 {
                return true
            }
        }
    }

    // Search in PATH
    if path, err := exec.LookPath(exe); err == nil {
        return true
    }

    return false
}

// FileSystem operations with error handling
type FileOps struct {
    errors []error
}

func NewFileOps() *FileOps {
    return &FileOps{
        errors: make([]error, 0),
    }
}

// CreateDir creates a directory and tracks any errors
func (f *FileOps) CreateDir(path string, mode os.FileMode) bool {
    if err := os.MkdirAll(path, mode); err != nil {
        f.errors = append(f.errors, fmt.Errorf("failed to create directory %s: %w", path, err))
        return false
    }
    return true
}

// WriteFile writes data to a file and tracks any errors
func (f *FileOps) WriteFile(path string, data []byte, mode os.FileMode) bool {
    if err := os.WriteFile(path, data, mode); err != nil {
        f.errors = append(f.errors, fmt.Errorf("failed to write file %s: %w", path, err))
        return false
    }
    return true
}

// CopyFile copies a file and tracks any errors
func (f *FileOps) CopyFile(src, dst string) bool {
    srcFile, err := os.Open(src)
    if err != nil {
        f.errors = append(f.errors, fmt.Errorf("failed to open source file %s: %w", src, err))
        return false
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        f.errors = append(f.errors, fmt.Errorf("failed to create destination file %s: %w", dst, err))
        return false
    }
    defer dstFile.Close()

    if _, err := io.Copy(dstFile, srcFile); err != nil {
        f.errors = append(f.errors, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err))
        return false
    }

    return true
}

// Errors returns all accumulated errors
func (f *FileOps) Errors() []error {
    return f.errors
}

// Path manipulation utilities
type PathUtil struct {
    base string
}

func NewPathUtil(base string) *PathUtil {
    return &PathUtil{base: base}
}

// JoinPath joins path elements using the correct separator
func (p *PathUtil) JoinPath(elem ...string) string {
    return filepath.Join(append([]string{p.base}, elem...)...)
}

// Clean cleans a path
func (p *PathUtil) Clean(path string) string {
    return filepath.Clean(path)
}

// RelPath returns a relative path
func (p *PathUtil) RelPath(target string) (string, error) {
    return filepath.Rel(p.base, target)
}

// IsSubPath checks if a path is under the base path
func (p *PathUtil) IsSubPath(path string) bool {
    rel, err := p.RelPath(path)
    if err != nil {
        return false
    }
    return !strings.Contains(rel, "..")
}

// Common configuration directories
func GetConfigDir() (string, error) {
    if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
        return dir, nil
    }
    
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    
    return filepath.Join(home, ".config"), nil
}

// Package level functions for backward compatibility
func QuoteString(str string) string {
    return QuoteAndDeslashEachItem([]string{str})[0]
}

func IsExecutable(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.Mode()&0111 != 0
}