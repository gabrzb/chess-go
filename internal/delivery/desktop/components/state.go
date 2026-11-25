package components

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/application"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
	"github.com/gabrzb/chess-go/internal/engine"
)

type Options struct {
	// Difficulty, when set, starts a game against the AI without the menu.
	Difficulty engine.Difficulty
	HumanColor chess.Color
	Seed       int64
	// Update schedules a function on the UI thread and redraws the window.
	Update func(func())
	Quit   func()
}

type Screen int

const (
	MenuScreen Screen = iota
	SetupScreen
	GameScreen
)

const (
	whiteSide = iota
	blackSide
	randomSide
)

const (
	defaultDifficulty = 1
	livePosition      = -1
)

var (
	difficulties = []engine.Difficulty{engine.Easy, engine.Medium, engine.Hard}
	LevelLabels  = []string{"Easy", "Medium", "Hard"}
	SideLabels   = []string{"White", "Black", "Random"}
)

// App state is owned by the UI thread.
type App struct {
	ctx    context.Context
	seed   int64
	update func(func())
	Quit   func()

	Screen          Screen
	vsAI            bool
	Level, Side     int
	game            *game.Game
	opponent        application.Opponent
	session         uint64
	Flipped         bool
	selected        chess.Square
	dragTarget      chess.Square
	sliding         *chess.Move
	targets         []chess.Move
	promoting       []chess.Move
	review          int
	moveHistory     []string
	message         string
	choosing        bool
	paused          bool
	confirmResign   bool
	confirmLeave    bool
	drawOffered     bool
	showResult      bool
	resultDismissed bool
}

func DifficultyDescription(level int) string {
	depth := difficulties[level].Depth()
	if depth == 1 {
		return "Searches one ply ahead."
	}
	return fmt.Sprintf("Searches up to %d plies ahead.", depth)
}

func New(ctx context.Context, options Options) (*App, error) {
	a := &App{ctx: ctx, seed: options.Seed, update: options.Update, Quit: options.Quit, Level: defaultDifficulty, Side: whiteSide}
	if options.HumanColor == chess.Black {
		a.Side = blackSide
	}
	if options.Difficulty == "" {
		return a, nil
	}
	a.Level = slices.Index(difficulties, options.Difficulty)
	if a.Level < 0 {
		return nil, engine.ErrInvalidDifficulty
	}
	return a, a.start(true)
}

func (a *App) start(vsAI bool) error {
	a.opponent.Stop()
	a.session++
	a.vsAI = vsAI
	a.game = game.NewGame()
	a.opponent = application.Opponent{}
	a.Screen, a.Flipped = GameScreen, false
	a.selected, a.targets, a.promoting, a.dragTarget, a.sliding = chess.NoSquare, nil, nil, chess.NoSquare, nil
	a.review, a.moveHistory, a.message = livePosition, nil, ""
	a.choosing, a.paused, a.confirmResign, a.confirmLeave = false, false, false, false
	a.drawOffered, a.showResult, a.resultDismissed = false, false, false
	if !vsAI {
		return nil
	}
	player, err := engine.NewPlayer(difficulties[a.Level], a.seed)
	if err != nil {
		return err
	}
	human := chess.White
	if a.Side == blackSide || (a.Side == randomSide && rand.IntN(2) == blackSide) {
		human = chess.Black
	}
	a.Flipped = human == chess.Black
	a.opponent, err = application.NewOpponent(a.ctx, a.game, player, other(human), a.session)
	return err
}

func (a *App) StartGame(vsAI bool) {
	if err := a.start(vsAI); err != nil {
		a.message = err.Error()
	}
	a.StartAI()
}

func (a *App) toMenu() {
	a.opponent.Stop()
	a.session++
	a.Screen = MenuScreen
	a.paused, a.confirmLeave, a.showResult = false, false, false
}

func (a *App) record(before game.Snapshot, move chess.Move) {
	after, _ := a.game.Snapshot()
	a.moveHistory = append(a.moveHistory, game.SAN(before, after, move))
	a.showOutcome()
}

func (a *App) showOutcome() {
	if s, _ := a.game.Snapshot(); s.Status != chess.StatusOngoing && !a.resultDismissed {
		a.opponent.Stop()
		a.paused, a.showResult = false, true
	}
}

func (a *App) play(move chess.Move) {
	before, _ := a.game.Snapshot()
	if _, err := a.game.Play(move); err != nil {
		a.message = err.Error()
		return
	}
	a.message, a.sliding = "", nil
	a.record(before, move)
	a.StartAI()
}

func (a *App) resign() {
	a.opponent.Stop()
	if err := a.game.Resign(a.humanColor()); err != nil {
		a.message = err.Error()
	}
	a.showOutcome()
}

func (a *App) claimDraw() {
	a.opponent.Stop()
	if err := a.game.ClaimDraw(); err != nil {
		a.message = err.Error()
	}
	a.showOutcome()
}

func (a *App) offerDraw() {
	s, _ := a.game.Snapshot()
	if err := a.game.OfferDraw(s.SideToMove); err != nil {
		a.message = err.Error()
		return
	}
	a.paused, a.drawOffered = false, true
}

func (a *App) answerDraw(accept bool) {
	s, _ := a.game.Snapshot()
	a.drawOffered = false
	if err := a.game.RespondDraw(other(s.DrawOffer), accept); err != nil {
		a.message = err.Error()
	}
	a.showOutcome()
}

func (a *App) resume() {
	a.paused = false
	a.StartAI()
}

func (a *App) setReview(index int) {
	a.selected, a.targets = chess.NoSquare, nil
	if index >= a.game.HistoryLen()-1 || index < 0 {
		a.review = livePosition
		return
	}
	a.review = index
}

func (a *App) step(delta int) {
	at := a.review
	if at < 0 {
		at = a.game.HistoryLen() - 1
	}
	if at+delta >= 0 {
		a.setReview(at + delta)
	}
}

func (a *App) displayedPosition() game.Snapshot {
	if a.review >= 0 {
		s, _ := a.game.HistorySnapshot(a.review)
		return s
	}
	s, _ := a.game.Snapshot()
	return s
}

func (a *App) humanColor() chess.Color {
	if a.opponent.Enabled() {
		return other(a.opponent.Color())
	}
	s, _ := a.game.Snapshot()
	return s.SideToMove
}

func other(c chess.Color) chess.Color {
	if c == chess.White {
		return chess.Black
	}
	return chess.White
}

func colorName(c chess.Color) string {
	if c == chess.White {
		return "White"
	}
	return "Black"
}

func (a *App) SetUpdate(update func(func())) { a.update = update }

func (a *App) GameShortcuts(c *ui.Context) {
	switch {
	case c.Shortcut(0, ui.KeyEscape):
		if a.selected.Valid() {
			a.selected, a.targets = chess.NoSquare, nil
		} else {
			a.openGameMenu()
		}
	case c.Shortcut(0, ui.KeyLeft):
		a.step(-1)
	case c.Shortcut(0, ui.KeyRight):
		a.step(1)
	case c.Shortcut(0, ui.KeyF):
		a.Flipped = !a.Flipped
	}
}

func (a *App) openGameMenu() {
	if a.over() {
		a.showResult = true
	} else {
		a.paused = true
	}
}

func (a *App) over() bool {
	s, _ := a.game.Snapshot()
	return s.Status != chess.StatusOngoing
}
