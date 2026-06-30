# my-ls Implementation

A Go implementation of the `ls` command with support for the following flags:
- `-l`: Long format output (permissions, links, owner, group, size, date, name)
- `-R`: Recursive directory listing
- `-a`: Include hidden files (starting with '.')
- `-r`: Reverse sort order
- `-t`: Sort by modification time

## Package Structure

```
my-ls-1/
├── cmd/my-ls/main.go      # Application entry point
├── pkg/config/flags.go    # Flag parsing logic
├── pkg/display/print.go   # Output formatting (-l format, standard output, recursive)
├── pkg/fs/
│   ├── reader.go          # Directory reading and file metadata extraction
│   ├── sorter.go          # Manual sorting implementation (selection sort)
│   └── types.go           # FileInfo struct definition
└── go.mod                 # Go module definition
```

## Allowed Packages Used

- `fmt`: Formatted I/O operations
- `os`: File system operations (ReadDir, Stat, Lstat)
- `os/user`: User/group name resolution from UIDs/GIDs
- `strconv`: Integer to string conversion
- `strings`: String manipulation (hidden file detection)
- `syscall`: Permission bit constants (via octal literals)
- `time`: Time formatting for modification dates
- `errors`: Error creation

## Error Handling

- Directory read errors are returned and displayed to stderr
- Non-existent paths produce appropriate error messages
- Failed stat operations skip entries but don't halt execution
- User/group lookup failures fall back to numeric IDs

## Execution & Interaction Rules

This section describes how to run and test flag permutations inside the application framework. 

### Core Binary Compilations
To compile the standalone binary executable tool asset directly into your workspace project root workspace directory, utilize the primary automation rule:
```bash
make build
```

---

## Foundational Interactions

### Standard Workspace Traversal
Executing the binary directly without passing trailing arguments or location descriptors defaults to parsing the local present directory environment, listing items alphabetically:
```bash
./my-ls
```

### Target Vector Specifications
You can pass single or multiple unique absolute or relative directory paths as arguments to scan distant targets independently:
```bash
./my-ls /usr/bin /var/log
```

---

## Composite Flag Configurations

The parsing matrix supports combined, stacked POSIX shorthand arguments seamlessly.

### Chronological Reversed Listings (-ltr)
Evaluates target file systems long-format properties, sorting the results chronologically by last modified date in reverse order. This orientation pins your most recently updated assets to the very bottom of the terminal output block:
```bash
./my-ls -ltr
```

### Deep Recursive Inspection (-Ra)
Instructs the execution loop to initiate deep Depth-First Search recursive loops down nested subdirectory trees while showing hidden system files starting with a dot configuration prefix:
```bash
./my-ls -Ra pkg/
```

---

## Output Visual Styling Guide

When writing text streams to standard terminal outputs, the display formatting layer applies distinct ANSI color highlights based on resource type definitions:

* **Blue Text Highlight:** Identifies an explicit File System Directory.
* **Green Text Highlight:** Identifies an Executable Script or Binary Application asset.
* **Default Standard Text:** Signifies baseline operational files (.txt, .go, .json, etc.).