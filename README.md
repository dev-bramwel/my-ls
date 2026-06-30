# MY-LS-1
A lightweight, high-performance clone of the Unix `ls` command written from scratch in Go. This project explores low-level Unix system programming, file system traversals, custom sorting layouts, and explicit flag evaluation without relying on high-level orchestration packages like `os/exec`.

## Objectives
The primary goal of this project is to dive deep into Go's fundamentals and Unix system interfaces. It showcases:
* Low-level Unix filesystem interaction.
* Custom argument/flag evaluation using Go slice mechanics.
* Clean package decoupling and data pipeline separation.

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

## Features
* **Standard Listing:** Lists file and directory entries alphabetically. Defaults to the current directory (`.`) if no target path is supplied.
* **Core Flag Configurations:**
  * `-l` : Long listing layout detailing file types, permissions, hard link count, owner, group, file byte size, and last modified timestamp.
  * `-a` : Forces visibility of hidden entities (files or folders starting with a `.`).
  * `-r` : Inverts the active sorting sequence.
  * `-t` : Priorities sorting by modification time (newest entries first).
* **Bonus Elements:**
  * `-R` : Enables recursive Depth-First Search traversal throughout the directory hierarchy.

## Getting Started
### Installation & Compilation
Build the standalone executable utility locally using the project automation rules via your terminal configuration.

```bash
make build
```

### Running the Test Engine
Execute all decoupled unit testing suites directly across localized packages using your project configuration tools.

```bash
make test
```
