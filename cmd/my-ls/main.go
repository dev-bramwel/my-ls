package main

import (
	"fmt"
	"os"
	"my-ls/pkg/config"
	"my-ls/pkg/display"
	"my-ls/pkg/fs"
)

func main() {
	// opts stores the boolean states of our flags (-l, -a, -r, -t, -R)
	// paths stores an array of target strings indicating files or directories to list (defaults to ".")
	opts, paths := config.ParseArgs(os.Args[1:])

	// Sort the target path parameters case-insensitively before reading any file metadata.
	sortPaths(paths, opts.Reverse)

	// filesOnly stores items explicitly passed via CLI that are regular files or symlinks to files
	var filesOnly []fs.FileInfo
	// dirsOnly stores paths passed via CLI that point to valid directory descriptors
	var dirsOnly []string

	// Segregate file paths from directory paths to mirror system layout ordering
	for _, path := range paths {
		isDir, err := fs.IsDirectory(path)
		if err != nil {
			// If path validation fails, unwrap the OS error details to log identical system text
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

	// Boolean tracking indicators used to format block spacing layout cleanly
	hasFiles := len(filesOnly) > 0
	hasDirs := len(dirsOnly) > 0
	// multipleDirs tracks whether directory path headers (e.g., "dir:") should print out
	multipleDirs := len(dirsOnly) > 1 || (hasFiles && hasDirs)

	// Phase 1: Output standalone individual files first
	if hasFiles {
		fs.SortFiles(filesOnly, opts.TimeSort, opts.Reverse)

		if opts.LongFormat {
			display.PrintLong(filesOnly, false) // false omits the "total X" block header row
		} else {
			display.PrintStandard(filesOnly)
		}

		if hasDirs {
			fmt.Print("\n") // Clean separator before starting directory printing block
		}
	}

	// Phase 2: Traverse and output directory contents sequentially
	for i, path := range dirsOnly {
		if multipleDirs && !opts.Recursive {
			fmt.Printf("%s:\n", path)
		}

		if opts.Recursive {
			_ = display.PrintRecursive(path, opts.ShowAll, opts.LongFormat, opts.TimeSort, opts.Reverse)
		} else {
			files, err := fs.ReadDir(path, opts.ShowAll)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: %s\n", err.Error())
				continue
			}

			fs.SortFiles(files, opts.TimeSort, opts.Reverse)

			if opts.LongFormat {
				display.PrintLong(files, true) // true prints the summary filesystem "total X" block row
			} else {
				display.PrintStandard(files)
			}
		}

		// Print newline spacing if there are more directories remaining in the list
		if i < len(dirsOnly)-1 {
			fmt.Print("\n")
		}
	}
}

// sortPaths uses an in-place selection sort on raw target parameters before file read evaluation
func sortPaths(paths []string, reverse bool) {
	n := len(paths) // n stores the length of the path slice
	for i := 0; i < n-1; i++ {
		extreme := i // extreme stores the index of the highest/lowest sorting element found
		for j := i + 1; j < n; j++ {
			nameJ := fs.ToLower(paths[j])
			nameExt := fs.ToLower(paths[extreme])

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