package application

import "github.com/gabrzb/chess-go/internal/engine"

type sideChoice int

const (
	whiteSide sideChoice = iota
	blackSide
	randomSide
)

var sideNames = [...]string{"White", "Black", "Random"}

var difficulties = [...]struct {
	value engine.Difficulty
	label string
}{
	{engine.Easy, "Easy"},
	{engine.Medium, "Medium"},
	{engine.Hard, "Hard"},
}

func difficultyLabel(value engine.Difficulty) string {
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
