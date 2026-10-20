package application

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui/match"
)

type menuAction uint8

const (
	startLocal menuAction = iota
	configureAI
	quitApplication
	changeSide
	changeDifficulty
	startMatch
	returnToMenu
	continueMatch
	offerDraw
	showResignation
	showAbandonment
	confirmResignation
	rejectDraw
	acceptDraw
	showFinalBoard
	claimDraw
	retryAI
)

type menuItem struct {
	label  string
	action menuAction
}

// Labels and actions live together so optional entries cannot shift a command.
func (m application) menuItems() []menuItem {
	switch m.screen {
	case mainScreen:
		return []menuItem{{"Jogador vs Jogador", startLocal}, {"Jogador vs IA", configureAI}, {"Sair", quitApplication}}
	case configScreen:
		return []menuItem{
			{"Suas peças    ‹ " + sideNames[m.config.side] + " ›", changeSide},
			{"Dificuldade   ‹ " + difficultyLabel(m.config.level) + " ›", changeDifficulty},
			{"Iniciar partida", startMatch},
			{"Voltar", returnToMenu},
		}
	case pauseScreen:
		items := []menuItem{{"Continuar", continueMatch}}

		if !m.config.computer {
			items = append(items, menuItem{"Oferecer empate", offerDraw})
		}

		items = append(items, menuItem{"Desistir", showResignation}, menuItem{"Voltar ao menu principal", showAbandonment})

		if m.match.Game().ClaimableDraw() != game.NoReason {
			items = append(items, menuItem{"Reivindicar empate", claimDraw})
		}

		return items
	case resignScreen:
		return []menuItem{{"Continuar jogando", continueMatch}, {"Confirmar desistência", confirmResignation}}
	case abandonScreen:
		return []menuItem{{"Continuar jogando", continueMatch}, {"Encerrar e voltar ao menu", returnToMenu}}
	case drawScreen:
		return []menuItem{{"Recusar e continuar", rejectDraw}, {"Aceitar empate", acceptDraw}}
	case resultScreen:
		return []menuItem{{"Jogar novamente", startMatch}, {"Ver tabuleiro", showFinalBoard}, {"Menu principal", returnToMenu}}
	case errorScreen:
		return []menuItem{{"Tentar novamente", retryAI}, {"Voltar ao menu principal", showAbandonment}}
	default:
		return nil
	}
}

func (m application) menuOptions() []string {
	items := m.menuItems()
	labels := make([]string, len(items))

	for index, item := range items {
		labels[index] = item.label
	}

	return labels
}

func (m *application) activate() tea.Cmd {
	items := m.menuItems()

	if m.choice < 0 || m.choice >= len(items) {
		return nil
	}

	switch action := items[m.choice].action; action {
	case startLocal:
		m.config.computer = false
		return m.beginMatch()
	case configureAI:
		m.config.computer = true
		m.show(configScreen)
	case quitApplication:
		m.match.StopAI()
		return tea.Quit
	case changeSide, changeDifficulty:
		m.adjustConfig(1)
	case startMatch:
		return m.beginMatch()
	case returnToMenu:
		m.match.StopAI()
		m.match = match.Model{}
		m.show(mainScreen)
	case continueMatch:
		return m.resume()
	case retryAI:
		m.match.ClearMessage()
		return m.resume()
	case offerDraw:
		if err := m.match.Game().OfferDraw(m.match.Game().Position().SideToMove); err != nil {
			m.logger.Print(err)
			return nil
		}
		m.show(drawScreen)
	case showResignation:
		m.show(resignScreen)
	case showAbandonment:
		m.show(abandonScreen)
	case confirmResignation:
		if err := m.match.Game().Resign(m.resigningSide()); err != nil {
			m.logger.Print(err)
			return nil
		}
		m.show(resultScreen)
	case rejectDraw, acceptDraw:
		snapshot, err := m.match.Game().Snapshot()
		if err != nil {
			m.logger.Print(err)
			return nil
		}
		if err := m.match.Game().RespondDraw(opposite(snapshot.DrawOffer), action == acceptDraw); err != nil {
			m.logger.Print(err)
			return nil
		}
		if action == rejectDraw {
			return m.resume()
		}
		m.show(resultScreen)
	case showFinalBoard:
		m.match.ClearSelection()
		m.match.EndReview()
		m.show(finalBoardScreen)
	case claimDraw:
		if err := m.match.Game().ClaimDraw(); err != nil {
			m.logger.Print(err)
			return nil
		}
		m.show(resultScreen)
	}
	return nil
}
