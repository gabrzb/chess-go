package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess/game"
)

func Run() error {
	program := tea.NewProgram(NewModel(game.NewGame()))
	_, err := program.Run()
	return err
}
