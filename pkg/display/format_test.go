package display

import (
	// testing - provides testing framework
	"testing"
	// my-ls/pkg/fs - provides FileInfo type
	"my-ls/pkg/fs"
)

// TestFormatLong tests the FormatLong function.
func TestFormatLong(t *testing.T) {
	t.Run("formats file correctly", func(t *testing.T) {
		file := fs.FileInfo{
			Name:      "test.txt",
			IsDir:     false,
			Size:      100,
			LinkCount: 1,
			Owner:     "user",
			Group:     "group",
			Mode:      0o644, // rw-r--r--
		}
		output := FormatLong(file)

		if output == "" {
			t.Error("Expected non-empty output")
		}
		// Check that it starts with "-" for regular file
		if output[0] != '-' {
			t.Errorf("Expected output to start with '-', got '%c'", output[0])
		}
	})

	t.Run("formats directory correctly", func(t *testing.T) {
		file := fs.FileInfo{
			Name:      "mydir",
			IsDir:     true,
			Size:      0,
			LinkCount: 2,
			Owner:     "root",
			Group:     "root",
			Mode:      0o755, // rwxr-xr-x
		}
		output := FormatLong(file)

		if output[0] != 'd' {
			t.Errorf("Expected output to start with 'd', got '%c'", output[0])
		}
	})
}

// TestFormatPermissions tests the permission string formatting.
func TestFormatPermissions(t *testing.T) {
	t.Run("read-only permissions", func(t *testing.T) {
		// 0o444 = r--r--r--
		perms := formatPermissions(false, 0o444)
		expected := "-r--r--r--"
		if perms != expected {
			t.Errorf("Expected '%s', got '%s'", expected, perms)
		}
	})

	t.Run("full permissions", func(t *testing.T) {
		// 0o777 = rwxrwxrwx
		perms := formatPermissions(true, 0o777)
		expected := "drwxrwxrwx"
		if perms != expected {
			t.Errorf("Expected '%s', got '%s'", expected, perms)
		}
	})

	t.Run("setuid bit", func(t *testing.T) {
		// 0o4755 = rwsr-xr-x
		perms := formatPermissions(true, 0o4755)
		expected := "drwsr-xr-x"
		if perms != expected {
			t.Errorf("Expected '%s', got '%s'", expected, perms)
		}
	})
}