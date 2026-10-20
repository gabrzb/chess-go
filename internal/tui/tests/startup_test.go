package tui_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestLoggingDoesNotPreventStartup(t *testing.T) {
	var cacheVariable string
	switch runtime.GOOS {
	case "windows":
		cacheVariable = "LocalAppData"
	case "linux":
		cacheVariable = "XDG_CACHE_HOME"
	default:
		t.Skip("requires a platform with a configurable user cache directory")
	}
	dir := t.TempDir()
	t.Setenv(cacheVariable, dir)
	// Invalid options stop before opening a terminal, after logger initialization.
	options := tui.Options{Difficulty: ai.Difficulty("invalid")}
	if err := tui.RunWithOptions(options); !errors.Is(err, ai.ErrInvalidDifficulty) {
		t.Fatal("startup failed before validating options", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "chess-go", "chess-go.log")); err != nil {
		t.Fatal("startup did not create its diagnostic log", err)
	}
	blocked := filepath.Join(dir, "file")
	if err := os.WriteFile(blocked, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(cacheVariable, blocked)
	if err := tui.RunWithOptions(options); !errors.Is(err, ai.ErrInvalidDifficulty) {
		t.Fatal("logging failure prevented option validation", err)
	}
}
