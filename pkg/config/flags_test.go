package config

import (
	// testing - provides testing framework
	"testing"
)

// TestParseArgs tests the flag parsing functionality.
func TestParseArgs(t *testing.T) {
	t.Run("parses no flags - defaults to current directory", func(t *testing.T) {
		opts, paths := ParseArgs([]string{})
		if opts.LongFormat {
			t.Error("Expected LongFormat to be false")
		}
		if opts.ShowAll {
			t.Error("Expected ShowAll to be false")
		}
		if opts.Reverse {
			t.Error("Expected Reverse to be false")
		}
		if opts.TimeSort {
			t.Error("Expected TimeSort to be false")
		}
		if opts.Recursive {
			t.Error("Expected Recursive to be false")
		}
		if len(paths) != 1 || paths[0] != "." {
			t.Errorf("Expected paths=['.'], got %v", paths)
		}
	})

	t.Run("parses -l flag", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-l"})
		if !opts.LongFormat {
			t.Error("Expected LongFormat to be true")
		}
	})

	t.Run("parses -a flag", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-a"})
		if !opts.ShowAll {
			t.Error("Expected ShowAll to be true")
		}
	})

	t.Run("parses -r flag", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-r"})
		if !opts.Reverse {
			t.Error("Expected Reverse to be true")
		}
	})

	t.Run("parses -t flag", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-t"})
		if !opts.TimeSort {
			t.Error("Expected TimeSort to be true")
		}
	})

	t.Run("parses -R flag", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-R"})
		if !opts.Recursive {
			t.Error("Expected Recursive to be true")
		}
	})

	t.Run("parses combined flags -la", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-la"})
		if !opts.LongFormat || !opts.ShowAll {
			t.Error("Expected LongFormat and ShowAll to be true")
		}
	})

	t.Run("parses path argument", func(t *testing.T) {
		_, paths := ParseArgs([]string{"pkg"})
		if len(paths) != 1 || paths[0] != "pkg" {
			t.Errorf("Expected paths=['pkg'], got %v", paths)
		}
	})

	t.Run("parses multiple paths", func(t *testing.T) {
		_, paths := ParseArgs([]string{"pkg", "cmd"})
		if len(paths) != 2 {
			t.Errorf("Expected 2 paths, got %d", len(paths))
		}
	})

	t.Run("parses flags before paths", func(t *testing.T) {
		opts, paths := ParseArgs([]string{"-l", "pkg"})
		if !opts.LongFormat {
			t.Error("Expected LongFormat to be true")
		}
		if len(paths) != 1 || paths[0] != "pkg" {
			t.Errorf("Expected paths=['pkg'], got %v", paths)
		}
	})

	t.Run("parses paths before flags", func(t *testing.T) {
		opts, paths := ParseArgs([]string{"pkg", "-l"})
		if !opts.LongFormat {
			t.Error("Expected LongFormat to be true")
		}
		if len(paths) != 1 || paths[0] != "pkg" {
			t.Errorf("Expected paths=['pkg'], got %v", paths)
		}
	})
}