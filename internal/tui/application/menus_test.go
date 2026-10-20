package application

import (
	"context"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess/ai"
)

func TestMenuActionsMatchDisplayedOptions(t *testing.T) {
	for _, computer := range []bool{false, true} {
		model, err := New(context.Background(), Options{StartGame: true, NoAnimations: true})
		if err != nil {
			t.Fatal(err)
		}
		m := model.(application)
		m.config.computer = computer
		m.show(pauseScreen)
		for index, item := range m.menuItems() {
			if m.menuOptions()[index] != item.label {
				t.Fatal("menu label diverged from its action")
			}
			if item.action == showResignation {
				m.choice = index
				m.activate()
				if m.screen != resignScreen {
					t.Fatal("optional draw item shifted resignation action")
				}
				break
			}
		}
	}
}

func TestConfigurationCyclesInBothDirections(t *testing.T) {
	model, err := New(context.Background(), Options{NoAnimations: true})
	if err != nil {
		t.Fatal(err)
	}
	m := model.(application)
	m.show(configScreen)
	for index, item := range m.menuItems() {
		m.choice = index
		switch item.action {
		case changeSide:
			m.adjustConfig(-1)
			if m.config.side != randomSide {
				t.Fatal("color did not wrap backwards")
			}
			m.adjustConfig(1)
			if m.config.side != whiteSide {
				t.Fatal("color did not wrap forwards")
			}
		case changeDifficulty:
			for range len(difficulties) {
				m.adjustConfig(1)
			}
			if m.config.level != ai.Medium {
				t.Fatal("difficulty did not wrap")
			}
		}
	}
}
