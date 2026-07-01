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

to create a test environment, you can Copy and paste this script directly into your terminal. It sets up files, nested subdirectories, a dash-named folder, device lookups, synchronized timestamps, and symbolic links exactly as the questionnaire mandates:

Bash

# 1. Base Structure & Files
mkdir -p test_ls/dir1/subdir1 test_ls/dir2
touch test_ls/file1.txt test_ls/file2.txt test_ls/dir1/subdir1/nested.txt

# 2. Folder with '-' as a name
mkdir -p test_ls/-

# 3. Symbolic Links (Files and Directories)
ln -sf test_ls/file1.txt test_ls/sym_file
ln -sf test_ls/dir1 test_ls/sym_dir

# 4. Force uniform timestamps for the matching modification time test
find test_ls -exec touch -t 202607011200 {} +

Once testing is complete, you can drop the temporary environment instantly using:

# Bash
rm -rf test_ls