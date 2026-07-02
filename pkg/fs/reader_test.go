package fs

import (
	"testing"
	"time"
)

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
			if len(f.Name) > 0 && f.Name[0] == '.' {
				t.Errorf("Unexpected dotfile found: %s", f.Name)
			}
		}
	})

	t.Run("shows dotfiles with showHidden=true", func(t *testing.T) {
		files, _ := ReadDir(".", true)
		foundDot := false // foundDot evaluates true if relative folder links register correctly
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
}

func TestSortFiles(t *testing.T) {
	t.Run("sorts alphabetically ascending by default", func(t *testing.T) {
		files := []FileInfo{
			{Name: "cherry"},
			{Name: "banana"},
			{Name: "apple"},
		}
		SortFiles(files, false, false)

		if files[0].Name != "apple" { t.Errorf("Expected 'apple', got '%s'", files[0].Name) }
		if files[1].Name != "banana" { t.Errorf("Expected 'banana', got '%s'", files[1].Name) }
		if files[2].Name != "cherry" { t.Errorf("Expected 'cherry', got '%s'", files[2].Name) }
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

		if files[0].Name != "new" { t.Errorf("Expected 'new' first, got '%s'", files[0].Name) }
		if files[2].Name != "old" { t.Errorf("Expected 'old' last, got '%s'", files[2].Name) }
	})
}

func TestToLower(t *testing.T) {
	t.Run("converts to lowercase and strips single leading dot", func(t *testing.T) {
		// tests slice houses structured parameter matrices to run automated evaluation iterations
		tests := []struct {
			input    string // input parameter parameter string
			expected string // expected translation target value
		}{
			{"File.txt", "file.txt"},
			{".Hidden", "hidden"},
			{"..", ".."},
			{"A", "a"},
		}

		for _, test := range tests {
			result := ToLower(test.input) // result captures active computational output mappings
			if result != test.expected {
				t.Errorf("ToLower(%q) = %q; expected %q", test.input, result, test.expected)
			}
		}
	})
}