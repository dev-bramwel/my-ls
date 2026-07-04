package display

import (
	"fmt"
	"my-ls/pkg/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ANSI terminal escape variables targeting specific layout color highlights
const (
	Reset      = "\033[0m"        // Disables colors, resetting the terminal brush back to native values
	Blue       = "\033[1;34m"     // Bold Blue escape tag for Directories
	Green      = "\033[1;32m"     // Bold Green escape tag for Executable Binaries
	Cyan       = "\033[1;36m"     // Bold Cyan escape tag for Symbolic Links
	DeviceOpts = "\033[40;33;01m" // Bold Yellow text over a Black background block for character/block devices
)

// FormatLongWithPadding calculates layout spaces dynamically using cell lengths computed in PrintLong
// FormatLongWithPadding calculates layout spaces dynamically using cell lengths computed in PrintLong
func FormatLongWithPadding(file fs.FileInfo, maxLinks, maxOwner, maxGroup, maxSize int) string {
	perms := formatPermissions(file)
	linkCount := strconv.FormatUint(file.LinkCount, 10)
	size := formatSizeOrDevice(file)
	date := formatDate(file.ModTime)
	coloredName := GetColorizedName(file.Name, file.Mode)

	name := coloredName
	if file.IsSymlink && file.SymlinkTarget != "" {
		targetColor := ""
		baseDir := filepath.Dir(file.Path)
		fullTargetCheckPath := filepath.Join(baseDir, file.SymlinkTarget)

		if info, err := os.Stat(fullTargetCheckPath); err == nil {
			if info.IsDir() {
				targetColor = Blue
			} else if info.Mode()&0111 != 0 {
				targetColor = Green
			}
		} else {
			targetColor = "\033[31m"
		}

		colorizedTarget := file.SymlinkTarget
		if targetColor != "" {
			colorizedTarget = fmt.Sprintf("%s%s%s", targetColor, file.SymlinkTarget, Reset)
		}

		name = coloredName + " -> " + colorizedTarget
	}

	// FIX: Explicitly match every layout verb width operator (* or -*) with its corresponding max constraint integer parameter
	return fmt.Sprintf("%s %*s %-*s %-*s %*s %s %s\n",
		perms,
		maxLinks, linkCount,
		maxOwner, file.Owner,
		maxGroup, file.Group,
		maxSize, size,
		date,
		name,
	)
}

// formatPermissions processes POSIX mode mask numbers to construct the exact 10-character string string layout
func formatPermissions(file fs.FileInfo) string {
	fileType := file.Mode & 0o170000 // fileType isolates the core system bits defining the type of file descriptor
	perms := file.Mode & 0o7777      // perms extracts the execution rights, setuid, setgid, and sticky bits

	var firstChar string
	switch {
	case fileType == 0o040000: // S_IFDIR
		firstChar = "d"
	case fileType == 0o120000: // S_IFLNK
		firstChar = "l"
	case file.IsCharDevice || fileType == 0o020000: // S_IFCHR
		firstChar = "c"
	case file.IsBlockDevice || fileType == 0o060000: // S_IFBLK
		firstChar = "b"
	case fileType == 0o010000: // S_IFIFO
		firstChar = "p"
	case fileType == 0o140000: // S_IFSOCK
		firstChar = "s"
	default:
		firstChar = "-"
	}

	result := firstChar // result string is aggregated progressively matching owner, group, and other bit scopes

	// Owner Bit Flags (S_IRUSR, S_IWUSR, S_IXUSR)
	if perms&0o400 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perms&0o200 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perms&0o100 != 0 {
		if perms&0o4000 != 0 {
			result += "s"
		} else {
			result += "x"
		} // S_ISUID check
	} else {
		if perms&0o4000 != 0 {
			result += "S"
		} else {
			result += "-"
		}
	}

	// Group Bit Flags (S_IRGRP, S_IWGRP, S_IXGRP)
	if perms&0o40 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perms&0o20 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perms&0o10 != 0 {
		if perms&0o2000 != 0 {
			result += "s"
		} else {
			result += "x"
		} // S_ISGID check
	} else {
		if perms&0o2000 != 0 {
			result += "S"
		} else {
			result += "-"
		}
	}

	// Other/World Bit Flags (S_IROTH, S_IWOTH, S_IXOTH)
	if perms&0o4 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perms&0o2 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perms&0o1 != 0 {
		if perms&0o1000 != 0 {
			result += "t"
		} else {
			result += "x"
		} // S_ISVTX (Sticky bit) check
	} else {
		if perms&0o1000 != 0 {
			result += "T"
		} else {
			result += "-"
		}
	}

	return result
}

func formatSizeOrDevice(file fs.FileInfo) string {
	if isDevice(file) {
		return fmt.Sprintf("%d, %d", major(file.Rdev), minor(file.Rdev))
	}
	return strconv.FormatInt(file.Size, 10)
}

func isDevice(file fs.FileInfo) bool {
	if file.IsCharDevice || file.IsBlockDevice {
		return true
	}
	fileType := file.Mode & 0o170000
	return fileType == 0o020000 || fileType == 0o060000
}

func major(dev uint64) uint64 {
	return ((dev >> 8) & 0xfff) | ((dev >> 32) & 0xfffff000)
}

func minor(dev uint64) uint64 {
	return (dev & 0xff) | ((dev >> 12) & 0xffffff00)
}

// formatDate converts system time objects into localized layout variations matching standard ls windows
func formatDate(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0) // sixMonthsAgo sets the boundary window threshold

	// If the file was modified within the last 6 months, show Month Day Hour:Minute.
	// If it's older, show Month Day Year instead.
	if t.After(sixMonthsAgo) {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}
