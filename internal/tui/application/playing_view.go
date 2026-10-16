package application

import (
	"github.com/gabrzb/chess-go/internal/tui/render"

	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func (m application) playerInfo(s game.Snapshot) string {
	status := "Vez das " + sideText(s.SideToMove)
	if m.config.computer {
		status += "\nVocê joga de " + sideText(opposite(m.match.State().AIColor)) + "."
		if s.SideToMove == m.match.State().AIColor {
			status += "\nComputador · " + levelNames[m.config.level]
		}
	} else {
		status += "\nJogador " + map[chess.Color]string{chess.White: "1", chess.Black: "2"}[s.SideToMove]
	}
	if m.match.State().Thinking {
		status += "\n" + m.thinkingText()
	}
	if s.Check {
		warning := "O rei " + map[chess.Color]string{chess.White: "branco", chess.Black: "preto"}[s.SideToMove] + " está ameaçado."
		if m.config.computer && s.SideToMove != m.match.State().AIColor {
			warning = "Seu rei está ameaçado. Elimine o xeque."
		}
		status += "\n" + render.Paint("XEQUE · "+strings.ToUpper(sideText(s.SideToMove)), render.Alert) + "\n" + warning
	}
	if m.match.State().Reviewing {
		status += fmt.Sprintf("\nREVISÃO · posição %d/%d", m.match.State().HistoryIndex+1, m.match.Game().HistoryLen())
	}
	if m.screen == finalBoardScreen {
		title, body, _ := m.resultText()
		paragraphs := strings.Split(body, "\n\n")
		status = "PARTIDA ENCERRADA\n" + title + "\n" + paragraphs[len(paragraphs)-1]
		if m.match.State().Reviewing {
			status += fmt.Sprintf("\nREVISÃO · posição %d/%d", m.match.State().HistoryIndex+1, m.match.Game().HistoryLen())
		}
	}
	return status
}

func (m application) gameFooter() string {
	if m.match.State().Reviewing {
		return "←/A anterior · →/D próximo · H/Esc retornar · F inverter · ? ajuda"
	}
	if m.screen == finalBoardScreen {
		return "F inverter · H histórico · Esc voltar ao resultado"
	}
	if m.match.State().Selected.Valid() {
		piece, _ := m.match.Game().Position().Board.PieceAt(m.match.State().Selected)
		name := map[chess.PieceType]string{chess.Pawn: "PEÃO", chess.Knight: "CAVALO", chess.Bishop: "BISPO", chess.Rook: "TORRE", chess.Queen: "DAMA", chess.King: "REI"}[piece.Type]
		return name + " · " + m.match.State().Selected.String() + " — escolha o destino.\nEnter confirmar · Esc cancelar seleção · ? ajuda"
	}
	if m.match.State().Thinking {
		return "IA calculando · Esc menu · H histórico · ? ajuda · Ctrl+C sair"
	}
	return "Setas/WASD mover · Enter selecionar\nEsc menu · F inverter · H histórico · ? ajuda"
}

func (m application) historyPanel(end, rows int) string {
	content := "HISTÓRICO\n    Brancas   Pretas"
	if end == 0 {
		return content + "\n    —         —"
	}
	last := (end - 1) / 2
	first := max(0, last-rows+1)
	for full := first; full <= last && rows > 0; full++ {
		white, black := m.labels[full*2], ""
		if full*2+1 < end {
			black = m.labels[full*2+1]
		}
		content += fmt.Sprintf("\n%2d. %-9s %s", full+1, white, black)
	}
	return content
}

func (m application) boardContent() string {
	s, err := m.match.DisplayedSnapshot()
	if err != nil {
		return "Não foi possível carregar a posição."
	}
	width := min(76, max(18, m.width-4))
	wide := m.width >= 72 && m.height >= 22
	mode := "LOCAL · JOGADOR VS JOGADOR"
	if m.config.computer {
		mode = "JOGADOR VS IA · " + levelNames[m.config.level]
	}
	if !wide {
		mode = "LOCAL"
		if m.config.computer {
			mode = "VS IA"
		}
	}
	bar := render.Paint(strings.Repeat("=", width), render.Muted)
	header := bar + "\nCHESS · " + mode + "\n" + bar
	foot := render.Paint(ansi.Wrap(m.gameFooter(), width, " "), render.Muted)
	status := m.playerInfo(s)
	lastMove := "Último lance: —"
	if s.LastMove != nil {
		lastMove = "Último lance: " + s.LastMove.From.String() + " → " + s.LastMove.To.String()
	}
	board := m.match.RenderBoard(s.Board)
	var body string
	if wide {
		panelWidth := width - lipgloss.Width(board) - 4
		info := "PARTIDA\n\n" + ansi.Wrap(status, panelWidth, " ") + "\n\n" + lastMove
		available := m.height - lipgloss.Height(header) - lipgloss.Height(foot) - 4
		if m.match.State().Message != "" {
			available--
		}
		rows := min(6, max(0, available-lipgloss.Height(info)-4))
		if rows > 0 {
			end := len(m.labels)
			if m.match.State().Reviewing {
				end = min(end, m.match.State().HistoryIndex)
			}
			info += "\n\n" + m.historyPanel(end, rows)
		}
		body = lipgloss.JoinHorizontal(lipgloss.Top, "TABULEIRO\n\n"+board, "    ", info)
	} else {
		// Compact mode prioritizes the board and current turn over the move list.
		status = strings.Split(status, "\n")[0]
		if m.config.computer {
			status += " · Você: " + sideText(opposite(m.match.State().AIColor))
		}
		if s.Check {
			status += "\n" + render.Paint("XEQUE · elimine a ameaça ao rei", render.Alert)
		}
		if m.match.State().Thinking {
			status += "\n" + m.thinkingText()
		}
		if m.match.State().Reviewing {
			status += fmt.Sprintf("\nRevisão %d/%d", m.match.State().HistoryIndex+1, m.match.Game().HistoryLen())
		}
		body = board + "\n\n" + ansi.Wrap(status, width, " ") + "\n" + lastMove
	}
	if m.match.State().Message != "" {
		message := "Movimento inválido. Escolha uma casa destacada."
		if m.match.State().Message == chess.ErrEmptySource.Error() {
			message = "Selecione uma peça do jogador da vez."
		}
		body += "\n" + render.Paint(ansi.Wrap(message, width, " "), render.Alert)
	}
	return header + "\n\n" + body + "\n\n" + foot
}

func (m application) boardFits() bool {
	if m.match.Game() == nil {
		return false
	}
	content := m.boardContent()
	return lipgloss.Width(content) <= m.width && lipgloss.Height(content) <= m.height
}
