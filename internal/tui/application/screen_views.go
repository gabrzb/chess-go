package application

import (
	"github.com/gabrzb/chess-go/internal/tui/render"

	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func (m application) menuOptions() []string {
	switch m.screen {
	case mainScreen:
		return []string{"Jogador vs Jogador", "Jogador vs IA", "Sair"}
	case configScreen:
		return []string{"Suas peças    ‹ " + sideNames[m.config.side] + " ›", "Dificuldade   ‹ " + levelNames[m.config.level] + " ›", "Iniciar partida", "Voltar"}
	case pauseScreen:
		if m.config.computer {
			return []string{"Continuar", "Desistir", "Voltar ao menu principal"}
		}
		return []string{"Continuar", "Oferecer empate", "Desistir", "Voltar ao menu principal"}
	case resignScreen:
		return []string{"Continuar jogando", "Confirmar desistência"}
	case abandonScreen:
		return []string{"Continuar jogando", "Encerrar e voltar ao menu"}
	case drawScreen:
		return []string{"Recusar e continuar", "Aceitar empate"}
	case resultScreen:
		return []string{"Jogar novamente", "Ver tabuleiro", "Menu principal"}
	case errorScreen:
		return []string{"Tentar novamente", "Voltar ao menu principal"}
	default:
		return nil
	}
}

func (m application) card(title, body, footer string, double bool) string {
	width := min(52, max(18, m.width-4))
	return render.Card(title, body, footer, m.band(width, double), m.menuOptions(), m.choice, width)
}

func (m application) View() tea.View {
	var content string
	footer := "↑↓ navegar · Enter selecionar"
	switch m.screen {
	case mainScreen:
		title := m.logoText()
		if m.width < 48 || m.height < 23 {
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
			content = render.Promotion(m.match.State().Cursor, m.match.Game().Position().SideToMove, m.choice, min(48, max(18, m.width-4)))
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
