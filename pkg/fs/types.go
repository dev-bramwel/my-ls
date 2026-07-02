package fs

import (
	"time"
)

type FileInfo struct {
	Name          string
	Path          string
	IsDir         bool
	IsSymlink     bool
	SymlinkTarget string
	Size          int64
	ModTime       time.Time
	ModeString    string
	Mode          uint32
	LinkCount     uint64
	Owner         string
	Group         string
	Blocks        int64
}