package main

import (
	// fmt - provides formatted I/O functions
	"fmt"
	// os - provides command line arguments and directory functions
	"os"
	// my-ls/pkg/config - provides flag parsing functionality
	"my-ls/pkg/config"
	// my-ls/pkg/display - provides output formatting functionality
	"my-ls/pkg/display"
	// my-ls/pkg/fs - provides file reading and sorting functionality
	"my-ls/pkg/fs"
)

func main() {
	// os.Args contains the command-line arguments, where os.Args[0] is the program name
	// We skip the first argument to get only the user-provided flags and paths
	// config.ParseArgs returns parsed Options struct and slice of target paths
	opts, paths := config.ParseArgs(os.Args[1:])

	// Sort the top-level target paths case-insensitively before iterating
	sortPaths(paths, opts.Reverse)

	var filesOnly []fs.FileInfo
	var dirsOnly []string

	// First pass: Validate paths and segregate files from directories
	for _, path := range paths {
		// Check if the path exists and whether it's a file or directory
		isDir, err := fs.IsDirectory(path)
		if err != nil {
			// AUDIT FIX: Use err.Error() or unwrap the path error to extract the exact OS message
			// This dynamically outputs "Not a directory" or "No such file or directory" to match system ls.
			if pathErr, ok := err.(*os.PathError); ok {
				fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %s\n", path, pathErr.Err.Error())
			} else {
				fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %s\n", path, err.Error())
			}
			continue
		}

		if isDir {
			dirsOnly = append(dirsOnly, path)
		} else {
			file, err := fs.ReadFile(path)
			if err == nil {
				filesOnly = append(filesOnly, *file)
			}
		}
	}

	// Track if we have a mixed output or multiple directory headers to print
	hasFiles := len(filesOnly) > 0
	hasDirs := len(dirsOnly) > 0
	multipleDirs := len(dirsOnly) > 1 || (hasFiles && hasDirs)

	// Print all standalone files first on a single line (unless using -l)
	if hasFiles {
		// Sort regular file targets according to sorting flags
		fs.SortFiles(filesOnly, opts.TimeSort, opts.Reverse)

		if opts.LongFormat {
			display.PrintLong(filesOnly, false)
		} else {
			// Utilize standard display logic to safely print inline file names with proper colorization
			display.PrintStandard(filesOnly)
		}

		// Add an extra newline separation if directories follow the files
		if hasDirs {
			fmt.Print("\n")
		}
	}

	// Second pass: Process each directory target cleanly
	for i, path := range dirsOnly {
		// Print header for directories if there are multiple targets and we are NOT in recursive mode
		// (PrintRecursive handles its own headers automatically)
		if multipleDirs && !opts.Recursive {
			fmt.Printf("%s:\n", path)
		}

		if opts.Recursive {
			// -R flag: recursively list all subdirectories
			// display.PrintRecursive handles the entire recursive traversal
			_ = display.PrintRecursive(path, opts.ShowAll, opts.LongFormat, opts.TimeSort, opts.Reverse)
		} else {
			// Path is a directory - read contents
			// fs.ReadDir returns FileInfo slice for directory entries
			// showHidden controls whether dotfiles (starting with '.') are included
			files, err := fs.ReadDir(path, opts.ShowAll)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: %s\n", err.Error())
				continue
			}

			// fs.SortFiles sorts the FileInfo slice in place
			// timeSort=true for -t flag (sort by modification time)
			// reverse=true for -r flag (reverse the sort order)
			fs.SortFiles(files, opts.TimeSort, opts.Reverse)

			// Choose output format based on flags
			// -l flag: use long format with details
			// default: use simple space-separated format
			// showTotal=true for directories (print total blocks line)
			if opts.LongFormat {
				display.PrintLong(files, true)
			} else {
				display.PrintStandard(files)
			}
		}

		// Print newline between multiple directory segments
		if i < len(dirsOnly)-1 {
			fmt.Print("\n")
		}
	}
}

// sortPaths handles sorting the top-level path strings case-insensitively using selection sort.
func sortPaths(paths []string, reverse bool) {
	n := len(paths)
	for i := 0; i < n-1; i++ {
		extreme := i
		for j := i + 1; j < n; j++ {
			nameJ := toLowerStr(paths[j])
			nameExt := toLowerStr(paths[extreme])

			if !reverse {
				if nameJ < nameExt || (nameJ == nameExt && paths[j] < paths[extreme]) {
					extreme = j
				}
			} else {
				if nameJ > nameExt || (nameJ == nameExt && paths[j] > paths[extreme]) {
					extreme = j
				}
			}
		}
		if extreme != i {
			paths[i], paths[extreme] = paths[extreme], paths[i]
		}
	}
}

// toLowerStr converts an ASCII string argument to lowercase for evaluation comparisons, ignoring leading dots.
func toLowerStr(s string) string {
	if len(s) > 1 && s[0] == '.' && s != ".." {
		s = s[1:]
	}
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}