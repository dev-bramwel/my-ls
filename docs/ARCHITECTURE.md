# Architecture

## Component Overview

### cmd/my-ls/main.go
Application entry point that:
1. Parses command-line arguments using config.ParseArgs
2. Processes each path (directory or file)
3. Handles -R flag for recursive traversal
4. For directories: calls fs.ReadDir and fs.SortFiles
5. For files with -l: calls fs.ReadFile for metadata
6. Outputs using display.PrintStandard or display.PrintLong

### pkg/config/flags.go
Flag parsing component:
- `Options` struct holds all flag states
- `ParseArgs()` function separates flags from path arguments
- Supports combined flags (e.g., `-la`, `-ltr`)
- Defaults to current directory (`.`) when no paths specified

### pkg/fs/reader.go
File system reading component:
- `ReadDir(path, showHidden)` returns FileInfo slice for directory contents
- `ReadFile(path)` returns FileInfo for a single file
- `IsDirectory(path)` checks if path is a directory
- `IsSymlink(path)` checks for symbolic links
- `getOwnership(info)` extracts owner/group from syscall.Stat_t

### pkg/fs/sorter.go
Manual selection sort implementation:
- Alphabetical sorting (default)
- Time-based sorting (-t flag)
- Reverse ordering (-r flag)
- No external sorting packages used (per requirements)

### pkg/fs/types.go
Defines `FileInfo` struct with:
- Name, Path: file identification
- IsDir: directory flag
- Size: file size in bytes
- ModTime: modification timestamp
- Mode: Unix permission bits
- LinkCount: hard link count
- Owner, Group: ownership information

### pkg/display/print.go
Output formatting component:
- `FormatLong(file)`: generates single `-l` format line
- `PrintStandard(files)`: space-separated output
- `PrintLong(files, showTotal)`: long format with optional total line


- `PrintRecursive(path, showHidden, longFormat)`: recursive directory traversal

## Data Flow

```
main() 
  -> config.ParseArgs(os.Args[1:])
  -> fs.ReadDir(path, showAll) OR fs.ReadFile(path)
  -> fs.SortFiles(files, timeSort, reverse)
  -> display.PrintStandard() OR display.PrintLong() OR display.PrintRecursive()
```

## Key Design Decisions

1. **No os/exec**: All file operations use direct os.ReadDir/Stat calls
2. **Manual sorting**: Selection sort algorithm since sort package is not allowed
3. **Permission handling**: Direct bit manipulation for permission string generation
4. **Error handling**: Graceful continuation on individual file errors
5. **ACL support**: Not fully implemented (system ls shows `+` suffix for ACL)