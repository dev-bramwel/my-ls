## Project Structure

```text
my-ls-1/
├── cmd/
│   └── my-ls/
│       └── main.go         # Application entry point (keeps main clean)
├── pkg/
│   ├── config/
│   │   └── flags.go        # Flag parsing logic (-l, -a, -r, -t, -R)
│   ├── display/
│   │   ├── formatter.go    # Long-format (-l) styling, padding, colors
│   │   └── print.go        # Standard output logic
│   └── fs/
│       ├── reader.go       # os.ReadDir wrappers and error handling
│       ├── types.go        # Custom structs for file metadata
│       └── sorter.go       # Alphabetical, time (-t), and reverse (-r) sorting
├── docs/
│   ├── ARCHITECTURE.md     # System design and component layout
│   └── USAGE.md            # Examples of how to run and test flags
├── Makefile                # Build, test, and cleanup automation
├── go.mod                  # Go module file
└── README.md               # Main project overview
```