package fs

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

// ReadDir scans targeted folder locations and converts files into generalized metadata struct lists
func ReadDir(path string, showHidden bool) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("cannot read directory: " + err.Error())
	}

	var files []FileInfo // files stores the aggregated slice of mapped descriptors

	// Invert parameters to prepend virtual directory targets if hidden tracking flags are active (-a)
	if showHidden {
		info, err := os.Lstat(path)
		if err == nil {
			owner, group := getOwnership(info)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				files = append(files, FileInfo{
					Name:          ".",
					Path:          path,
					IsDir:         true,
					IsCharDevice:  isCharDevice(info.Mode()),
					IsBlockDevice: isBlockDevice(info.Mode()),
					Size:          info.Size(),
					ModTime:       info.ModTime(),
					ModeString:    info.Mode().String(),
					Mode:          uint32(stat.Mode),
					LinkCount:     uint64(stat.Nlink),
					Owner:         owner,
					Group:         group,
					Blocks:        stat.Blocks,
					Rdev:          uint64(stat.Rdev),
				})
			}
		}

		parentPath := path + string(os.PathSeparator) + ".."
		info, err = os.Lstat(parentPath)
		if err == nil {
			owner, group := getOwnership(info)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				files = append(files, FileInfo{
					Name:          "..",
					Path:          path,
					IsDir:         true,
					IsCharDevice:  isCharDevice(info.Mode()),
					IsBlockDevice: isBlockDevice(info.Mode()),
					Size:          info.Size(),
					ModTime:       info.ModTime(),
					ModeString:    info.Mode().String(),
					Mode:          uint32(stat.Mode),
					LinkCount:     uint64(stat.Nlink),
					Owner:         owner,
					Group:         group,
					Blocks:        stat.Blocks,
					Rdev:          uint64(stat.Rdev),
				})
			}
		}
	}

	for _, entry := range entries {
		// Filter files beginning with a single dot if hidden verification state is inactive
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := path + string(os.PathSeparator) + entry.Name()

		// Lstat reads the link itself without resolving down to the linked file parameter destination
		lstat, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		info := lstat // Lstat keeps symlink entries and special files as their own filesystem nodes
		isSymlink := info.Mode()&os.ModeSymlink != 0

		owner, group := getOwnership(info)

		var mode uint32 // mode extracts raw numeric bits required to execute safe bitwise sorting masking
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			mode = uint32(stat.Mode)
		}

		var symlinkTarget string // symlinkTarget tracks the pointed location context string
		if isSymlink {
			target, err := os.Readlink(fullPath)
			if err == nil {
				symlinkTarget = target
			}
		}

		file := FileInfo{
			Name:          entry.Name(),
			Path:          fullPath,
			IsDir:         entry.IsDir() && !isSymlink,
			IsSymlink:     isSymlink,
			IsCharDevice:  isCharDevice(info.Mode()),
			IsBlockDevice: isBlockDevice(info.Mode()),
			SymlinkTarget: symlinkTarget,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			Mode:          mode,
			LinkCount:     uint64(info.Mode()),
			Owner:         owner,
			Group:         group,
		}

		// Extract raw low-level syscall properties directly out of POSIX wrappers
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			file.LinkCount = uint64(stat.Nlink)
			file.Blocks = stat.Blocks
			file.Rdev = uint64(stat.Rdev)
		}

		files = append(files, file)
	}

	return files, nil
}

// getOwnership extracts user account records mapped matching internal UID and GID constants
func getOwnership(info os.FileInfo) (owner string, group string) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		u, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
		if err != nil {
			owner = strconv.Itoa(int(stat.Uid))
		} else {
			owner = u.Username
		}

		g, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
		if err != nil {
			group = strconv.Itoa(int(stat.Gid))
		} else {
			group = g.Name
		}
	}
	return owner, group
}

// IsDirectory performs an isolated type validation pass on individual tracking paths
func IsDirectory(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// IsSymlink tracks validation metrics for symlinks to see if their targets are valid or broken
func IsSymlink(path string) (bool, bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, false, err
	}

	isLink := info.Mode()&os.ModeSymlink != 0
	if !isLink {
		return false, false, nil
	}

	_, err = os.Stat(path)
	targetExists := err == nil // targetExists evaluates true if the targeted link path opens correctly

	return true, targetExists, nil
}

// ReadFile processes target file metrics directly if parameter lists point directly to explicit files
func ReadFile(path string) (*FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0

	var symlinkTarget string
	if isSymlink {
		target, err := os.Readlink(path)
		if err == nil {
			symlinkTarget = target
		}
	}

	owner, group := getOwnership(info)

	file := &FileInfo{
		Name:          path, // Retains the full path parameter input value to match output syntax
		Path:          path,
		IsDir:         info.IsDir() && !isSymlink,
		IsSymlink:     isSymlink,
		IsCharDevice:  isCharDevice(info.Mode()),
		IsBlockDevice: isBlockDevice(info.Mode()),
		SymlinkTarget: symlinkTarget,
		Size:          info.Size(),
		ModTime:       info.ModTime(),
		Mode:          0,
		Owner:         owner,
		Group:         group,
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		file.LinkCount = uint64(stat.Nlink)
		file.Mode = uint32(stat.Mode)
		file.Blocks = stat.Blocks
		file.Rdev = uint64(stat.Rdev)
	}

	return file, nil
}

func isCharDevice(mode os.FileMode) bool {
	return mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0
}

func isBlockDevice(mode os.FileMode) bool {
	return mode&os.ModeDevice != 0 && mode&os.ModeCharDevice == 0
}
