// util/lfs/lfs.go
package lfs

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"
)

// FileStat represents file information
type FileStat struct {
    Dev     int64
    Ino     int64
    Mode    uint32
    Nlink   int64
    Uid     int64
    Gid     int64
    Size    int64
    Atime   int64
    Mtime   int64
    Ctime   int64
    Blksize int64
    Blocks  int64
}

// IsFile checks if a path is a regular file
func IsFile(fname string) bool {
    if fi, err := os.Stat(fname); err == nil {
        return fi.Mode().IsRegular()
    }
    return false
}

// IsReadable checks if a file is readable
func IsReadable(fname string) bool {
    file, err := os.Open(fname)
    if err != nil {
        return false
    }
    file.Close()
    return true
}

// Stat returns file information handling large files correctly
func Stat(fname string) (*FileStat, error) {
    // Try native stat first
    fi, err := os.Stat(fname)
    if err != nil {
        return nil, fmt.Errorf("failed to stat file: %w", err)
    }

    sys := fi.Sys().(*syscall.Stat_t)
    stat := &FileStat{
        Dev:     int64(sys.Dev),
        Ino:     int64(sys.Ino),
        Mode:    uint32(sys.Mode),
        Nlink:   int64(sys.Nlink),
        Uid:     int64(sys.Uid),
        Gid:     int64(sys.Gid),
        Size:    fi.Size(),
        Atime:   sys.Atim.Sec,
        Mtime:   sys.Mtim.Sec,
        Ctime:   sys.Ctim.Sec,
        Blksize: int64(sys.Blksize),
        Blocks:  sys.Blocks,
    }

    // For large files on 32-bit systems, use external stat command
    if stat.Blocks == -1 || stat.Blocks > 4194303 {
        externStat, err := statExternal(fname)
        if err == nil {
            stat = externStat
        }
    }

    return stat, nil
}

// FileSize returns the size of a file handling large files correctly
func FileSize(fname string) (int64, error) {
    stat, err := Stat(fname)
    if err != nil {
        return 0, err
    }
    return stat.Size, nil
}

// FileMTime returns the modification time of a file
func FileMTime(fname string) (int64, error) {
    stat, err := Stat(fname)
    if err != nil {
        return 0, err
    }
    return stat.Mtime, nil
}

// statExternal uses the external stat command for large files
func statExternal(fname string) (*FileStat, error) {
    // Run stat command with format string
    cmd := exec.Command("stat",
        "-c%d:%i:%f:%h:%u:%g:%s:%X:%Y:%Z:%B:%b",
        fname)
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("external stat failed: %w", err)
    }

    // Parse stat output
    parts := strings.Split(strings.TrimSpace(string(output)), ":")
    if len(parts) != 12 {
        return nil, fmt.Errorf("invalid stat output format")
    }

    stat := &FileStat{}
    
    // Parse each field
    var parseErr error
    if stat.Dev, parseErr = strconv.ParseInt(parts[0], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse dev: %w", parseErr)
    }
    if stat.Ino, parseErr = strconv.ParseInt(parts[1], 16, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse ino: %w", parseErr)
    }
    if mode, parseErr := strconv.ParseUint(parts[2], 16, 32); parseErr != nil {
        return nil, fmt.Errorf("failed to parse mode: %w", parseErr)
    } else {
        stat.Mode = uint32(mode)
    }
    if stat.Nlink, parseErr = strconv.ParseInt(parts[3], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse nlink: %w", parseErr)
    }
    if stat.Uid, parseErr = strconv.ParseInt(parts[4], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse uid: %w", parseErr)
    }
    if stat.Gid, parseErr = strconv.ParseInt(parts[5], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse gid: %w", parseErr)
    }
    if stat.Size, parseErr = strconv.ParseInt(parts[6], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse size: %w", parseErr)
    }
    if stat.Atime, parseErr = strconv.ParseInt(parts[7], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse atime: %w", parseErr)
    }
    if stat.Mtime, parseErr = strconv.ParseInt(parts[8], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse mtime: %w", parseErr)
    }
    if stat.Ctime, parseErr = strconv.ParseInt(parts[9], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse ctime: %w", parseErr)
    }
    if stat.Blksize, parseErr = strconv.ParseInt(parts[10], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse blksize: %w", parseErr)
    }
    if stat.Blocks, parseErr = strconv.ParseInt(parts[11], 10, 64); parseErr != nil {
        return nil, fmt.Errorf("failed to parse blocks: %w", parseErr)
    }

    return stat, nil
}

// Test executes a test command on a file
func Test(fname string, flag string) bool {
    cmd := exec.Command("test", "-"+flag, fname)
    return cmd.Run() == nil
}

// GetMinFilePerms checks if a file has at least the specified permissions
func GetMinFilePerms(fname string) bool {
    fi, err := os.Stat(fname)
    if err != nil {
        return false
    }

    // Check minimum permissions (read and execute)
    minPerms := os.FileMode(0500) // r-x------
    return fi.Mode()&minPerms == minPerms
}