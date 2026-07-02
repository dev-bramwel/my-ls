package config

import "strings"

type Options struct {
	LongFormat   bool // -l
	ShowAll      bool // -a
	Reverse      bool // -r
	TimeSort     bool // -t
	Recursive    bool // -R
}

func ParseArgs(args []string) (Options, []string) {
	var opts Options
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
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

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	return opts, paths
}