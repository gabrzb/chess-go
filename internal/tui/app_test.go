package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	dir := t.TempDir()
	logger, closeLog := defaultLogger(dir)
	logger.Print("diagnostic")
	closeLog()
	data, err := os.ReadFile(filepath.Join(dir, "chess-go", "chess-go.log"))
	if err != nil || !strings.Contains(string(data), "diagnostic") {
		t.Fatalf("log was not persisted: %q, %v", data, err)
	}

	blocked := filepath.Join(dir, "file")
	if err := os.WriteFile(blocked, nil, 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"", blocked} {
		logger, closeLog := defaultLogger(path)
		logger.Print("logging failure must not prevent playing")
		closeLog()
	}
}
