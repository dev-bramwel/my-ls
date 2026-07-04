package display

import (
	"my-ls/pkg/fs"
	"strings"
	"testing"
)

func TestFormatLongWithPadding(t *testing.T) {
	t.Run("formats regular file correctly with dynamic columns", func(t *testing.T) {
		file := fs.FileInfo{
			Name:      "test.txt",
			IsDir:     false,
			Size:      100,
			LinkCount: 1,
			Owner:     "user",
			Group:     "group",
			Mode:      0o100644, // S_IFREG character type bit mixed with rw-r--r-- rights
		}
		output := FormatLongWithPadding(file, 1, 4, 5, 3, 0, 0) // output string captures grid output formatting text

		if output == "" {
			t.Fatal("Expected populated text row, got empty response block")
		}
		if output[0] != '-' {
			t.Errorf("Expected regular file line tracking prefix to map '-', got '%c'", output[0])
		}
		if !strings.Contains(output, "user") || !strings.Contains(output, "group") || !strings.Contains(output, "100") {
			t.Errorf("Output data alignments missing required file attributes. Got: %q", output)
		}
	})

	t.Run("formats directory row with proper type marker", func(t *testing.T) {
		file := fs.FileInfo{
			Name:      "mydir",
			IsDir:     true,
			Size:      4096,
			LinkCount: 2,
			Owner:     "root",
			Group:     "root",
			Mode:      0o040755, // S_IFDIR character type bit mixed with rwxr-xr-x rights
		}
		output := FormatLongWithPadding(file, 1, 4, 4, 4, 0, 0)

		if output[0] != 'd' {
			t.Errorf("Expected directory descriptor output rows to begin with prefix 'd', got '%c'", output[0])
		}
	})

	t.Run("formats character device with major and minor instead of size", func(t *testing.T) {
		file := fs.FileInfo{
			Name:         "null",
			Size:         0,
			LinkCount:    1,
			Owner:        "root",
			Group:        "root",
			Mode:         0o020666,
			IsCharDevice: true,
			Major:        1,
			Minor:        3,
		}
		output := FormatLongWithPadding(file, 1, 4, 4, 6, 2, 3)

		if output[0] != 'c' {
			t.Errorf("Expected character device row to begin with prefix 'c', got '%c'", output[0])
		}
		if !strings.Contains(output, " 1,   3") {
			t.Errorf("Expected device major/minor in size column. Got: %q", output)
		}
	})

	t.Run("formats block device with major and minor instead of size", func(t *testing.T) {
		file := fs.FileInfo{
			Name:          "loop0",
			Size:          0,
			LinkCount:     1,
			Owner:         "root",
			Group:         "root",
			Mode:          0o060660,
			IsBlockDevice: true,
			Major:         7,
			Minor:         0,
		}
		output := FormatLongWithPadding(file, 1, 4, 4, 6, 2, 3)

		if output[0] != 'b' {
			t.Errorf("Expected block device row to begin with prefix 'b', got '%c'", output[0])
		}
		if !strings.Contains(output, " 7,   0") {
			t.Errorf("Expected device major/minor in size column. Got: %q", output)
		}
	})

	t.Run("appends ACL marker to permission string", func(t *testing.T) {
		file := fs.FileInfo{
			Name:      "secure",
			Size:      1,
			LinkCount: 1,
			Owner:     "root",
			Group:     "root",
			Mode:      0o100600,
			ACLMarker: "+",
		}
		output := FormatLongWithPadding(file, 1, 4, 4, 1, 0, 0)

		if !strings.HasPrefix(output, "-rw-------+") {
			t.Errorf("Expected permission string to include ACL marker. Got: %q", output)
		}
	})
}

func TestGetColorizedName(t *testing.T) {
	t.Run("wraps directories in bold blue escape markers", func(t *testing.T) {
		name := "src"
		result := GetColorizedName(name, 0o040000) // 0o040000 represents the S_IFDIR constant directory mask
		if !strings.HasPrefix(result, Blue) || !strings.HasSuffix(result, Reset) {
			t.Errorf("Color mapping wrappers failing terminal color verification tests. Got: %q", result)
		}
	})

	t.Run("wraps executable binaries in bold green escape markers", func(t *testing.T) {
		name := "my-ls"
		result := GetColorizedName(name, 0o100755) // Standard active execution permissions mask
		if !strings.HasPrefix(result, Green) || !strings.HasSuffix(result, Reset) {
			t.Errorf("Execution tracking properties failing to colorize correctly. Got: %q", result)
		}
	})
}
