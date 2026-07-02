package display

import (
	"fmt"
	"my-ls/pkg/fs"
	"strconv"
)

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

func PrintLong(files []fs.FileInfo, showTotal bool) {
	if showTotal {
		var totalBlocks int64
		for _, f := range files {
			totalBlocks += f.Blocks
		}
		fmt.Printf("total %d\n", totalBlocks/2)
	}

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
		if len(strconv.FormatInt(f.Size, 10)) > maxSize {
			maxSize = len(strconv.FormatInt(f.Size, 10))
		}
	}

	for _, file := range files {
		fmt.Print(FormatLongWithPadding(file, maxLinks, maxOwner, maxGroup, maxSize))
	}
}

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

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir && file.Name != "." && file.Name != ".." {
			fmt.Println()
			PrintRecursive(file.Path, showHidden, longFormat, timeSort, reverse)
		}
	}

	return nil
}

func GetColorizedName(name string, mode uint32) string {
	fileType := mode & 0o170000

	switch fileType {
	case 0o040000: // S_IFDIR (Directory)
		return fmt.Sprintf("%s%s%s", Blue, name, Reset)
	case 0o120000: // S_IFLNK (Symbolic Link)
		return fmt.Sprintf("%s%s%s", Cyan, name, Reset)
	case 0o020000, 0o060000: // S_IFCHR or S_IFBLK (Devices)
		return fmt.Sprintf("%s%s%s", DeviceOpts, name, Reset)
	default:
		if mode&0o0111 != 0 {
			return fmt.Sprintf("%s%s%s", Green, name, Reset)
		}
	}
	return name
}