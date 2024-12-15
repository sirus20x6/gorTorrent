// util/permission/permission.go
package permission

import (
    "os"
    "os/user"
    "strconv"
    "path/filepath"
)

// DoesUserHave checks if a user has the specified permissions on a file
func DoesUserHave(uid int, gids []int, file string, flags int) (bool, error) {
    // For root/superuser
    if uid <= 0 {
        if (flags & 0x0001) != 0 && !isDir(file) {
            fi, err := os.Stat(file)
            if err != nil {
                return false, err
            }
            return fi.Mode()&0o111 != 0, nil
        }
        return true, nil
    }

    // Handle symbolic links
    if isLink, err := isSymlink(file); err != nil {
        return false, err
    } else if isLink {
        target, err := filepath.EvalSymlinks(file)
        if err != nil {
            return false, err
        }
        file = target
    }

    // Check primary file permissions
    if ok, err := checkUserPermsPrimary(uid, gids, file, flags); err != nil {
        return false, err
    } else if !ok {
        return false, nil
    }

    // Check parent directory permissions
    dirFlags := 0x0001 // Read permission
    if (flags & 0x0002) != 0 && !isDir(file) {
        dirFlags = 0x0003 // Read and write permissions
    }

    return checkUserPermsPrimary(uid, gids, filepath.Dir(file), dirFlags)
}

// checkUserPermsPrimary checks primary file permissions
func checkUserPermsPrimary(uid int, gids []int, file string, flags int) (bool, error) {
    fi, err := os.Stat(file)
    if err != nil {
        return false, err
    }

    stat := fi.Sys().(*syscall.Stat_t)
    mode := fi.Mode()

    // Check world permissions
    if mode&os.FileMode(flags) == os.FileMode(flags) {
        return true, nil
    }

    // Check group permissions
    groupFlags := flags << 3
    for _, gid := range gids {
        if gid == int(stat.Gid) && mode&os.FileMode(groupFlags) == os.FileMode(groupFlags) {
            return true, nil
        }
    }

    // Check owner permissions
    ownerFlags := flags << 6
    return uid == int(stat.Uid) && mode&os.FileMode(ownerFlags) == os.FileMode(ownerFlags), nil
}

// Helper functions

func isDir(path string) bool {
    fi, err := os.Stat(path)
    return err == nil && fi.IsDir()
}

func isSymlink(path string) (bool, error) {
    fi, err := os.Lstat(path)
    if err != nil {
        return false, err
    }
    return fi.Mode()&os.ModeSymlink != 0, nil
}

// GetGroupIDs returns all group IDs for a given user
func GetGroupIDs(username string) ([]int, error) {
    u, err := user.Lookup(username)
    if err != nil {
        return nil, err
    }

    gids := []int{}
    groups, err := u.GroupIds()
    if err != nil {
        return nil, err
    }

    for _, gstr := range groups {
        gid, err := strconv.Atoi(gstr)
        if err != nil {
            continue
        }
        gids = append(gids, gid)
    }

    // Add primary group
    if pgid, err := strconv.Atoi(u.Gid); err == nil {
        gids = append(gids, pgid)
    }

    return gids, nil
}

// Common permission flag constants
const (
    PermRead    = 0x0001
    PermWrite   = 0x0002
    PermExecute = 0x0004
    PermRW      = PermRead | PermWrite
    PermRWX     = PermRead | PermWrite | PermExecute
)