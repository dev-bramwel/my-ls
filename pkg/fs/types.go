package fs

import "time"

type FileInfo struct {
	Name        string
	Path        string
	IsDir       bool
	Size        int64
	ModTime     time.Time
	ModeString  string // e.g., "-rw-r--r--"
	LinkCount   uint64 // Needed for -l
	Owner       string // Needed for -l
	Group       string // Needed for -l
}
