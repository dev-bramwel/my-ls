package fs

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func ReadDir(path string, showHidden bool) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("cannot read directory: " + err.Error())
	}

	var files []FileInfo

	if showHidden {
		info, err := os.Stat(path)
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
					Blocks:     stat.Blocks,
				})
			}
		}

		parentPath := path + string(os.PathSeparator) + ".."
		info, err = os.Stat(parentPath)
		if err == nil {
			owner, group := getOwnership(info)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				files = append(files, FileInfo{
					Name:       "..",
					Path:       path,
					IsDir:      true,
					Size:       0,
					ModTime:    info.ModTime(),
					ModeString: info.Mode().String(),
					Mode:       uint32(stat.Mode),
					LinkCount:  uint64(stat.Nlink),
					Owner:      owner,
					Group:      group,
					Blocks:     stat.Blocks,
				})
			}
		}
	}

	for _, entry := range entries {
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := path + string(os.PathSeparator) + entry.Name()

		lstat, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		isSymlink := lstat.Mode()&os.ModeSymlink != 0

		var info os.FileInfo
		if isSymlink {
			info = lstat
		} else {
			info, err = entry.Info()
			if err != nil {
				continue
			}
		}

		owner, group := getOwnership(info)

		var mode uint32
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			mode = uint32(stat.Mode)
		}

		var symlinkTarget string
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
			SymlinkTarget: symlinkTarget,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			Mode:          mode,
			LinkCount:     uint64(info.Mode()),
			Owner:         owner,
			Group:         group,
		}

		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			file.LinkCount = uint64(stat.Nlink)
			file.Blocks = stat.Blocks
		}

		files = append(files, file)
	}

	return files, nil
}

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

func IsDirectory(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

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
	targetExists := err == nil

	return true, targetExists, nil
}

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
		Name:          path,
		Path:          path,
		IsDir:         info.IsDir() && !isSymlink,
		IsSymlink:     isSymlink,
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
	}

	return file, nil
}