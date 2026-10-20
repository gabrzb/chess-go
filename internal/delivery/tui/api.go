package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/delivery/tui/application"
	"github.com/gabrzb/chess-go/internal/delivery/tui/match"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

type Options = application.Options

var ErrInvalidAIConfig = match.ErrInvalidAIConfig

func NewApplication(ctx context.Context, options Options) (tea.Model, error) {
	return application.New(ctx, options)
}

func NewModel(g *game.Game) tea.Model { return match.New(g, 0) }

func NewModelWithAI(ctx context.Context, g *game.Game, chooser chess.MoveChooser, color chess.Color) (tea.Model, error) {
	m, err := match.NewWithAI(ctx, g, chooser, color, 0)
	if err != nil {
		return nil, err
	}

	return m, nil
}
