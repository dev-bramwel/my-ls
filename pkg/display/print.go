package display

import (
	"fmt"
	"my-ls/pkg/fs"
	"strconv"
)

// PrintStandard formats filenames horizontally across space-separated layouts
func PrintStandard(files []fs.FileInfo) {
	if len(files) == 0 {
		return
	}

	for _, file := range files {
		coloredName := GetColorizedName(file.Name, file.Mode)
		fmt.Print(coloredName + "  ")
	}
	fmt.Print("\n")
}

// PrintLong parses columns inside target slices to isolate cell character sizes for padding grid generation
func PrintLong(files []fs.FileInfo, showTotal bool) {
	if showTotal {
		var totalBlocks int64 // totalBlocks aggregates raw filesystem architecture allocations
		for _, f := range files {
			totalBlocks += f.Blocks
		}
		// System ls normalizes 512-byte blocks down to standard 1024-byte chunks, requiring us to divide by 2
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Padding constraint max widths used to keep textual column layouts cleanly aligned
	maxLinks := 0
	maxOwner := 0
	maxGroup := 0
	maxSize := 0

	for _, f := range files {
		if len(strconv.FormatUint(f.LinkCount, 10)) > maxLinks {
			maxLinks = len(strconv.FormatUint(f.LinkCount, 10))
		}
		if len(f.Owner) > maxOwner {
			maxOwner = len(f.Owner)
		}
		if len(f.Group) > maxGroup {
			maxGroup = len(f.Group)
		}
		size := formatSizeOrDevice(f)
		if len(size) > maxSize {
			maxSize = len(size)
		}
	}

	for _, file := range files {
		fmt.Print(FormatLongWithPadding(file, maxLinks, maxOwner, maxGroup, maxSize))
	}
}

// PrintRecursive deep-dives down sub-elements within a folder tree map sequence
func PrintRecursive(path string, showHidden bool, longFormat bool, timeSort bool, reverse bool) error {
	files, err := fs.ReadDir(path, showHidden)
	if err != nil {
		return err
	}

	fs.SortFiles(files, timeSort, reverse)

	fmt.Printf("%s:\n", path)

	if longFormat {
		PrintLong(files, true)
	} else {
		PrintStandard(files)
	}

	// Recurse straight down matching directory structures while ignoring relative pointers ('.' and '..')
	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir && file.Name != "." && file.Name != ".." {
			fmt.Println() // Prints separating empty spacing row between recursive tracking lists
			PrintRecursive(file.Path, showHidden, longFormat, timeSort, reverse)
		}
	}

	return nil
}

// GetColorizedName queries type bit flags to apply appropriate escape styling parameters
func GetColorizedName(name string, mode uint32) string {
	fileType := mode & 0o170000 // fileType stores isolated bit masks characterizing system descriptor shapes

	switch fileType {
	case 0o040000: // S_IFDIR: Directory entry structure
		return fmt.Sprintf("%s%s%s", Blue, name, Reset)
	case 0o120000: // S_IFLNK: Symbolic reference line target
		return fmt.Sprintf("%s%s%s", Cyan, name, Reset)
	case 0o020000, 0o060000: // S_IFCHR or S_IFBLK: Device node entry points
		return fmt.Sprintf("%s%s%s", DeviceOpts, name, Reset)
	default:
		// Execute permission mask check matching user, group, or global execution targets
		if mode&0o0111 != 0 {
			return fmt.Sprintf("%s%s%s", Green, name, Reset)
		}
	}
	return name
}
