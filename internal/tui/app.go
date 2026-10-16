package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func Run() error {
	program := tea.NewProgram(NewModel(game.NewGame()))
	_, err := program.Run()
	return err
}

func RunWithAI(difficulty ai.Difficulty, human chess.Color, seed int64) error {
	if human != chess.White && human != chess.Black {
		return chess.ErrInvalidColor
	}
	player, err := ai.NewPlayer(difficulty, seed)
	if err != nil {
		return err
	}
	color := chess.Black
	if human == chess.Black {
		color = chess.White
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := NewModelWithAI(ctx, game.NewGame(), player, color)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m).Run()
	return err
}
