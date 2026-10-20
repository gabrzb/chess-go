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
		status = "PARTIDA ENCERRADA\n" + title + "\n" + paragraphs[len(paragraphs)-1]
	} else {
		status = "Vez das " + sideText(s.SideToMove)

		if m.config.computer {
			if compact {
				status += " · Você: " + sideText(opposite(state.AIColor))
			} else {
				status += "\nVocê joga de " + sideText(opposite(state.AIColor)) + "."

				if s.SideToMove == state.AIColor {
					status += "\nComputador · " + difficultyLabel(m.config.level)
				}
			}
		} else if !compact {
			status += "\nJogador " + map[chess.Color]string{chess.White: "1", chess.Black: "2"}[s.SideToMove]
		}

		if state.Thinking {
			status += "\n" + render.Thinking(m.frame, !m.noAnimations)
		}
	}

	if state.Reviewing {
		status += fmt.Sprintf("\nREVISÃO · posição %d/%d", state.HistoryIndex+1, m.match.Game().HistoryLen())
	}

	return status
}

func (m application) gameFooter() string {
	state := m.match.State()

	if state.Reviewing {
		return "←/A anterior · →/D próximo · H/Esc retornar · F inverter · ? ajuda"
	}

	if m.screen == finalBoardScreen {
		return "F inverter · H histórico · Esc voltar ao resultado"
	}

	if state.Selected.Valid() {
		piece, _ := m.match.Game().Position().Board.PieceAt(state.Selected)
		name := strings.ToUpper(render.PieceName(piece.Type))

		return name + " · " + state.Selected.String() + " — escolha o destino.\nEnter confirmar · Esc cancelar seleção · ? ajuda"
	}

	if state.Thinking {
		return "IA calculando · Esc menu · H histórico · ? ajuda · Ctrl+C sair"
	}

	return "Setas/WASD mover · Enter selecionar\nEsc menu · F inverter · H histórico · ? ajuda"
}

func (m application) boardContent() string {
	s, err := m.match.DisplayedSnapshot()
	if err != nil {
		return "Não foi possível carregar a posição."
	}

	state := m.match.State()
	mode, compactMode := "LOCAL · JOGADOR VS JOGADOR", "LOCAL"

	if m.config.computer {
		mode, compactMode = "JOGADOR VS IA · "+difficultyLabel(m.config.level), "VS IA · "+difficultyLabel(m.config.level)
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
		notices = append(notices, "EMPATE DISPONÍVEL · Esc abre o menu para reivindicar.")
	}

	if s.Check && m.screen != finalBoardScreen {
		warning := "O rei " + map[chess.Color]string{chess.White: "branco", chess.Black: "preto"}[s.SideToMove] + " está ameaçado."

		if m.config.computer && s.SideToMove != state.AIColor {
			warning = "Seu rei está ameaçado. Elimine o xeque."
		}

		notices = append(notices, "XEQUE · "+strings.ToUpper(sideText(s.SideToMove))+" · "+warning)
	}
	if state.Message != "" {
		message := "Movimento inválido. Escolha uma casa destacada."

		if state.Message == chess.ErrEmptySource.Error() {
			message = "Selecione uma peça do jogador da vez."
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
