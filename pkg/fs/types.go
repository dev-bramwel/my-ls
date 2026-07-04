package fs

import (
	"time"
)

// FileInfo models normalized low-level filesystem property dimensions used across layout displays
type FileInfo struct {
	Name          string    // Base display title string identifying the file target
	Path          string    // Complete filepath tracking coordinate line mapping metrics
	IsDir         bool      // Active truth flag determining directory verification contexts
	IsSymlink     bool      // Active truth flag tracking symlink classification
	IsCharDevice  bool      // Active truth flag tracking Unix character device descriptors
	IsBlockDevice bool      // Active truth flag tracking Unix block device descriptors
	ACLMarker     string    // Permission suffix marker for ACL/security metadata
	SymlinkTarget string    // String destination target text path pointing to target files
	Size          int64     // File allocation byte counter volume matching storage limits
	ModTime       time.Time // Hardware modification time entry mapping parameters
	ModeString    string    // Readable permission string format layout matching legacy wrappers
	Mode          uint32    // Numeric type mask and flag bit tracking properties
	LinkCount     uint64    // Hard link parameter reference counter metrics
	Owner         string    // Username tracking metric identifying administrative creators
	Group         string    // Group name validation mapping identifier tracking profile paths
	Blocks        int64     // Raw physical system 512-byte partition block count metrics
	Rdev          uint64    // Raw device identifier used to print major/minor values for device nodes
	Major         uint64    // Device major number extracted from Rdev
	Minor         uint64    // Device minor number extracted from Rdev
}
