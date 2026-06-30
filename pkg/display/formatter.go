// Package display provides output formatting functions for my-ls.
// This file contains the long format (-l) output generation.
// The functions here create output matching the standard ls -l command format.
package display

import (
	"fmt"
	//"os"
	"my-ls/pkg/fs"
)

// This file intentionally left minimal.
// Long format functionality is in print.go to keep all output functions together.

const (
	Reset  = "\033[0m"
	Blue   = "\033[1;34m" // For Directories
	Green  = "\033[1;32m" // For Executables
	Cyan   = "\033[1;36m" // For Symlinks (Optional bonus styling)
)

// FormatName applies ANSI color codes based on the file type and permissions.
func FormatName(file fs.FileInfo) string {
	// 1. Check if it's a directory
	if file.IsDir {
		return fmt.Sprintf("%s%s%s", Blue, file.Name, Reset)
	}

	// 2. Check if the file is executable (looks at owner, group, or other execute bits)
	// os.ModePerm is a bitmask for all runtime permissions (0777)
	const executeMask = 0111 // Looks for any 'x' bit (--x--x--x)
	if file.Mode&executeMask != 0 {
		return fmt.Sprintf("%s%s%s", Green, file.Name, Reset)
	}

	// 3. Default plain output for standard files
	return file.Name
}
