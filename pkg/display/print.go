package display

import (
	// fmt - provides output formatting functions
	"fmt"
	// my-ls/pkg/fs - provides FileInfo type for formatting
	"my-ls/pkg/fs"
	// strconv - provides integer to string conversion
	"strconv"
	// time - provides time formatting utilities
	
)

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
		fmt.Print(coloredName + "  ")
	}
	fmt.Print("\n")
}

// PrintLong outputs files in long format (-l flag) with dynamic column padding.
// Parameters:
//   - files: slice of fs.FileInfo to display in detail
//   - showTotal: when true, prints "total X" line at start
//
// Returns:
//   - outputs formatted detail lines to stdout
//
// Scope: Scans all files to determine maximum width profiles for links, owner, group, and size.
// Iterates over files and calls FormatLongWithPadding for each.
// Prints "total X" line at start when showTotal is true.
// No return value - writes directly to stdout.
func PrintLong(files []fs.FileInfo, showTotal bool) {
	if showTotal {
		var totalBlocks int64
		// Sum up only the pre-calculated system blocks of the listed entries[cite: 4]
		for _, f := range files {
			totalBlocks += f.Blocks
		}
		// Linux kernel tracks in 512B units; standard ls displays in 1024B blocks
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Dynamic padding tracking variables
	maxLinks := 0
	maxOwner := 0
	maxGroup := 0
	maxSize := 0

	// Scan through entries to extract highest character counts for column metrics
	for _, f := range files {
		if len(strconv.FormatUint(f.LinkCount, 10)) > maxLinks {
			maxLinks = len(strconv.FormatUint(f.LinkCount, 10))
		}
		if len(f.Owner) > maxOwner {
			maxOwner = len(f.Owner)
		}
		if len(f.Group) > maxGroup {
			maxGroup = len(f.Group)
		}
		if len(strconv.FormatInt(f.Size, 10)) > maxSize {
			maxSize = len(strconv.FormatInt(f.Size, 10))
		}
	}

	for _, file := range files {
		fmt.Print(FormatLongWithPadding(file, maxLinks, maxOwner, maxGroup, maxSize))
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