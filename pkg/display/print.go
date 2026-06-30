package display

import (
	// fmt - provides output formatting functions
	"fmt"
	// my-ls/pkg/fs - provides FileInfo type for formatting
	"my-ls/pkg/fs"
	// strconv - provides integer to string conversion
	"strconv"
	// time - provides time formatting utilities
	"time"
)

// FormatLong generates a single line of ls -l output for a file entry.
// Parameters:
//   - file: FileInfo struct containing file metadata
//
// Returns:
//   - string: formatted line resembling "drwxr-xr-x  2 user group 4096 Jan  2 15:04 dirName"
//
// Scope: Creates formatted string with proper spacing based on typical ls column widths.
// The output format matches system ls -l with:
//   - Permissions (10 chars): drwxrwxrwx format including special bits
//   - Link count (right-aligned): number of hard links
//   - Owner (left-aligned): username
//   - Group (left-aligned): group name
//   - Size (right-aligned): file size in bytes
//   - Date (12 chars): "Jan  2 15:04" or "Jan  2  2006" for old files
//   - Name: filename (with trailing -> for symlinks)
func FormatLong(file fs.FileInfo) string {
	// Build the permission string from mode bits
	perms := formatPermissions(file.Mode)

	// Format link count with standard ls width
	linkCount := strconv.FormatUint(file.LinkCount, 10)

	// Format size with standard ls width
	size := strconv.FormatInt(file.Size, 10)

	// Format modification time like ls does
	date := formatDate(file.ModTime)

	coloredName := getColorizedName(file.Name, file.Mode)

	// Build the name - include symlink target if applicable
	name := coloredName
	if file.IsSymlink && file.SymlinkTarget != "" {
		name = coloredName + " -> " + file.SymlinkTarget
	}

	// Build final output with proper spacing
	return fmt.Sprintf("%s %s %s %s %s %s %s\n",
		perms, linkCount, file.Owner, file.Group, size, date, name)
}

// formatPermissions creates the 10-character permission string.
// Parameters:
//   - mode: Unix file mode bits (includes file type in high bits)
//
// Returns:
//   - string: 10-character permission string like "drwxr-xr-x" or "lrwxrwxrwx"
//
// Scope: Maps mode bits to permission characters using standard Unix permissions.
// Handles setuid (s), setgid (s), and sticky (t) bits.
// File type indicators: d (dir), - (regular), l (symlink), c (char dev), b (block dev)
func formatPermissions(mode uint32) string {
	// Extract file type from high bits (S_IFMT = 0o170000)
	fileType := mode & 0o170000

	// Extract permission bits (bits 0-9)
	perms := mode & 0o7777

	// Determine file type character
	var firstChar string
	switch fileType {
	case 0o040000: // S_IFDIR
		firstChar = "d"
	case 0o120000: // S_IFLNK
		firstChar = "l"
	case 0o020000: // S_IFBLK
		firstChar = "b"
	case 0o060000: // S_IFCHR
		firstChar = "c"
	case 0o010000: // S_IFIFO
		firstChar = "p"
	case 0o140000: // S_IFSOCK
		firstChar = "s"
	default:
		firstChar = "-"
	}

	result := firstChar

	// Owner permissions (S_IRUSR=0o400, S_IWUSR=0o200, S_IXUSR=0o100)
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
		if perms&0o4000 != 0 { // S_ISUID
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

	// Group permissions (S_IRGRP=0o40, S_IWGRP=0o20, S_IXGRP=0o10)
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
		if perms&0o2000 != 0 { // S_ISGID
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

	// Other permissions (S_IROTH=0o4, S_IWOTH=0o2, S_IXOTH=0o1)
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
		if perms&0o1000 != 0 { // S_ISVTX
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

// formatDate creates the date string like ls displays.
// Parameters:
//   - t: modification time to format
//
// Returns:
//   - string: formatted date like "Jan  2 15:04" or "Jan  2  2006"
//
// Scope: Uses time.Time formatting to match ls output.
// Files older than 6 months show year instead of time.
func formatDate(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if t.After(sixMonthsAgo) {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}

// PrintStandard outputs filenames in standard ls format (space-separated on one line).
// Parameters:
//   - files: slice of fs.FileInfo to display
//
// Returns:
//   - outputs formatted filenames to stdout
//
// Scope: Iterates over files slice and prints each name followed by space.
// Each file name is printed on the same line, matching standard ls behavior.
// No return value - writes directly to stdout.
func PrintStandard(files []fs.FileInfo) {
	for _, file := range files {
		coloredName := getColorizedName(file.Name, file.Mode)
		fmt.Print(coloredName + " ")
	}
	fmt.Print("\n")
}

// PrintLong outputs files in long format (-l flag).
// Parameters:
//   - files: slice of fs.FileInfo to display in detail
//   - showTotal: when true, prints "total X" line at start
//
// Returns:
//   - outputs formatted detail lines to stdout
//
// Scope: Iterates over files and calls FormatLong for each.
// Prints "total X" line at start when showTotal is true.
// No return value - writes directly to stdout.
func PrintLong(files []fs.FileInfo, showTotal bool) {
	if showTotal {
		totalBlocks := 0
		for _, f := range files {
			totalBlocks++
			if f.Size > 0 {
				totalBlocks += int(f.Size / 1024)
			}
		}
		fmt.Printf("total %d\n", totalBlocks)
	}

	for _, file := range files {
		fmt.Print(FormatLong(file))
	}
}

// PrintRecursive outputs directory contents recursively (-R flag).
// Parameters:
//   - path: root directory path to start recursion
//   - showHidden: whether to include dotfiles
//   - longFormat: whether to use -l format for entries
//
// Returns:
//   - error: non-nil on filesystem errors
//
// Scope: Recursively traverses directories, printing each with a header.
// Uses depth-first traversal with sorted output at each level.
// No return value on success - writes directly to stdout.
func PrintRecursive(path string, showHidden bool, longFormat bool) error {
	files, err := fs.ReadDir(path, showHidden)
	if err != nil {
		return err
	}

	fs.SortFiles(files, false, false)

	fmt.Printf("%s:\n", path)

	if longFormat {
		PrintLong(files, true)
	} else {
		PrintStandard(files)
	}

	for _, file := range files {
		if file.IsDir && file.Name != "." && file.Name != ".." {
			fmt.Println()
			PrintRecursive(file.Path, showHidden, longFormat)
		}
	}

	return nil
}

// getColorizedName wraps the filename in ANSI escape sequence codes depending on type.
// getColorizedName wraps the filename in ANSI escape sequence codes depending on type.
func getColorizedName(name string, mode uint32) string {
	const (
		reset = "\033[0m"
		blue  = "\033[1;34m" // Bold Blue
		green = "\033[1;32m" // Bold Green
	)

	fileType := mode & 0o170000

	switch fileType {
	case 0o040000: // S_IFDIR (Directory)
		return fmt.Sprintf("%s%s%s", blue, name, reset)
	default:
		// Check for any executable permission bit (S_IXUSR, S_IXGRP, S_IXOTH)
		if mode&0o0111 != 0 {
			return fmt.Sprintf("%s%s%s", green, name, reset)
		}
	}
	return name
}