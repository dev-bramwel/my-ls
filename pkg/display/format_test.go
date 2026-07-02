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
		output := FormatLongWithPadding(file, 1, 4, 5, 3) // output string captures grid output formatting text

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
		output := FormatLongWithPadding(file, 1, 4, 4, 4)

		if output[0] != 'd' {
			t.Errorf("Expected directory descriptor output rows to begin with prefix 'd', got '%c'", output[0])
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