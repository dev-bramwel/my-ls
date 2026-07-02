package config

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	t.Run("parses no flags - defaults to current directory", func(t *testing.T) {
		opts, paths := ParseArgs([]string{})
		if opts.LongFormat || opts.ShowAll || opts.Reverse || opts.TimeSort || opts.Recursive {
			t.Error("Expected all options to evaluate false under default parameters")
		}
		if len(paths) != 1 || paths[0] != "." {
			t.Errorf("Expected default paths=['.'], got %v", paths)
		}
	})

	t.Run("parses flags correctly", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-l", "-a", "-r", "-t", "-R"})
		if !opts.LongFormat || !opts.ShowAll || !opts.Reverse || !opts.TimeSort || !opts.Recursive {
			t.Error("Expected all flags to map explicitly to true values")
		}
	})

	t.Run("parses combined flags -la", func(t *testing.T) {
		opts, _ := ParseArgs([]string{"-la"})
		if !opts.LongFormat || !opts.ShowAll {
			t.Error("Expected combined parsing fields to activate option states cleanly")
		}
	})

	t.Run("parses path argument", func(t *testing.T) {
		_, paths := ParseArgs([]string{"pkg", "cmd"})
		if len(paths) != 2 || paths[0] != "pkg" || paths[1] != "cmd" {
			t.Errorf("Expected path string collection to persist, got %v", paths)
		}
	})
}