# my-ls
## Technical User Guide

**Developers:** [Moses Amani](https://learn.zone01kisumu.ke/git/mamani) & [Bramwel Mutugi](https://learn.zone01kisumu.ke/git/mumutugi)  
**Language:** Go (Golang)  
**Environment:** Unix / Linux Environment Alignment  
**Date:** July 2026  

---

### Project Overview
`my-ls` is a production-grade, POSIX-compliant implementation of the core storage-listing utility `ls` engineered entirely in Go. Built from the ground up without relying on high-level diagnostic frameworks, the software replicates the behavior, output format, spacing configurations, colorization patterns, and error reporting matrices of the standard GNU/Linux `ls` engine.

---

### Installation & Build Workflow
The project features an automated pipeline driven by a structured `Makefile`. To ensure clean execution environments, developers should rigidly adhere to the following command workflow:

#### 1. Building the Binary
To compile the codebase into a standalone executable at the repository root, run:
```bash
$ make build
```
This invokes the Go compiler toolchain and outputs an optimized ./my-ls binary block ready for native terminal execution.

#### 2. Interacting with the Binary
In compliance with rigid alignment guidelines, do not use make run to execute the app. Instead, run the natively compiled target file directly from the terminal, supplying any desired flags or folder paths:

```bash
$ ./my-ls -l test_ls/
$ ./my-ls -laR
```
#### 3. Cleanup Routine
To clear temporary artifacts, stale binaries, and garbage objects generated during test cycles, always execute the clean target immediately following your evaluation loop:
```bash
$ make clean
```
### Interactive Usage & Supported Flags
The application seamlessly parses individual and aggregated flag string clusters to dynamically adjust directory traversal and visual presentations:
```text
-l : Triggers the long format, printing file type, permission strings, link counts, owners, groups, sizes, timestamps, and target paths.

-a : Includes hidden entries (files and folders beginning with a dot .), alongside the implicit . and .. context pointers.

-r : Inverts the active sorting sequence (reversing alphanumeric or chronological orderings).

-t : Sorts entries chronologically by modification timestamps (newest first) instead of alphabetically.

-R : Activates deep recursive traversal down the entire sub-directory matrix, appending distinct headers for every node.
```
### Test Directory Structure & Audit Usage
The repository includes a dedicated test_ls/ suite structured specifically to validate low-level filesystem corner cases during system audits. This folder contains specialized file types designed to verify that your flags and formatting match the system ls exactly.

#### Key Testing Scenarios Provided:

- Directory Symlink with Trailing Slash (test_ls/sym_dir/): Tests if the binary correctly resolves the target directory (dir1) and displays its interior contents along with a directory block summary (total X).

- Directory Symlink without Trailing Slash (test_ls/sym_dir): Verifies that the program lists the symbolic link itself on a single line showing its structural arrow reference (-> dir1), while completely omitting the total block line.

- File Symlink with Forced Slash (test_ls/sym_file/): Validates system error string reproduction. Passing a trailing slash on a regular file symlink triggers a native system override, throwing a Not a directory error instead of a generic file-missing notice.
#### Command Examples & Verification
##### Example A: Standard short-form space-separated output with implicit current directory
```bash
$ ./my-ls
```
pkg  main.go  Makefile  test_ls

##### Example B: Long listing of a specific directory target with block count totals
```bash
$ ./my-ls -l test_ls/sym_dir/
```
total 4
drwxrwxr-x 2 mamani mamani 4096 Jul  1 14:01 subdir1

##### Example C: Long listing of a raw symbolic link asset without a trailing slash
```bash
$ ./my-ls -l test_ls/sym_dir
```
lrwxrwxrwx 1 mamani mamani 4 Jul  1 14:01 test_ls/sym_dir -> dir1