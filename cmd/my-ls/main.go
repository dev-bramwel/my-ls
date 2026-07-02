package main

import (
	"fmt"
	"os"
	"my-ls/pkg/config"
	"my-ls/pkg/display"
	"my-ls/pkg/fs"
)

func main() {
	opts, paths := config.ParseArgs(os.Args[1:])

	sortPaths(paths, opts.Reverse)

	var filesOnly []fs.FileInfo
	var dirsOnly []string

	for _, path := range paths {
		isDir, err := fs.IsDirectory(path)
		if err != nil {
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

	hasFiles := len(filesOnly) > 0
	hasDirs := len(dirsOnly) > 0
	multipleDirs := len(dirsOnly) > 1 || (hasFiles && hasDirs)

	if hasFiles {
		fs.SortFiles(filesOnly, opts.TimeSort, opts.Reverse)

		if opts.LongFormat {
			display.PrintLong(filesOnly, false)
		} else {
			display.PrintStandard(filesOnly)
		}

		if hasDirs {
			fmt.Print("\n")
		}
	}

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
				display.PrintLong(files, true)
			} else {
				display.PrintStandard(files)
			}
		}

		if i < len(dirsOnly)-1 {
			fmt.Print("\n")
		}
	}
}

func sortPaths(paths []string, reverse bool) {
	n := len(paths)
	for i := 0; i < n-1; i++ {
		extreme := i
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