package display

import (
	"fmt"
	"my-ls/pkg/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	Reset      = "\033[0m"
	Blue       = "\033[1;34m" // For Directories
	Green      = "\033[1;32m" // For Executables
	Cyan       = "\033[1;36m" // For Symlinks
	DeviceOpts = "\033[40;33;01m" // Bold Yellow text on Black background for Devices
)

func FormatLongWithPadding(file fs.FileInfo, maxLinks, maxOwner, maxGroup, maxSize int) string {
	perms := formatPermissions(file.Mode)
	linkCount := strconv.FormatUint(file.LinkCount, 10)
	size := strconv.FormatInt(file.Size, 10)
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
			targetColor = "\033[31m" // Red for broken/orphaned link targets
		}

		colorizedTarget := file.SymlinkTarget
		if targetColor != "" {
			colorizedTarget = fmt.Sprintf("%s%s%s", targetColor, file.SymlinkTarget, Reset)
		}

		name = coloredName + " -> " + colorizedTarget
	}

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

func formatPermissions(mode uint32) string {
	fileType := mode & 0o170000
	perms := mode & 0o7777

	var firstChar string
	switch fileType {
	case 0o040000:
		firstChar = "d"
	case 0o120000:
		firstChar = "l"
	case 0o020000:
		firstChar = "b"
	case 0o060000:
		firstChar = "c"
	case 0o010000:
		firstChar = "p"
	case 0o140000:
		firstChar = "s"
	default:
		firstChar = "-"
	}

	result := firstChar

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
		}
	} else {
		if perms&0o4000 != 0 {
			result += "S"
		} else {
			result += "-"
		}
	}

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
		}
	} else {
		if perms&0o2000 != 0 {
			result += "S"
		} else {
			result += "-"
		}
	}

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
		}
	} else {
		if perms&0o1000 != 0 {
			result += "T"
		} else {
			result += "-"
		}
	}

	return result
}

func formatDate(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if t.After(sixMonthsAgo) {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}