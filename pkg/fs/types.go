package fs

import (
	// time - provides Time type for modification timestamp storage
	"time"
)

type FileInfo struct {
	Name        string
	Path        string
	IsDir       bool
	IsSymlink   bool      // True if this is a symbolic link
	SymlinkTarget string  // Target path for symlinks (for -> display)
	Size        int64
	ModTime     time.Time
	ModeString  string  // e.g., "crw-r--r--" for character devices
	Mode        uint32  // Full mode bits (file type + permissions)
	LinkCount   uint64  // Needed for -l
	Owner       string  // Needed for -l
	Group       string  // Needed for -l
	Blocks      int64  // Add this field to track raw 512B system blocks
}

