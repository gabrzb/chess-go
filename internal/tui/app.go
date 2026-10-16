package tui

import (
	"context"
	"log"
	"os"

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
		file, err := os.OpenFile("chess-go.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		options.Logger = log.New(file, "", log.LstdFlags)
	}
	m, err := NewApplication(ctx, options)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m).Run()
	return err
}
