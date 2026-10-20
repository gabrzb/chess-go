package match

import (
	"context"
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

var ErrInvalidAIConfig = errors.New("AI requires a context, game, move chooser and valid color")

type startAIMessage struct{}

type aiResultMessage struct {
	requestID uint64
	sessionID uint64
	position  chess.Position
	move      chess.Move
	err       error
}

func NewWithAI(ctx context.Context, g *game.Game, chooser ai.MoveChooser, color chess.Color, session uint64) (Model, error) {
	if ctx == nil || g == nil || chooser == nil || (color != chess.White && color != chess.Black) {
		return Model{}, ErrInvalidAIConfig
	}

	m := New(g, session)
	m.aiContext = ctx
	m.opponent = chooser
	m.aiColor = color

	if color == chess.White {
		m.cursor = chess.A7
	}

	return m, nil
}

func (m Model) isAITurn() bool {
	return m.opponent != nil && m.game.Position().SideToMove == m.aiColor
}

func (m *Model) StartAI() tea.Cmd {
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
	m.aiError = nil
	m.requestID++
	id, session, position, chooser := m.requestID, m.sessionID, m.game.Position(), m.opponent

	return func() tea.Msg {
		move, err := chooser.ChooseMove(ctx, position)

		return aiResultMessage{requestID: id, sessionID: session, position: position, move: move, err: err}
	}
}

func (m *Model) StopAI() {
	if m.cancelAI != nil {
		m.cancelAI()
		m.cancelAI = nil
	}

	m.thinking = false
}

func (m Model) receiveAIResult(result aiResultMessage) (tea.Model, tea.Cmd) {
	// A canceled search may finish after review closes and another search starts.
	if !m.thinking || result.requestID != m.requestID || result.sessionID != m.sessionID {
		return m, nil
	}

	m.StopAI()

	if result.position != m.game.Position() {
		m.aiError = errors.New("AI result discarded: the position changed")
		m.message = "AI result discarded: the position changed. Press Enter to retry."
		return m, nil
	}

	if result.err != nil {
		m.aiError = fmt.Errorf("AI failed: %w", result.err)
		m.message = "AI failed: " + result.err.Error() + ". Press Enter to retry."
		return m, nil
	}

	if err := m.playMove(result.move); err != nil {
		m.aiError = fmt.Errorf("AI failed: %w", err)
		m.message = "AI failed: " + err.Error() + ". Press Enter to retry."
		return m, nil
	}

	m.cursor = result.move.To
	return m, nil
}

func IsAIResult(msg tea.Msg, session uint64) bool {
	result, ok := msg.(aiResultMessage)
	return ok && result.sessionID == session
}
