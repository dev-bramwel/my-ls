// Package display provides output formatting functions for my-ls.[cite: 2]
// This file contains the long format (-l) output generation.[cite: 2]
// The functions here create output matching the standard ls -l command format.[cite: 2]
package display

import (
	"fmt"
	"my-ls/pkg/fs"
	"strconv"
	"time"
)

const (
	Reset  = "\033[0m"   //[cite: 2]
	Blue   = "\033[1;34m" // For Directories[cite: 2]
	Green  = "\033[1;32m" // For Executables[cite: 2]
	Cyan   = "\033[1;36m" // For Symlinks (Optional bonus styling)[cite: 2]
)

// FormatName applies ANSI color codes based on the file type and permissions.[cite: 2]
func FormatName(file fs.FileInfo) string {
	// 1. Check if it's a directory[cite: 2]
	if file.IsDir {
		return fmt.Sprintf("%s%s%s", Blue, file.Name, Reset) //[cite: 2]
	}

	// 2. Check if the file is executable (looks at owner, group, or other execute bits)[cite: 2]
	// os.ModePerm is a bitmask for all runtime permissions (0777)[cite: 2]
	const executeMask = 0111 // Looks for any 'x' bit (--x--x--x)[cite: 2]
	if file.Mode&executeMask != 0 {
		return fmt.Sprintf("%s%s%s", Green, file.Name, Reset) //[cite: 2]
	}

	// 3. Default plain output for standard files[cite: 2]
	return file.Name //[cite: 2]
}

// FormatLongWithPadding generates a single line of ls -l output for a file entry with active spacing.
// Parameters:
//   - file: FileInfo struct containing file metadata
//   - maxLinks, maxOwner, maxGroup, maxSize: integer padding metrics computed dynamically from directory
//
// Returns:
//   - string: formatted line resembling "drwxr-xr-x  2 user group 4096 Jan  2 15:04 dirName"[cite: 3]
func FormatLongWithPadding(file fs.FileInfo, maxLinks, maxOwner, maxGroup, maxSize int) string {
	// Build the permission string from mode bits[cite: 3]
	perms := formatPermissions(file.Mode)

	// Format link count with standard ls width[cite: 3]
	linkCount := strconv.FormatUint(file.LinkCount, 10)

	// Format size with standard ls width[cite: 3]
	size := strconv.FormatInt(file.Size, 10)

	// Format modification time like ls does[cite: 3]
	date := formatDate(file.ModTime)

	coloredName := getColorizedName(file.Name, file.Mode)

	// Build the name - include symlink target if applicable[cite: 3]
	name := coloredName
	if file.IsSymlink && file.SymlinkTarget != "" {
		name = coloredName + " -> " + file.SymlinkTarget
	}

	// Build final output with proper spacing[cite: 3]
	// %*s handles standard right-alignment for metrics (Link count and sizes)
	// %-*s handles standard left-alignment for text blocks (Owners and groups)
	// Change the spaces between %-*s columns from double to single
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

// formatPermissions creates the 10-character permission string.[cite: 3]
// Parameters:
//   - mode: Unix file mode bits (includes file type in high bits)[cite: 3]
//
// Returns:
//   - string: 10-character permission string like "drwxr-xr-x" or "lrwxrwxrwx"[cite: 3]
func formatPermissions(mode uint32) string {
	// Extract file type from high bits (S_IFMT = 0o170000)[cite: 3]
	fileType := mode & 0o170000

	// Extract permission bits (bits 0-9)[cite: 3]
	perms := mode & 0o7777

	// Determine file type character[cite: 3]
	var firstChar string
	switch fileType {
	case 0o040000: // S_IFDIR[cite: 3]
		firstChar = "d"
	case 0o120000: // S_IFLNK[cite: 3]
		firstChar = "l"
	case 0o020000: // S_IFBLK[cite: 3]
		firstChar = "b"
	case 0o060000: // S_IFCHR[cite: 3]
		firstChar = "c"
	case 0o010000: // S_IFIFO[cite: 3]
		firstChar = "p"
	case 0o140000: // S_IFSOCK[cite: 3]
		firstChar = "s"
	default:
		firstChar = "-"
	}

	result := firstChar

	// Owner permissions (S_IRUSR=0o400, S_IWUSR=0o200, S_IXUSR=0o100)[cite: 3]
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
		if perms&0o4000 != 0 { // S_ISUID[cite: 3]
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

	// Group permissions (S_IRGRP=0o40, S_IWGRP=0o20, S_IXGRP=0o10)[cite: 3]
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
		if perms&0o2000 != 0 { // S_ISGID[cite: 3]
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

	// Other permissions (S_IROTH=0o4, S_IWOTH=0o2, S_IXOTH=0o1)[cite: 3]
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
		if perms&0o1000 != 0 { // S_ISVTX[cite: 3]
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

// formatDate creates the date string like ls displays.[cite: 3]
// Parameters:
//   - t: modification time to format[cite: 3]
//
// Returns:
//   - string: formatted date like "Jan  2 15:04" or "Jan  2  2006"[cite: 3]
func formatDate(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0) //[cite: 3]

	if t.After(sixMonthsAgo) {
		return t.Format("Jan _2 15:04") //[cite: 3]
	} else {
		return t.Format("Jan _2  2006") //[cite: 3]
	}
}