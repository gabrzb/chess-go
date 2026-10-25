package match

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/application"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

var ErrInvalidAIConfig = application.ErrInvalidAIConfig

type startAIMessage struct{}

func NewWithAI(ctx context.Context, g *game.Game, chooser chess.MoveChooser, color chess.Color, session uint64) (Model, error) {
	opponent, err := application.NewOpponent(ctx, g, chooser, color, session)
	if err != nil {
		return Model{}, err
	}

	m := New(g)
	m.opponent = opponent

	if color == chess.White {
		m.cursor = chess.A7
	}

	return m, nil
}

func (m Model) isAITurn() bool {
	return m.opponent.IsTurn()
}

func (m *Model) StartAI() tea.Cmd {
	if m.reviewing || len(m.promotionMoves) > 0 {
		return nil
	}

	search, err := m.opponent.Start()
	if err != nil {
		m.message = err.Error()
		return nil
	}

	if search == nil {
		return nil
	}

	m.message = ""
	m.aiError = nil
	return func() tea.Msg {
		timer := time.NewTimer(500 * time.Millisecond)
		defer timer.Stop()
		result := search()
		<-timer.C
		return result
	}
}

func (m *Model) StopAI() {
	m.opponent.Stop()
}

func (m Model) receiveAIResult(result application.SearchResult) (tea.Model, tea.Cmd) {
	accepted, err := m.opponent.Accept(result)
	if !accepted {
		return m, nil
	}

	if errors.Is(err, application.ErrPositionChanged) {
		m.aiError = err
		m.message = "AI result discarded: the position changed. Press Enter to retry."
		return m, nil
	}

	if err != nil {
		m.aiError = fmt.Errorf("AI failed: %w", err)
		m.message = "AI failed: " + err.Error() + ". Press Enter to retry."
		return m, nil
	}

	m.ClearSelection()
	m.cursor = result.Move.To
	return m, nil
}

func IsAIResult(msg tea.Msg, session uint64) bool {
	result, ok := msg.(application.SearchResult)
	return ok && result.Session == session
}
