package application

import "github.com/gabrzb/chess-go/internal/chess/ai"

type sideChoice int

const (
	whiteSide sideChoice = iota
	blackSide
	randomSide
)

var sideNames = [...]string{"Brancas", "Pretas", "Aleatório"}

var difficulties = [...]struct {
	value ai.Difficulty
	label string
}{
	{ai.Easy, "Fácil"},
	{ai.Medium, "Normal"},
	{ai.Hard, "Difícil"},
}

func difficultyLabel(value ai.Difficulty) string {
	for _, difficulty := range difficulties {
		if difficulty.value == value {
			return difficulty.label
		}
	}
	return string(value)
}

func cycleIndex(index, step, count int) int {
	return (index + step + count) % count
}

func (m *application) adjustConfig(step int) {
	items := m.menuItems()
	if m.choice < 0 || m.choice >= len(items) {
		return
	}
	switch items[m.choice].action {
	case changeSide:
		m.config.side = sideChoice(cycleIndex(int(m.config.side), step, len(sideNames)))
	case changeDifficulty:
		for index, difficulty := range difficulties {
			if difficulty.value == m.config.level {
				m.config.level = difficulties[cycleIndex(index, step, len(difficulties))].value
				return
			}
		}
	}
}
