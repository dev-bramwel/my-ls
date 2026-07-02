package config

import "strings"

// Options encapsulates the execution states determined by active CLI flags
type Options struct {
	LongFormat   bool // -l: Shows metadata grid columns, ownership strings, and sizes
	ShowAll      bool // -a: Includes dotfiles like '.' and '..' in the directory stream
	Reverse      bool // -r: Flips the resulting sorting algorithm sequencing match output
	TimeSort     bool // -t: Sequences entries by filesystem modification time instead of names
	Recursive    bool // -R: Deep-dives through internal subdirectories down the directory tree
}

// ParseArgs scans terminal parameters to isolate flags from structural file path elements
func ParseArgs(args []string) (Options, []string) {
	var opts Options
	var paths []string // paths stores clean non-flag parameter targets collected from inputs

	for _, arg := range args {
		// Detect flags beginning with a single dash (ignoring empty dashes or single dashes like "-")
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Iterate across runes to process clustered strings natively (e.g., -laR)
			for _, char := range arg[1:] {
				switch char {
				case 'l':
					opts.LongFormat = true
				case 'a':
					opts.ShowAll = true
				case 'r':
					opts.Reverse = true
				case 't':
					opts.TimeSort = true
				case 'R':
					opts.Recursive = true
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}

	// Default to targeting the local executing directory if no explicit path parameters exist
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	return opts, paths
}