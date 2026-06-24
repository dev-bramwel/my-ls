package main

import (
	"fmt"
	"os"
	"my-ls/pkg/config"
)

func main() {
	// Skip the first argument which is the binary name itself
	opts, paths := config.ParseArgs(os.Args[1:])

	fmt.Printf("Paths to process: %v\n", paths)
	fmt.Printf("Active Flags: Long(-l): %t, All(-a): %t, Rev(-r): %t, Time(-t): %t, Rec(-R): %t\n",
		opts.LongFormat, opts.ShowAll, opts.Reverse, opts.TimeSort, opts.Recursive)
}
