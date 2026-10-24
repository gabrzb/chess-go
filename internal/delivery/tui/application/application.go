package application

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand/v2"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/delivery/tui/match"
	"github.com/gabrzb/chess-go/internal/delivery/tui/render"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
	"github.com/gabrzb/chess-go/internal/engine"
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

type Options struct {
	Difficulty   engine.Difficulty
	HumanColor   chess.Color
	RandomColor  bool
	Seed         int64
	StartGame    bool
	NoAnimations bool
	Logger       *log.Logger
}

type settings struct {
	computer bool
	side     sideChoice
	level    engine.Difficulty
	seed     int64
}

type application struct {
	ctx                context.Context
	logger             *log.Logger
	match              match.Model
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

func New(ctx context.Context, options Options) (tea.Model, error) {
	if ctx == nil {
		return nil, errors.New("application requires a context")
	}

	if options.HumanColor != chess.NoColor && options.HumanColor != chess.White && options.HumanColor != chess.Black {
		return nil, chess.ErrInvalidColor
	}

	if options.Difficulty != "" && options.Difficulty.Depth() == 0 {
		return nil, engine.ErrInvalidDifficulty
	}
	m := application{
		ctx: ctx, width: render.DefaultWidth, height: render.DefaultHeight, noAnimations: options.NoAnimations, logger: options.Logger,
		config: settings{computer: options.Difficulty != "", level: engine.Medium, seed: options.Seed},
	}

	if m.logger == nil {
		m.logger = log.New(io.Discard, "", 0)
	}

	if options.HumanColor == chess.Black {
		m.config.side = blackSide
	}

	if options.RandomColor {
		m.config.side = randomSide
	}

	if options.Difficulty != "" {
		m.config.level = options.Difficulty
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
		m.match.StopAI()
	}

	m.screen, m.choice = next, 0
	m.frame = 0
	m.animationID++
}

func (m *application) resume() tea.Cmd {
	m.show(playingScreen)
	return m.match.StartAI()
}

func (m *application) beginMatch() tea.Cmd {
	m.match.StopAI()
	m.session++
	m.match = match.New(game.NewGame())
	m.labels = nil

	if m.config.computer {
		human := chess.White

		if m.config.side == blackSide || (m.config.side == randomSide && rand.IntN(2) == 1) {
			human = chess.Black
		}

		player, err := engine.NewPlayer(m.config.level, m.config.seed)
		if err != nil {
			m.logger.Print(err)
			m.show(errorScreen)
			return nil
		}

		m.match, err = match.NewWithAI(m.ctx, m.match.Game(), player, opposite(human), m.session)
		if err != nil {
			m.logger.Print(err)
			m.show(errorScreen)

			return nil
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
		return "white"
	}

	if color == chess.Black {
		return "black"
	}

	return "—"
}

func (m application) resigningSide() chess.Color {
	if m.config.computer {
		return opposite(m.match.State().AIColor)
	}

	return m.match.Game().Position().SideToMove
}
