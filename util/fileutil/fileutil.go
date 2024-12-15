// util/fileutil/fileutil.go
package fileutil

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

var (
    profilePathInstance string
    profileMask        = os.FileMode(0777)
    mu                sync.RWMutex
)

// GetFileName extracts the filename from a path
func GetFileName(path string) string {
    return filepath.Base(path)
}

// AddSlash ensures path ends with a slash
func AddSlash(str string) string {
    if str == "" || str[len(str)-1] == filepath.Separator {
        return str
    }
    return str + string(filepath.Separator)
}

// DelSlash removes trailing slash if present
func DelSlash(str string) string {
    if str == "" || str[len(str)-1] != filepath.Separator {
        return str
    }
    return str[:len(str)-1]
}

// FullPath returns absolute path with proper base
func FullPath(path, base string) string {
    if path == "" {
        return base
    }
    if filepath.IsAbs(path) {
        return path
    }
    if base == "" {
        var err error
        base, err = os.Getwd()
        if err != nil {
            return path
        }
    }
    return filepath.Clean(filepath.Join(AddSlash(base), path))
}

// GetProfilePath returns user profile path
func GetProfilePath() string {
    mu.RLock()
    if profilePathInstance != "" {
        defer mu.RUnlock()
        return profilePathInstance
    }
    mu.RUnlock()

    mu.Lock()
    defer mu.Unlock()
    profilePathInstance = getProfilePathEx("")
    return profilePathInstance
}

// GetProfilePathEx returns profile path for specific user
func GetProfilePathEx(user string) string {
    // Get base profile path from config or environment
    profilePath := os.Getenv("RUTORRENT_PROFILE_PATH")
    if profilePath == "" {
        profilePath = filepath.Join("share", "users")
    }

    profilePath = FullPath(profilePath, "")

    if user != "" {
        profilePath = filepath.Join(profilePath, user)
        // Create user directories if they don't exist
        dirs := []string{
            profilePath,
            filepath.Join(profilePath, "settings"),
            filepath.Join(profilePath, "torrents"),
            filepath.Join(profilePath, "tmp"),
        }
        MakeDirectory(dirs)
    }

    return profilePath
}

// GetSettingsPath returns settings path for current user
func GetSettingsPath() string {
    return filepath.Join(GetProfilePath(), "settings")
}

// GetSettingsPathEx returns settings path for specific user
func GetSettingsPathEx(user string) string {
    return filepath.Join(GetProfilePathEx(user), "settings")
}

// GetUploadsPath returns uploads path for current user
func GetUploadsPath() string {
    return filepath.Join(GetProfilePath(), "torrents")
}

// GetUploadsPathEx returns uploads path for specific user
func GetUploadsPathEx(user string) string {
    return filepath.Join(GetProfilePathEx(user), "torrents")
}

// GetTempDirectory returns temporary directory
func GetTempDirectory() string {
    // Check common temp directories
    dirs := []string{
        os.Getenv("TEMP"),
        os.Getenv("TMP"),
        "/tmp",
    }

    for _, dir := range dirs {
        if dir != "" {
            if info, err := os.Stat(dir); err == nil && info.IsDir() {
                if isWritable(dir) {
                    return AddSlash(dir)
                }
            }
        }
    }

    // Fallback to user temp directory
    return AddSlash(filepath.Join(GetProfilePath(), "tmp"))
}

// GetTempFilename generates a unique temporary filename
func GetTempFilename(purpose, extension string) string {
    parts := []string{
        "rutorrent",
        purpose,
        fmt.Sprint(os.Getpid()),
    }
    name := strings.Join(parts, "-")

    if extension != "" && !strings.HasPrefix(extension, ".") {
        extension = "." + extension
    }

    for {
        path := filepath.Join(GetTempDirectory(), fmt.Sprintf("%s-%d%s",
            name, time.Now().UnixNano(), extension))
        if _, err := os.Stat(path); os.IsNotExist(err) {
            return path
        }
    }
}

// GetUniqueFilename ensures filename doesn't exist by adding counter
func GetUniqueFilename(fname string) string {
    if !FileExists(fname) {
        return fname
    }

    ext := filepath.Ext(fname)
    base := strings.TrimSuffix(fname, ext)
    counter := 1

    for {
        newName := fmt.Sprintf("%s(%d)%s", base, counter, ext)
        if !FileExists(newName) {
            return newName
        }
        counter++
    }
}

// GetUniqueUploadedFilename ensures unique filename in uploads directory
func GetUniqueUploadedFilename(fname string) string {
    return GetUniqueFilename(filepath.Join(GetUploadsPath(), fname))
}

// MakeDirectory creates directories with proper permissions
func MakeDirectory(dirs interface{}) error {
    mode := profileMask
    if mode == 0 {
        mode = 0777
    }

    createDir := func(dir string) error {
        if err := os.MkdirAll(dir, mode); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
        return nil
    }

    switch d := dirs.(type) {
    case string:
        return createDir(d)
    case []string:
        for _, dir := range d {
            if err := createDir(dir); err != nil {
                return err
            }
        }
    }
    return nil
}

// DeleteDirectory removes directory and contents
func DeleteDirectory(dir string) error {
    return os.RemoveAll(dir)
}

// FileExists checks if file exists
func FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// CopyFile copies a file with proper permissions
func CopyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destination.Close()

    if _, err := io.Copy(destination, source); err != nil {
        return fmt.Errorf("failed to copy file: %w", err)
    }

    sourceInfo, err := os.Stat(src)
    if err == nil {
        err = os.Chmod(dst, sourceInfo.Mode())
    }
    return err
}

// isWritable checks if directory is writable
func isWritable(path string) bool {
    // Try to create a temporary file
    tempFile := filepath.Join(path, fmt.Sprintf(".test_%d", time.Now().UnixNano()))
    file, err := os.Create(tempFile)
    if err != nil {
        return false
    }
    file.Close()
    os.Remove(tempFile)
    return true
}

// GetPluginConf returns plugin configuration file path
func GetPluginConf(plugin string) string {
    paths := []string{
        filepath.Join("plugins", plugin, "conf.php"),
        filepath.Join("conf", "users", GetProfilePath(), "plugins", plugin, "conf.php"),
    }

    for _, path := range paths {
        if FileExists(path) {
            return path
        }
    }
    return ""
}

// GetConfFile returns configuration file path
func GetConfFile(name string) string {
    path := filepath.Join("conf", "users", GetProfilePath(), name)
    if FileExists(path) {
        return path
    }
    return ""
}