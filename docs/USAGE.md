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

## Building

```bash
make build
# or
go build ./cmd/my-ls
```

## Testing

```bash
make test
# or
go test -v ./...
```