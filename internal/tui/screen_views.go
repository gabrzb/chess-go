package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

const (
	background = "#16161D"
	foreground = "#DCD7BA"
	muted      = "#72727F"
	accent     = "#7E9CD8"
	victory    = "#98BB6C"
	drawColor  = "#E6C384"
	alert      = "#E46876"
)

var logo = strings.Join([]string{
	" ██████╗██╗  ██╗███████╗███████╗███████╗",
	"██╔════╝██║  ██║██╔════╝██╔════╝██╔════╝",
	"██║     ███████║█████╗  ███████╗███████╗",
	"██║     ██╔══██║██╔══╝  ╚════██║╚════██║",
	"╚██████╗██║  ██║███████╗███████║███████║",
	" ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝",
}, "\n")

func paint(text, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func centeredBlock(text string, width int) string {
	left := strings.Repeat(" ", max(0, (width-lipgloss.Width(text))/2))
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		lines[index] = left + line + strings.Repeat(" ", max(0, width-len(left)-lipgloss.Width(line)))
	}
	return strings.Join(lines, "\n")
}

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
	bar := m.band(width, double)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	content := bar + "\n\n" + center.Bold(true).Render(title) + "\n\n" + bar
	if body != "" {
		content += "\n\n" + center.Render(ansi.Wrap(body, width, " "))
	}
	options := m.menuOptions()
	if len(options) > 0 {
		items := make([]string, len(options))
		for index, option := range options {
			items[index] = "  " + option
			if index == m.choice {
				items[index] = paint("▸ "+option, accent)
			}
		}
		menu := strings.Join(items, "\n")
		content += "\n\n" + centeredBlock(menu, width)
	}
	if footer != "" {
		content += "\n\n" + center.Render(paint(ansi.Wrap(footer, width, " "), muted))
	}
	return content
}

func (m application) View() tea.View {
	var content string
	footer := "↑↓ navegar · Enter selecionar"
	switch m.screen {
	case mainScreen:
		title := m.logoText()
		if m.width < 48 || m.height < 23 {
			title = paint("CHESS", accent)
		}
		setting := "ligadas"
		if m.noAnimations {
			setting = "desligadas"
		}
		content = m.card(title, "", footer+"\n? ajuda · A animações: "+setting, true)
	case configScreen:
		content = m.card("NOVA PARTIDA\nJOGADOR VS IA", "", "↑↓ navegar · ←→ alterar\nEnter confirmar · Esc voltar", true)
	case playingScreen, finalBoardScreen:
		if len(m.match.promotionMoves) > 0 && m.screen == playingScreen {
			content = m.promotionCard()
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
		s, _ := m.match.game.Snapshot()
		content = m.card("OFERTA DE EMPATE", "As "+sideText(s.DrawOffer)+" propuseram empate.\nAs "+sideText(opposite(s.DrawOffer))+" desejam aceitar?", footer, false)
	case resultScreen:
		title, body, color := m.resultText()
		content = m.card(paint(title, color), body, footer, true)
	case helpScreen:
		body := "Setas / WASD   Mover cursor\nEnter          Selecionar / confirmar\nEsc            Cancelar / abrir menu\nF              Inverter tabuleiro\nH              Revisar histórico\n?              Abrir ajuda\nCtrl+C         Encerrar\n\nSelecione uma peça e depois o destino."
		content = m.card("AJUDA", body, "Esc voltar", false)
	case errorScreen:
		content = m.card("NÃO FOI POSSÍVEL CONTINUAR", "A IA não conseguiu gerar a jogada.\nSua posição atual foi preservada.", footer, false)
	}
	if lipgloss.Width(content) > m.width || lipgloss.Height(content) > m.height {
		content = ansi.Wrap("CHESS\n\nTerminal pequeno demais.\nAmplie a janela para continuar.\n\nEsc: menu · Ctrl+C: sair", max(1, m.width), " ")
		content = lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(content)
	}
	surface := lipgloss.NewStyle().Foreground(lipgloss.Color(foreground)).Background(lipgloss.Color(background))
	content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, surface.Render(content), lipgloss.WithWhitespaceStyle(surface))
	view := tea.NewView(content)
	view.AltScreen = true
	view.BackgroundColor = lipgloss.Color(background)
	return view
}

func (m application) promotionCard() string {
	width := min(48, max(18, m.width-4))
	bar := paint(strings.Repeat("=", width), muted)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	names := []string{"Dama", "Torre", "Bispo", "Cavalo"}
	kinds := []chess.PieceType{chess.Queen, chess.Rook, chess.Bishop, chess.Knight}
	var items []string
	for index, name := range names {
		option := pieceSymbol(chess.Piece{Type: kinds[index], Color: m.match.game.Position().SideToMove}) + "  " + name
		text := "  " + option
		if index == m.choice {
			text = paint("▸ "+option, accent)
		}
		items = append(items, text)
	}
	return bar + "\n\n" + center.Bold(true).Render("PROMOÇÃO DE PEÃO") + "\n" + center.Render("Escolha a nova peça para "+m.match.cursor.String()+".") + "\n\n" + bar + "\n\n" + centeredBlock(strings.Join(items, "\n"), width) + "\n\n" + center.Render("↑↓ escolher · Enter confirmar\nQ/R/B/N atalhos · Esc cancelar")
}

func (m application) resultText() (string, string, string) {
	s, _ := m.match.game.Snapshot()
	reason := map[game.EndReason]string{game.Checkmate: "Xeque-mate.", game.Stalemate: "Empate por afogamento.\nSem jogadas legais e sem xeque.", game.Resignation: "Partida encerrada por desistência.", game.Agreement: "Empate por acordo."}[s.Outcome.Reason]
	if s.Outcome.Winner == chess.NoColor {
		return "EMPATE", "Nenhum vencedor.\n\n" + reason, drawColor
	}
	if !m.config.computer {
		return strings.ToUpper(sideText(s.Outcome.Winner)) + " VENCERAM", "Jogador " + map[chess.Color]string{chess.White: "1", chess.Black: "2"}[s.Outcome.Winner] + "\n\n" + reason, accent
	}
	if s.Outcome.Winner == m.match.aiColor {
		return "DERROTA", "O computador venceu.\n\n" + reason, alert
	}
	return "VITÓRIA", "Você venceu de " + sideText(s.Outcome.Winner) + ".\n\n" + reason, victory
}
