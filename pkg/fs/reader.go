package fs

import (
	// os - provides ReadDir and Stat for directory/file operations
	"os"
	// os/user - provides user lookup functionality for owner/group resolution
	"os/user"
	// strconv - provides string conversion utilities
	"strconv"
	// strings - provides string manipulation utilities for hidden file detection
	"strings"
	// errors - provides error creation utilities
	"errors"
	// syscall - provides Stat_t type for Unix file metadata extraction
	"syscall"
)

// ReadDir reads all entries from a directory and returns FileInfo structs.
// Parameters:
//   - path: directory path to read (required)
//   - showHidden: when true, includes files starting with '.' (default false)
//
// Returns:
//   - []FileInfo: slice of file metadata structs
//   - error: non-nil if directory cannot be read or stat fails
//
// Scope: Performs os.ReadDir followed by os.Stat on each entry to populate all metadata fields.
// When showHidden is true, adds "." and ".." entries at the beginning.
func ReadDir(path string, showHidden bool) ([]FileInfo, error) {
	// os.ReadDir reads directory entries without requiring separate stat calls for basic info
	entries, err := os.ReadDir(path)
	if err != nil {
		// errors.New creates a formatted error string
		return nil, errors.New("cannot read directory: " + err.Error())
	}

	// Result slice holds FileInfo structs for valid entries
	var files []FileInfo

	// When showHidden is true, we include "." and ".." entries
	// These represent current directory and parent directory respectively
	if showHidden {
		// Add "." (current directory) entry
		info, err := os.Stat(path)
		if err == nil {
			owner, group := getOwnership(info)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				files = append(files, FileInfo{
					Name:       ".",
					Path:       path,
					IsDir:      true,
					Size:       0, // directories show 0 size in ls -l
					ModTime:    info.ModTime(),
					ModeString: info.Mode().String(),
					Mode:       uint32(stat.Mode),
					LinkCount:  uint64(stat.Nlink),
					Owner:      owner,
					Group:      group,
					Blocks:     stat.Blocks, // Capture blocks
				})
			}
		}

		// Add ".." (parent directory) entry
		parentPath := path + string(os.PathSeparator) + ".."
		info, err = os.Stat(parentPath)
		if err == nil {
			owner, group := getOwnership(info)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				files = append(files, FileInfo{
					Name:       ".",
					Path:       path,
					IsDir:      true,
					Size:       0, 
					ModTime:    info.ModTime(),
					ModeString: info.Mode().String(),
					Mode:       uint32(stat.Mode),
					LinkCount:  uint64(stat.Nlink),
					Owner:      owner,
					Group:      group,
					Blocks:     stat.Blocks, // Capture blocks
				})
			}
		}
	}

	for _, entry := range entries {
		// Filter hidden files unless showHidden is true
		// strings.HasPrefix checks for '.' prefix on filenames
		// Note: "." and ".." are already handled above
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Build full path for stat operations
		fullPath := path + string(os.PathSeparator) + entry.Name()

		// Use Lstat to detect symlinks without following them
		lstat, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		// Check if this is a symlink
		isSymlink := lstat.Mode()&os.ModeSymlink != 0

		// Get file info - use Lstat for symlinks, Stat for other files
		var info os.FileInfo
		if isSymlink {
			info = lstat
		} else {
			info, err = entry.Info()
			if err != nil {
				continue
			}
		}

		// Extract ownership information from stat result
		owner, group := getOwnership(info)

		// Detect file type and set appropriate mode string
		var mode uint32
		var modeStr string
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			mode = uint32(stat.Mode)
		}

		// For symlinks, read the target
		var symlinkTarget string
		if isSymlink {
			target, err := os.Readlink(fullPath)
			if err == nil {
				symlinkTarget = target
			}
		}

		// Create FileInfo struct with all required metadata
		file := FileInfo{
			Name:          entry.Name(),
			Path:          fullPath,
			IsDir:         entry.IsDir() && !isSymlink, // Symlinks report their link status
			IsSymlink:     isSymlink,
			SymlinkTarget: symlinkTarget,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			ModeString:    modeStr, // Will be computed from mode
			Mode:          mode,
			LinkCount:     uint64(info.Mode()), // Will be corrected below
			Owner:         owner,
			Group:         group,
		}

		// Extract link count from syscall.Stat_t
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			file.LinkCount = uint64(stat.Nlink)
			file.Blocks = stat.Blocks // Capture blocks here safely
		}

		files = append(files, file)
	}

	return files, nil
}

// getOwnership extracts owner and group names from file info.
// Uses os/user.LookupId to convert numeric UIDs/GIDs to names.
//
// Parameters:
//   - info: os.FileInfo from which to extract ownership
//
// Returns:
//   - string: username of file owner
//   - string: group name
//
// Scope: Creates temporary user/group lookup results. Falls back to numeric IDs on lookup failure.
func getOwnership(info os.FileInfo) (owner string, group string) {
	// Sys() returns platform-specific stat data; cast to syscall.Stat_t for Unix systems
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		// LookupUid converts numeric user ID to username
		// Error handling: falls back to numeric ID string if lookup fails
		u, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
		if err != nil {
			// Return numeric UID as string on lookup failure
			owner = strconv.Itoa(int(stat.Uid))
		} else {
			owner = u.Username
		}

		// LookupGroupId converts numeric group ID to group name
		// Error handling: falls back to numeric GID string if lookup fails
		g, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
		if err != nil {
			// Return numeric GID as string on lookup failure
			group = strconv.Itoa(int(stat.Gid))
		} else {
			group = g.Name
		}
	}

	return owner, group
}

// IsDirectory checks if a path refers to a directory.
// Uses os.Stat to follow symlinks (unlike os.Lstat which does not).
//
// Parameters:
//   - path: filesystem path to check
//
// Returns:
//   - bool: true if path is a directory
//   - error: non-nil if stat fails
//
// Scope: Performs single stat call, no side effects.
func IsDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// IsSymlink checks if a path is a symbolic link.
// Uses os.Lstat to not follow symlinks.
//
// Parameters:
//   - path: filesystem path to check
//
// Returns:
//   - bool: true if path is a symlink
//   - bool: true if symlink target exists and is accessible
//   - error: non-nil on stat failure
//
// Scope: Single Lstat call, no side effects.
func IsSymlink(path string) (bool, bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, false, err
	}

	// Mode()&os.ModeSymlink checks if this is a symlink type
	isLink := info.Mode()&os.ModeSymlink != 0
	if !isLink {
		return false, false, nil
	}

	// Check if target exists by following the link
	_, err = os.Stat(path)
	targetExists := err == nil

	return true, targetExists, nil
}

// ReadFile reads metadata for a single file or directory path.
// Used when -l is applied to a file path (not a directory).
//
// Parameters:
//   - path: filesystem path to read metadata for
//
// Returns:
//   - *FileInfo: pointer to file metadata struct
//   - error: non-nil on stat failure
//
// Scope: Creates single FileInfo from os.Stat result.
func ReadFile(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Check if it's a symlink
	isSymlink := info.Mode()&os.ModeSymlink != 0

	// Extract ownership
	owner, group := getOwnership(info)

	file := &FileInfo{
		Name:     info.Name(),
		Path:     path,
		IsDir:    info.IsDir(),
		IsSymlink: isSymlink,
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		Mode:     0,
		Owner:    owner,
		Group:    group,
	}

	// Extract link count and mode bits from syscall.Stat_t
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		file.LinkCount = uint64(stat.Nlink)
		file.Mode = uint32(stat.Mode)
	}

	return file, nil
}