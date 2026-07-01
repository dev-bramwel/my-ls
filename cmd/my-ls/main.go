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
	// os.Args contains the command-line arguments, where os.Args[0] is the program name[cite: 2]
	// We skip the first argument to get only the user-provided flags and paths[cite: 2]
	// config.ParseArgs returns parsed Options struct and slice of target paths[cite: 2]
	opts, paths := config.ParseArgs(os.Args[1:]) //[cite: 2]

	// Sort the top-level target paths case-insensitively before iterating[cite: 2]
	sortPaths(paths, opts.Reverse) //[cite: 2]

	var filesOnly []fs.FileInfo
	var dirsOnly []string

	// First pass: Validate paths and segregate files from directories[cite: 2]
	for _, path := range paths {
		// Check if the path exists and whether it's a file or directory[cite: 2]
		isDir, err := fs.IsDirectory(path) //[cite: 2]
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': No such file or directory\n", path) //[cite: 2]
			continue
		}

		if isDir {
			dirsOnly = append(dirsOnly, path) //[cite: 2]
		} else {
			file, err := fs.ReadFile(path) //[cite: 2]
			if err == nil {
				filesOnly = append(filesOnly, *file) //[cite: 2]
			}
		}
	}

	// Track if we have a mixed output or multiple directory headers to print[cite: 2]
	hasFiles := len(filesOnly) > 0 //[cite: 2]
	hasDirs := len(dirsOnly) > 0   //[cite: 2]
	multipleDirs := len(dirsOnly) > 1 || (hasFiles && hasDirs) //[cite: 2]

	// Print all standalone files first on a single line (unless using -l)[cite: 2]
	if hasFiles {
		// Sort regular file targets according to sorting flags[cite: 2]
		fs.SortFiles(filesOnly, opts.TimeSort, opts.Reverse) //[cite: 2]

		if opts.LongFormat {
			display.PrintLong(filesOnly, false) //[cite: 2]
		} else {
			// Utilize standard display logic to safely print inline file names with proper colorization[cite: 2]
			display.PrintStandard(filesOnly) //[cite: 2]
		}

		// Add an extra newline separation if directories follow the files[cite: 2]
		if hasDirs {
			fmt.Print("\n") //[cite: 2]
		}
	}

	// Second pass: Process each directory target cleanly[cite: 2]
	for i, path := range dirsOnly {
		// Print header for directories if there are multiple targets[cite: 2]
		if multipleDirs {
			fmt.Printf("%s:\n", path) //[cite: 2]
		}

		if opts.Recursive {
			// -R flag: recursively list all subdirectories[cite: 2]
			// display.PrintRecursive handles the entire recursive traversal[cite: 2]
			_ = display.PrintRecursive(path, opts.ShowAll, opts.LongFormat, opts.TimeSort, opts.Reverse)
		} else {
			// Path is a directory - read contents[cite: 2]
			// fs.ReadDir returns FileInfo slice for directory entries[cite: 2]
			// showHidden controls whether dotfiles (starting with '.') are included[cite: 2]
			files, err := fs.ReadDir(path, opts.ShowAll) //[cite: 2]
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: %s\n", err.Error()) //[cite: 2]
				continue
			}

			// fs.SortFiles sorts the FileInfo slice in place[cite: 2]
			// timeSort=true for -t flag (sort by modification time)[cite: 2]
			// reverse=true for -r flag (reverse the sort order)[cite: 2]
			fs.SortFiles(files, opts.TimeSort, opts.Reverse) //[cite: 2]

			// Choose output format based on flags[cite: 2]
			// -l flag: use long format with details[cite: 2]
			// default: use simple space-separated format[cite: 2]
			// showTotal=true for directories (print total blocks line)[cite: 2]
			if opts.LongFormat {
				display.PrintLong(files, true) //[cite: 2]
			} else {
				display.PrintStandard(files) //[cite: 2]
			}
		}

		// Print newline between multiple directory segments[cite: 2]
		if i < len(dirsOnly)-1 {
			fmt.Print("\n") //[cite: 2]
		}
	}
}

// sortPaths handles sorting the top-level path strings case-insensitively using selection sort.[cite: 2]
func sortPaths(paths []string, reverse bool) {
	n := len(paths) //[cite: 2]
	for i := 0; i < n-1; i++ {
		extreme := i //[cite: 2]
		for j := i + 1; j < n; j++ {
			nameJ := toLowerStr(paths[j])     //[cite: 2]
			nameExt := toLowerStr(paths[extreme]) //[cite: 2]

			if !reverse {
				if nameJ < nameExt || (nameJ == nameExt && paths[j] < paths[extreme]) { //[cite: 2]
					extreme = j //[cite: 2]
				}
			} else {
				if nameJ > nameExt || (nameJ == nameExt && paths[j] > paths[extreme]) { //[cite: 2]
					extreme = j //[cite: 2]
				}
			}
		}
		if extreme != i {
			paths[i], paths[extreme] = paths[extreme], paths[i] //[cite: 2]
		}
	}
}

// toLowerStr converts an ASCII string argument to lowercase for evaluation comparisons, ignoring leading dots.[cite: 2]
func toLowerStr(s string) string {
	if len(s) > 1 && s[0] == '.' && s != ".." { //[cite: 2]
		s = s[1:] //[cite: 2]
	}
	b := []byte(s) //[cite: 2]
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' { //[cite: 2]
			b[i] += 32 //[cite: 2]
		}
	}
	return string(b) //[cite: 2]
}