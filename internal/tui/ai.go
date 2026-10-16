package tui

import (
	"context"
	"errors"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

var ErrInvalidAIConfig = errors.New("AI requires a context, game, move chooser and valid color")

type startAIMessage struct{}

type aiResultMessage struct {
	requestID uint64
	position  chess.Position
	move      chess.Move
	err       error
}

func NewModelWithAI(ctx context.Context, g *game.Game, chooser ai.MoveChooser, color chess.Color) (tea.Model, error) {
	if ctx == nil || g == nil || chooser == nil || (color != chess.White && color != chess.Black) {
		return nil, ErrInvalidAIConfig
	}
	m := NewModel(g).(model)
	m.aiContext = ctx
	m.opponent = chooser
	m.aiColor = color
	return m, nil
}

func (m model) isAITurn() bool {
	return m.opponent != nil && m.game.Position().SideToMove == m.aiColor
}

func (m *model) startAITurn() tea.Cmd {
	if !m.isAITurn() || m.thinking || m.reviewing || len(m.promotionMoves) > 0 {
		return nil
	}
	snapshot, err := m.game.Snapshot()
	if err != nil {
		m.message = err.Error()
		return nil
	}
	if snapshot.Status != chess.StatusOngoing {
		return nil
	}
	ctx, cancel := context.WithCancel(m.aiContext)
	m.cancelAI = cancel
	m.thinking = true
	m.message = ""
	m.requestID++
	id, position, chooser := m.requestID, m.game.Position(), m.opponent
	return func() tea.Msg {
		move, err := chooser.ChooseMove(ctx, position)
		return aiResultMessage{requestID: id, position: position, move: move, err: err}
	}
}

func (m *model) stopAITurn() {
	if m.cancelAI != nil {
		m.cancelAI()
		m.cancelAI = nil
	}
	m.thinking = false
}

func (m model) receiveAIResult(result aiResultMessage) (tea.Model, tea.Cmd) {
	// A canceled search may finish after review closes and another search starts.
	if !m.thinking || result.requestID != m.requestID {
		return m, nil
	}
	m.stopAITurn()
	if result.position != m.game.Position() {
		m.message = "AI result discarded: the position changed. Press Enter to retry."
		return m, nil
	}
	if result.err != nil {
		m.message = "AI failed: " + result.err.Error() + ". Press Enter to retry."
		return m, nil
	}
	if err := m.playMove(result.move); err != nil {
		m.message = "AI failed: " + err.Error() + ". Press Enter to retry."
		return m, nil
	}
	m.cursor = result.move.To
	return m, nil
}
