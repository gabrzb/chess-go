package tui

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
)

func Run() error {
	return RunWithOptions(Options{Seed: 1})
}

func RunWithAI(difficulty ai.Difficulty, human chess.Color, seed int64) error {
	if human != chess.White && human != chess.Black {
		return chess.ErrInvalidColor
	}

	return RunWithOptions(Options{Difficulty: difficulty, HumanColor: human, Seed: seed, StartGame: true})
}

func RunWithOptions(options Options) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if options.Logger == nil {
		dir, _ := os.UserCacheDir()
		logger, closeLog := defaultLogger(dir)
		defer closeLog()
		options.Logger = logger
	}

	m, err := NewApplication(ctx, options)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m).Run()

	return err
}

func defaultLogger(cacheDir string) (*log.Logger, func()) {
	logger := log.New(io.Discard, "", log.LstdFlags)
	noop := func() {}

	if cacheDir == "" {
		return logger, noop
	}

	dir := filepath.Join(cacheDir, "chess-go")

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return logger, noop
	}

	file, err := os.OpenFile(filepath.Join(dir, "chess-go.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return logger, noop
	}

	logger.SetOutput(file)
	return logger, func() { _ = file.Close() }
}
