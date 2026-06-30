package fs

import (
	// testing - provides testing framework
	"testing"
	// time - provides time utilities for test data
	"time"
)

// TestReadDir tests the ReadDir function for various scenarios.
func TestReadDir(t *testing.T) {
	t.Run("reads current directory", func(t *testing.T) {
		files, err := ReadDir(".", false)
		if err != nil {
			t.Fatalf("ReadDir failed: %v", err)
		}
		if len(files) == 0 {
			t.Error("Expected files in current directory, got none")
		}
	})

	t.Run("hides dotfiles by default", func(t *testing.T) {
		files, _ := ReadDir(".", false)
		for _, f := range files {
			if f.Name[0] == '.' {
				t.Errorf("Unexpected dotfile found: %s", f.Name)
			}
		}
	})

	t.Run("shows dotfiles with showHidden=true", func(t *testing.T) {
		files, _ := ReadDir(".", true)
		foundDot := false
		for _, f := range files {
			if f.Name == ".." || f.Name == "." {
				foundDot = true
				break
			}
		}
		if !foundDot {
			t.Error("Expected . and .. entries when showHidden=true")
		}
	})

	t.Run("Populates all metadata fields", func(t *testing.T) {
		files, _ := ReadDir(".", true)
		for _, f := range files {
			if f.Mode == 0 {
				t.Errorf("File %s has Mode=0, expected non-zero", f.Name)
			}
			if f.LinkCount == 0 {
				t.Errorf("File %s has LinkCount=0, expected non-zero", f.Name)
			}
		}
	})
}

// TestIsDirectory tests the IsDirectory function.
func TestIsDirectory(t *testing.T) {
	t.Run("returns true for directory", func(t *testing.T) {
		isDir, err := IsDirectory(".")
		if err != nil {
			t.Fatalf("IsDirectory failed: %v", err)
		}
		if !isDir {
			t.Error("Expected '.' to be a directory")
		}
	})

	t.Run("returns error for non-existent path", func(t *testing.T) {
		_, err := IsDirectory("nonexistent_file_12345")
		if err == nil {
			t.Error("Expected error for non-existent path")
		}
	})
}

// TestSortFiles tests the sorting functionality.
func TestSortFiles(t *testing.T) {
	t.Run("sorts alphabetically ascending by default", func(t *testing.T) {
		files := []FileInfo{
			{Name: "zebra", ModTime: time.Now()},
			{Name: "apple", ModTime: time.Now()},
			{Name: "mango", ModTime: time.Now()},
		}
		SortFiles(files, false, false)

		if files[0].Name != "apple" {
			t.Errorf("Expected 'apple' first, got '%s'", files[0].Name)
		}
		if files[1].Name != "mango" {
			t.Errorf("Expected 'mango' second, got '%s'", files[1].Name)
		}
		if files[2].Name != "zebra" {
			t.Errorf("Expected 'zebra' third, got '%s'", files[2].Name)
		}
	})

	t.Run("sorts alphabetically descending with reverse flag", func(t *testing.T) {
		files := []FileInfo{
			{Name: "zebra", ModTime: time.Now()},
			{Name: "apple", ModTime: time.Now()},
			{Name: "mango", ModTime: time.Now()},
		}
		SortFiles(files, false, true)

		if files[0].Name != "zebra" {
			t.Errorf("Expected 'zebra' first, got '%s'", files[0].Name)
		}
		if files[1].Name != "mango" {
			t.Errorf("Expected 'mango' second, got '%s'", files[1].Name)
		}
		if files[2].Name != "apple" {
			t.Errorf("Expected 'apple' third, got '%s'", files[2].Name)
		}
	})

	t.Run("sorts by time descending with timeSort flag", func(t *testing.T) {
		now := time.Now()
		oldTime := now.Add(-24 * time.Hour)
		newTime := now.Add(24 * time.Hour)
		midTime := now

		files := []FileInfo{
			{Name: "old", ModTime: oldTime},
			{Name: "new", ModTime: newTime},
			{Name: "mid", ModTime: midTime},
		}
		SortFiles(files, true, false)

		if files[0].Name != "new" {
			t.Errorf("Expected 'new' first (newest), got '%s'", files[0].Name)
		}
		if files[1].Name != "mid" {
			t.Errorf("Expected 'mid' second, got '%s'", files[1].Name)
		}
		if files[2].Name != "old" {
			t.Errorf("Expected 'old' third (oldest), got '%s'", files[2].Name)
		}
	})
}

// TestReadFile tests the ReadFile function.
func TestReadFile(t *testing.T) {
	t.Run("reads file metadata", func(t *testing.T) {
		file, err := ReadFile("go.mod")
		if err != nil {
			// Skip test if file doesn't exist (might be running from different directory)
			t.Skip("go.mod not found in current directory")
		}
		if file == nil {
			t.Fatal("Expected non-nil FileInfo")
		}
		if file.Name != "go.mod" {
			t.Errorf("Expected name 'go.mod', got '%s'", file.Name)
		}
		if file.Size == 0 {
			t.Error("Expected non-zero size for go.mod")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := ReadFile("nonexistent_file_12345")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}