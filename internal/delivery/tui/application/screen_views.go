package application

import (
	"github.com/gabrzb/chess-go/internal/delivery/tui/render"

	tea "charm.land/bubbletea/v2"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

func (m application) card(title, body, footer string, double bool) string {
	width := render.ContentWidth(m.width, render.MenuWidth)

	return render.Card(title, body, footer, m.band(width, double), m.menuOptions(), m.choice, width)
}

func (m application) View() tea.View {
	var content string
	footer := "↑↓ navigate · Enter select"

	switch m.screen {
	case mainScreen:
		title := m.logoText()
		if m.width < render.LogoMinWidth || m.height < render.LogoMinHeight {
			title = render.Paint("CHESS", render.Accent)
		}

		setting := "on"

		if m.noAnimations {
			setting = "off"
		}

		content = m.card(title, "", footer+"\n? help · A animations: "+setting, true)
	case configScreen:
		content = m.card("NEW GAME\nPLAYER VS AI", "", "↑↓ navigate · ←→ change\nEnter confirm · Esc back", true)
	case playingScreen, finalBoardScreen:
		if m.match.State().PromotionPending && m.screen == playingScreen {
			content = render.Promotion(m.match.State().Cursor, m.match.Game().Position().SideToMove, m.choice, m.match.State().PromotionChoices, render.ContentWidth(m.width, render.PromotionWidth))
		} else {
			content = m.boardContent()
		}
	case pauseScreen:
		content = m.card("GAME MENU", "", footer+"\nEsc resume", false)
	case resignScreen:
		content = m.card("CONFIRM RESIGNATION?", sideText(opposite(m.resigningSide()))+" will win.", footer+"\nEsc back", false)
	case abandonScreen:
		content = m.card("END GAME?", "The current game will be discarded.", footer+"\nEsc back", false)
	case drawScreen:
		s, _ := m.match.Game().Snapshot()
		content = m.card("DRAW OFFER", sideText(s.DrawOffer)+" offered a draw.\nDoes "+sideText(opposite(s.DrawOffer))+" accept?", footer, false)
	case resultScreen:
		title, body, color := m.resultText()
		content = m.card(render.Paint(title, color), body, footer, true)
	case helpScreen:
		body := "Arrows / WASD  Move cursor\nEnter          Select / confirm\nEsc            Cancel / open menu\nF              Flip board\nH              Review history\n?              Open help\nCtrl+C         Quit\n\nSelect a piece, then its destination."
		content = m.card("HELP", body, "Esc back", false)
	case errorScreen:
		content = m.card("UNABLE TO CONTINUE", "The AI could not generate a move.\nYour current position has been preserved.", footer, false)
	}
	return render.Window(content, m.width, m.height)
}

func (m application) resultText() (string, string, string) {
	s, _ := m.match.Game().Snapshot()
	reason := map[game.EndReason]string{game.Checkmate: "Checkmate.", game.Stalemate: "Draw by stalemate.\nNo legal moves and not in check.", game.Resignation: "Game ended by resignation.", game.Agreement: "Draw by agreement."}[s.Outcome.Reason]

	switch s.Outcome.Reason {
	case game.InsufficientMaterial:
		reason = "Draw by insufficient material."
	case game.Repetition:
		reason = "Draw claimed by threefold repetition."
	case game.FiftyMoves:
		reason = "Draw claimed under the 50-move rule."
	case game.FivefoldRepetition:
		reason = "Automatic draw by fivefold repetition."
	case game.SeventyFiveMoves:
		reason = "Automatic draw under the 75-move rule."
	}

	if s.Outcome.Winner == chess.NoColor {
		return "DRAW", "No winner.\n\n" + reason, render.Draw
	}

	if !m.config.computer {
		winner := map[chess.Color]string{chess.White: "Player 1(W)", chess.Black: "Player 2(B)"}[s.Outcome.Winner]
		return winner + " WINS", reason, render.Accent
	}

	if s.Outcome.Winner == m.match.State().AIColor {
		return "DEFEAT", "The computer won.\n\n" + reason, render.Alert
	}

	return "VICTORY", "You won as " + sideText(s.Outcome.Winner) + ".\n\n" + reason, render.Victory
}
