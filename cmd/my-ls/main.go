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
	// This ensures that multiple arguments match the global execution sequence of standard ls
	sortPaths(paths, opts.Reverse)

	// Track if we're processing multiple paths (affects output formatting)
	multiplePaths := len(paths) > 1

	// Process each path provided on the command line
	for i, path := range paths {
		// Check if the path exists and whether it's a file or directory
		// fs.IsDirectory returns true if path is a directory, false otherwise
		// Error handling: prints error message and continues to next path
		isDir, err := fs.IsDirectory(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': No such file or directory\n", path)
			continue
		}

		// Print header for multiple paths or for directories with -l/-R flag
		// This matches ls behavior of printing "path:\n" before contents
		if multiplePaths || (isDir && (opts.LongFormat || opts.Recursive)) {
			fmt.Printf("%s:\n", path)
		}

		if opts.Recursive && isDir {
			// -R flag: recursively list all subdirectories
			// display.PrintRecursive handles the entire recursive traversal
			_ = display.PrintRecursive(path, opts.ShowAll, opts.LongFormat)
		} else if isDir {
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
		} else {
			// Path is a file - handle file output
			if opts.LongFormat {
				file, err := fs.ReadFile(path)
				if err == nil {
					display.PrintLong([]fs.FileInfo{*file}, false)
				}
			} else {
				fmt.Print(path + " ")
			}
		}

		// Print newline between multiple paths (after each path's output)
		if i < len(paths)-1 && multiplePaths {
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

// toLowerStr converts an ASCII string argument to lowercase for evaluation comparisons.
func toLowerStr(s string) string {
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}