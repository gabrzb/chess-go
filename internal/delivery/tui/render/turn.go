package render

import (
	"charm.land/lipgloss/v2"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

const (
	compactTurnHeight = 6
	largeTurnHeight   = 8
)

type TurnState struct {
	Color               chess.Color
	Reviewing, Finished bool
}

func Turn(state TurnState, scale BoardScale) string {
	height := compactTurnHeight
	if scale == LargeBoard {
		height = largeTurnHeight
	}

	width := height * terminalCellAspect
	title, label := "TURN", "Player 1(W)"
	background, foreground := WhitePiece, BlackPiece

	if state.Color == chess.Black {
		label, background, foreground = "Player 2(B)", BlackPiece, WhitePiece
	}

	if state.Reviewing {
		title = "REVIEW"
	} else if state.Finished {
		title, label = "GAME", "OVER"
		background, foreground = Background, Foreground
	}

	block := lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color(foreground)).Background(lipgloss.Color(background)).Render(label)
	framed := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(Muted)).Render(block)

	return Center(title, lipgloss.Width(framed)) + "\n" + framed
}
