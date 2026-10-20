package application

import (
	"strings"

	"github.com/gabrzb/chess-go/internal/tui/render"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func (m application) card(title, body, footer string, double bool) string {
	width := render.ContentWidth(m.width, render.MenuWidth)

	return render.Card(title, body, footer, m.band(width, double), m.menuOptions(), m.choice, width)
}

func (m application) View() tea.View {
	var content string
	footer := "↑↓ navegar · Enter selecionar"

	switch m.screen {
	case mainScreen:
		title := m.logoText()
		if m.width < render.LogoMinWidth || m.height < render.LogoMinHeight {
			title = render.Paint("CHESS", render.Accent)
		}

		setting := "ligadas"

		if m.noAnimations {
			setting = "desligadas"
		}

		content = m.card(title, "", footer+"\n? ajuda · A animações: "+setting, true)
	case configScreen:
		content = m.card("NOVA PARTIDA\nJOGADOR VS IA", "", "↑↓ navegar · ←→ alterar\nEnter confirmar · Esc voltar", true)
	case playingScreen, finalBoardScreen:
		if m.match.State().PromotionPending && m.screen == playingScreen {
			content = render.Promotion(m.match.State().Cursor, m.match.Game().Position().SideToMove, m.choice, m.match.State().PromotionChoices, render.ContentWidth(m.width, render.PromotionWidth))
		} else {
			content = m.boardContent()
		}
	case pauseScreen:
		content = m.card("MENU DA PARTIDA", "", footer+"\nEsc continuar", false)
	case resignScreen:
		content = m.card("CONFIRMAR DESISTÊNCIA?", "As "+sideText(opposite(m.resigningSide()))+" serão as vencedoras.", footer+"\nEsc voltar", false)
	case abandonScreen:
		content = m.card("ENCERRAR PARTIDA?", "A partida atual será descartada.", footer+"\nEsc voltar", false)
	case drawScreen:
		s, _ := m.match.Game().Snapshot()
		content = m.card("OFERTA DE EMPATE", "As "+sideText(s.DrawOffer)+" propuseram empate.\nAs "+sideText(opposite(s.DrawOffer))+" desejam aceitar?", footer, false)
	case resultScreen:
		title, body, color := m.resultText()
		content = m.card(render.Paint(title, color), body, footer, true)
	case helpScreen:
		body := "Setas / WASD   Mover cursor\nEnter          Selecionar / confirmar\nEsc            Cancelar / abrir menu\nF              Inverter tabuleiro\nH              Revisar histórico\n?              Abrir ajuda\nCtrl+C         Encerrar\n\nSelecione uma peça e depois o destino."
		content = m.card("AJUDA", body, "Esc voltar", false)
	case errorScreen:
		content = m.card("NÃO FOI POSSÍVEL CONTINUAR", "A IA não conseguiu gerar a jogada.\nSua posição atual foi preservada.", footer, false)
	}
	return render.Window(content, m.width, m.height)
}

func (m application) resultText() (string, string, string) {
	s, _ := m.match.Game().Snapshot()
	reason := map[game.EndReason]string{game.Checkmate: "Xeque-mate.", game.Stalemate: "Empate por afogamento.\nSem jogadas legais e sem xeque.", game.Resignation: "Partida encerrada por desistência.", game.Agreement: "Empate por acordo."}[s.Outcome.Reason]

	switch s.Outcome.Reason {
	case game.InsufficientMaterial:
		reason = "Empate por material insuficiente."
	case game.Repetition:
		reason = "Empate reivindicado por tripla repetição."
	case game.FiftyMoves:
		reason = "Empate reivindicado pela regra dos 50 movimentos."
	case game.FivefoldRepetition:
		reason = "Empate automático por quíntupla repetição."
	case game.SeventyFiveMoves:
		reason = "Empate automático pela regra dos 75 movimentos."
	}

	if s.Outcome.Winner == chess.NoColor {
		return "EMPATE", "Nenhum vencedor.\n\n" + reason, render.Draw
	}

	if !m.config.computer {
		return strings.ToUpper(sideText(s.Outcome.Winner)) + " VENCERAM", "Jogador " + map[chess.Color]string{chess.White: "1", chess.Black: "2"}[s.Outcome.Winner] + "\n\n" + reason, render.Accent
	}

	if s.Outcome.Winner == m.match.State().AIColor {
		return "DERROTA", "O computador venceu.\n\n" + reason, render.Alert
	}

	return "VITÓRIA", "Você venceu de " + sideText(s.Outcome.Winner) + ".\n\n" + reason, render.Victory
}
