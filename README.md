my-ls

Technical User Guide

Developers: mamani & mumutugi
Language: Go (Golang)

Environment: Unix / Linux Environment Alignment

Date: July 2026

my-ls Project Documentation 1

Section 1: User Guide (README)

Project Overview
my-ls is a production-grade, POSIX-compliant implementation of the core storage-listing utility ls
engineered entirely in Go. Built from the ground up without relying on high-level diagnostic frameworks, the
software replicates the behavior, output format, spacing configurations, colorization patterns, and error
reporting matrices of the standard GNU/Linux ls engine.

Installation & Build Workflow
The project features a automated pipeline driven by a structured Makefile . To ensure clean execution
environments, developers should rigidly adhere to the following command workflow:
1. Building the Binary
To compile the codebase into a standalone executable at the repository root, run:

$ make build

This invokes the Go compiler toolchain and outputs an optimized ./my-ls binary block ready for native
terminal execution.
2. Interacting with the Binary
In compliance with rigid alignment guidelines, do not use make run to execute the app. Instead, run the
natively compiled target file directly from the terminal, supplying any desired flags or folder paths:

$ ./my-ls -l test_ls/
$ ./my-ls -laR

3. Cleanup Routine
To clear temporary artifacts, stale binaries, and garbage objects generated during test cycles, always execute
the clean target immediately following your evaluation loop:

$ make clean

my-ls Project Documentation 2

Interactive Usage & Supported Flags
The application seamlessly parses individual and aggregated flag string clusters to dynamically adjust
directory traversal and visual presentations:
-l : Triggers the long format, printing file type, permission strings, link counts, owners, groups, sizes,
timestamps, and target paths.
-a : Includes hidden entries (files and folders beginning with a dot . ), alongside the implicit .
and .. context pointers.
-r : Inverts the active sorting sequence (reversing alphanumeric or chronological orderings).
-t : Sorts entries chronologically by modification timestamps (newest first) instead of alphabetically.
-R : Activates deep recursive traversal down the entire sub-directory matrix, appending distinct headers
for every node.

Command Examples & Verification

# Example A: Standard short-form space-separated output with implicit current directory
$ ./my-ls
pkg main.go Makefile test_ls
# Example B: Long listing of a specific directory target with block count totals
$ ./my-ls -l test_ls/sym_dir/
total 4
drwxrwxr-x 2 mamani mamani 4096 Jul 1 14:01 subdir1