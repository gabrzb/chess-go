package tui

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand/v2"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

type screen uint8

const (
	mainScreen screen = iota
	configScreen
	playingScreen
	pauseScreen
	resignScreen
	abandonScreen
	drawScreen
	resultScreen
	finalBoardScreen
	helpScreen
	errorScreen
)

var levels = []ai.Difficulty{ai.Easy, ai.Medium, ai.Hard}
var levelNames = []string{"Fácil", "Normal", "Difícil"}
var sideNames = []string{"Brancas", "Pretas", "Aleatório"}

type Options struct {
	Difficulty   ai.Difficulty
	HumanColor   chess.Color
	RandomColor  bool
	Seed         int64
	StartGame    bool
	NoAnimations bool
	Logger       *log.Logger
}

type settings struct {
	computer    bool
	side, level int
	seed        int64
}

type application struct {
	ctx                context.Context
	logger             *log.Logger
	match              model
	config             settings
	screen, helpBack   screen
	choice, helpChoice int
	width, height      int
	session            uint64
	initial            tea.Cmd
	noAnimations       bool
	labels             []string
	introDone          bool
	frame              int
	animationID        uint64
}

func NewApplication(ctx context.Context, options Options) (tea.Model, error) {
	if ctx == nil {
		return nil, errors.New("application requires a context")
	}
	if options.HumanColor != chess.NoColor && options.HumanColor != chess.White && options.HumanColor != chess.Black {
		return nil, chess.ErrInvalidColor
	}
	if options.Difficulty != "" && options.Difficulty.Depth() == 0 {
		return nil, ai.ErrInvalidDifficulty
	}
	m := application{ctx: ctx, width: 80, height: 24, noAnimations: options.NoAnimations, logger: options.Logger,
		config: settings{computer: options.Difficulty != "", level: 1, seed: options.Seed}}
	if m.logger == nil {
		m.logger = log.New(io.Discard, "", 0)
	}
	if options.HumanColor == chess.Black {
		m.config.side = 1
	}
	if options.RandomColor {
		m.config.side = 2
	}
	for index, level := range levels {
		if level == options.Difficulty {
			m.config.level = index
		}
	}
	if options.StartGame {
		m.initial = m.beginMatch()
	}
	m.initial = combineCommands(m.initial, m.animationTick())
	return m, nil
}

func (m application) Init() tea.Cmd { return m.initial }

func (m *application) show(next screen) {
	if m.screen == mainScreen && next != mainScreen {
		m.introDone = true
	}
	if next != playingScreen && next != finalBoardScreen {
		m.match.stopAITurn()
	}
	m.screen, m.choice = next, 0
	m.frame = 0
	m.animationID++
}

func (m *application) resume() tea.Cmd {
	m.show(playingScreen)
	return m.match.startAITurn()
}

func (m *application) beginMatch() tea.Cmd {
	m.match.stopAITurn()
	m.session++
	m.match = NewModel(game.NewGame()).(model)
	m.match.sessionID = m.session
	m.labels = nil
	if m.config.computer {
		human := chess.White
		if m.config.side == 1 || (m.config.side == 2 && rand.IntN(2) == 1) {
			human = chess.Black
		}
		player, err := ai.NewPlayer(levels[m.config.level], m.config.seed)
		if err != nil {
			m.logger.Print(err)
			m.show(errorScreen)
			return nil
		}
		m.match.opponent, m.match.aiContext, m.match.aiColor = player, m.ctx, opposite(human)
		if human == chess.Black {
			m.match.cursor = chess.A7
		}
	}
	return m.resume()
}

func opposite(color chess.Color) chess.Color {
	if color == chess.White {
		return chess.Black
	}
	return chess.White
}

func sideText(color chess.Color) string {
	if color == chess.White {
		return "brancas"
	}
	if color == chess.Black {
		return "pretas"
	}
	return "—"
}

func (m application) resigningSide() chess.Color {
	if m.config.computer {
		return opposite(m.match.aiColor)
	}
	return m.match.game.Position().SideToMove
}
