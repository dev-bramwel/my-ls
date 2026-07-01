# Technical Architecture & Documentation

## 1. Architectural Overview & Folder Tree

The system utilizes a modular, decoupled layout splitting flag evaluation, filesystem ingestion, sorting matrices, and formatting operations into highly cohesive packages.

```text
my-ls-1/
├── main.go               # Application entrypoint and pipeline controller
├── Makefile              # Build automation matrix (build, clean targets)
└── pkg/
    ├── config/
    │   └── flags.go      # Flag tokenization and Option parsing structures
    ├── display/
    │   ├── formatter.go  # Long-format row generators, permissions, dates
    │   └── print.go      # Layout dispatchers (Standard, Long, Recursive)
    └── fs/
        ├── reader.go     # Low-level system file descriptors and directory reads
        ├── sorter.go     # In-place multi-criteria selection sorting engine
        └── types.go      # Primary global FileInfo metadata models
Package Breakdown
Main Controller (main.go): Coordinates the workflow. It processes command arguments via the config engine, performs an initial case-insensitive pass on target path parameters, filters files from directory streams, and feeds clean metadata rows down into the display layer.

Configuration Core (pkg/config/flags.go): Exposes the ParseArgs parser, scanning terminal arguments safely without throwing exceptions on combined character strings (e.g., -laR).

Filesystem Engine (pkg/fs/): Maps low-level operating system values directly to an internal FileInfo representation. It extracts system block allocations (syscall.Stat_t.Blocks), dynamic numerical ownership indices mapped to system user names via os/user lookup tables, and calculates explicit link counts.

Sorting Engine (pkg/fs/sorter.go): Implements a highly predictable, stable selection sorting architecture to guarantee consistent sorting results across Linux distributions, stripping leading dots out of case-insensitive strings to align hidden components correctly.

Display & Presenter Core (pkg/display/): Maps permission bitmasks to 10-character notation strings (e.g., lrwxrwxrwx), tracks terminal column widths dynamically for uniform right-aligned padding blocks, and feeds color escapes to standard output.

2. Collaboration & Work Division Strategy
The project was successfully designed and refined by a two-person engineering team, consisting of mumutugi and mamani. Responsibilities were split strategically based on foundational initialization versus fine-grained behavioral alignment.

Engineer	Primary Areas of Responsibility	Key Contributions & Milestone Metrics
mumutugi	Architecture Setup & Basic File Streaming	
- Initialized directory structure and base .gitignore definitions.


- Designed the primary directory ingestion scanners and initial Go file reads.


- Configured standard system ANSI code definitions for color mapping.


- Managed base integration branches and initial repository health workflows.

mamani	Flag Ingestion, Sorting Logic, and System Alignment	
- Built compound flag string parser configurations (-l, -a, -r, -t, -R).


- Developed case-insensitive alphabetical sorting and chronological sorting fallbacks.


- Engineered dynamic spacing calculators for column padding.


- Identified and refactored advanced symlink edge cases, error matching, and trailing newlines.

3. Engineering Challenges, Troubleshooting & Insights
During the intensive debugging and system-alignment phases, several significant technical issues were uncovered and resolved.

A. Symlink Traversal Ambiguity: os.Stat vs. os.Lstat
The Challenge: When evaluating target directories and symbolic links via arguments, passing a link pointing to a folder without a trailing slash (e.g., test_ls/sym_dir) caused the program to evaluate it as a standard directory, incorrectly generating block totals and traversing inside the directory.

The Insight: Go's standard os.Stat() function implicitly follows symbolic links to their end targets. To isolate and describe the link structure itself, the codebase was updated to utilize os.Lstat() across all entry evaluations. This ensures that a raw symlink argument correctly registers as a link file, while the addition of a trailing slash (e.g., sym_dir/) triggers natural operating system resolution down to the underlying folder node.

B. Block Count Totalization Overrides for Single File Arguments
The Challenge: When executing long listings on individual symbolic links or static text documents, the program initially output a total X header row above the file description block. System ls guidelines require block summaries to be completely omitted when listing explicit file parameters.

The Insight: The execution sequence in main.go was split into distinct phases. Files are grouped and processed together first, calling the formatting arrays with the block total flag explicitly disabled (showTotal = false). The block total header calculation is strictly bound to multi-entry directory lists.

C. Preservation of Full Terminal Input Path Arguments
The Challenge: When target files were queried directly, the application stripped their leading paths and displayed only the base filename (e.g., printing sym_dir instead of test_ls/sym_dir), breaking strict validation rules.

The Insight: The ReadFile metadata compiler was adjusted to explicitly bind the full, unmutated user input parameter path down into the display structure rather than querying info.Name(), matching native terminal expectations perfectly.

D. Dynamic Error Message and Signal String Mirroring
The Challenge: Appending a directory slash to a regular text file path string caused the program to throw a flat "No such file or directory" error, missing the system error string standard "Not a directory".

The Insight: Hardcoded custom text blocks were removed from the error reporting loop. By extracting and type-asserting native kernel errors down to *os.PathError wrappers, the program reads and logs the actual system message string (pathErr.Err.Error()), instantly matching the operating system output.