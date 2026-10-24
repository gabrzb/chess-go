package application

import (
	"fmt"
	"strings"

	"github.com/gabrzb/chess-go/internal/delivery/tui/render"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

func (m application) playerInfo(s game.Snapshot, compact bool) string {
	state := m.match.State()
	var status string

	if m.screen == finalBoardScreen {
		title, body, _ := m.resultText()
		paragraphs := strings.Split(body, "\n\n")
		status = "GAME OVER\n" + title + "\n" + paragraphs[len(paragraphs)-1]
	} else {
		status = "Turn: " + sideText(s.SideToMove)

		if m.config.computer {
			if compact {
				status += " · You: " + sideText(opposite(state.AIColor))
			} else {
				status += "\nYou play as " + sideText(opposite(state.AIColor)) + "."

				if s.SideToMove == state.AIColor {
					status += "\nComputer · " + difficultyLabel(m.config.level)
				}
			}
		} else if !compact {
			status += "\nPlayer " + map[chess.Color]string{chess.White: "1", chess.Black: "2"}[s.SideToMove]
		}

		if state.Thinking {
			status += "\n" + render.Thinking(m.frame, !m.noAnimations)
		}
	}

	if state.Reviewing {
		status += fmt.Sprintf("\nREVIEW · position %d/%d", state.HistoryIndex+1, m.match.Game().HistoryLen())
	}

	return status
}

func (m application) gameFooter() string {
	state := m.match.State()

	if state.Reviewing {
		return "←/A previous · →/D next · H/Esc return · F flip · ? help"
	}

	if m.screen == finalBoardScreen {
		return "F flip · H history · Esc return to result"
	}

	if state.Selected.Valid() {
		piece, _ := m.match.Game().Position().Board.PieceAt(state.Selected)
		name := strings.ToUpper(render.PieceName(piece.Type))

		return name + " · " + state.Selected.String() + " — choose a destination.\nEnter confirm · Esc cancel selection · ? help"
	}

	if state.Thinking {
		return "AI thinking · Esc menu · H history · ? help · Ctrl+C quit"
	}

	return "Arrows/WASD move · Enter select\nEsc menu · F flip · H history · ? help"
}

func (m application) boardContent() string {
	s, err := m.match.DisplayedSnapshot()
	if err != nil {
		return "Could not load the position."
	}

	state := m.match.State()
	mode, compactMode := "LOCAL · PLAYER VS PLAYER", "LOCAL"

	if m.config.computer {
		mode, compactMode = "PLAYER VS AI · "+difficultyLabel(m.config.level), "VS AI · "+difficultyLabel(m.config.level)
	}

	lastMove := "—"
	if s.LastMove != nil {
		lastMove = s.LastMove.From.String() + " → " + s.LastMove.To.String()
	}

	end := len(m.labels)
	if state.Reviewing {
		end = min(end, state.HistoryIndex)
	}

	var notices []string
	if !state.Reviewing && m.match.Game().ClaimableDraw() != game.NoReason {
		notices = append(notices, "DRAW AVAILABLE · Esc opens the menu to claim.")
	}

	if s.Check && m.screen != finalBoardScreen {
		warning := "The " + map[chess.Color]string{chess.White: "white", chess.Black: "black"}[s.SideToMove] + " king is in check."

		if m.config.computer && s.SideToMove != state.AIColor {
			warning = "Your king is in check. Get out of check."
		}

		notices = append(notices, "CHECK · "+strings.ToUpper(sideText(s.SideToMove))+" · "+warning)
	}
	if state.Message != "" {
		message := "Invalid move. Choose a highlighted square."

		if state.Message == chess.ErrEmptySource.Error() {
			message = "Select a piece belonging to the player whose turn it is."
		}

		notices = append(notices, message)
	}

	return render.Match(render.MatchPanel{
		Width: m.width, Height: m.height, Mode: mode, CompactMode: compactMode,
		Status: m.playerInfo(s, false), CompactStatus: m.playerInfo(s, true),
		LastMove: lastMove, Board: s.Board, BoardState: render.BoardState{Cursor: state.Cursor, Selected: state.Selected, LegalMoves: state.LegalMoves, Flipped: state.Flipped}, Footer: m.gameFooter(),
		Notice: render.Paint(strings.Join(notices, "\n"), render.Alert), Moves: m.labels[:end],
		Turn: render.TurnState{Color: s.SideToMove, Reviewing: state.Reviewing, Finished: s.Status != chess.StatusOngoing},
	})
}

func (m application) boardFits() bool {
	if m.match.Game() == nil {
		return false
	}

	content := m.boardContent()
	return render.Fits(content, m.width, m.height)
}
