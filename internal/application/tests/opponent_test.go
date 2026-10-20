package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/application"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

type chooserFunc func(context.Context, chess.Position) (chess.Move, error)

func (f chooserFunc) ChooseMove(ctx context.Context, p chess.Position) (chess.Move, error) {
	return f(ctx, p)
}

func opponent(t *testing.T, g *game.Game, chooser chess.MoveChooser, session uint64) application.Opponent {
	t.Helper()
	o, err := application.NewOpponent(t.Context(), g, chooser, chess.White, session)
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func start(t *testing.T, o *application.Opponent) func() application.SearchResult {
	t.Helper()
	search, err := o.Start()
	if err != nil || search == nil {
		t.Fatalf("missing search: %v", err)
	}
	return search
}

func TestOpponentSearchLifecycle(t *testing.T) {
	g := game.NewGame()
	calls := 0
	var searchContext context.Context
	move := chess.Move{From: chess.E2, To: chess.E4}
	o := opponent(t, g, chooserFunc(func(ctx context.Context, position chess.Position) (chess.Move, error) {
		calls++
		searchContext = ctx
		position.Board.RemovePiece(chess.E1)
		return move, nil
	}), 1)
	search := start(t, &o)
	if calls != 0 || !o.Thinking() {
		t.Fatal("search was not deferred")
	}
	if duplicate, err := o.Start(); duplicate != nil || err != nil {
		t.Fatal("started duplicate search")
	}
	result := search()
	if g.Position() != chess.NewInitialPosition() {
		t.Fatal("search changed live game")
	}
	if accepted, err := o.Accept(result); !accepted || err != nil {
		t.Fatalf("result rejected: %v", err)
	}
	if calls != 1 || o.Thinking() || g.HistoryLen() != 2 || !errors.Is(searchContext.Err(), context.Canceled) {
		t.Fatal("incorrect search completion")
	}
	if accepted, err := o.Accept(result); accepted || err != nil || g.HistoryLen() != 2 {
		t.Fatal("duplicate result applied")
	}
	if search, err := o.Start(); search != nil || err != nil {
		t.Fatal("searched on human turn")
	}
}

func TestOpponentDiscardsCanceledAndPreviousSessionResults(t *testing.T) {
	g := game.NewGame()
	var contexts []context.Context
	chooser := chooserFunc(func(ctx context.Context, _ chess.Position) (chess.Move, error) {
		contexts = append(contexts, ctx)
		return chess.Move{From: chess.E2, To: chess.E4}, nil
	})
	o := opponent(t, g, chooser, 1)
	first := start(t, &o)
	o.Stop()
	old := first()
	if !errors.Is(contexts[0].Err(), context.Canceled) {
		t.Fatal("search not canceled")
	}
	if accepted, _ := o.Accept(old); accepted {
		t.Fatal("canceled result accepted")
	}
	second := start(t, &o)
	if accepted, _ := o.Accept(old); accepted || !o.Thinking() {
		t.Fatal("stale result affected new search")
	}
	other := opponent(t, g, chooser, 2)
	otherSearch := start(t, &other)
	if accepted, _ := other.Accept(second()); accepted || !other.Thinking() {
		t.Fatal("previous session result accepted")
	}
	if g.HistoryLen() != 1 {
		t.Fatal("stale results changed game")
	}
	if accepted, err := other.Accept(otherSearch()); !accepted || err != nil {
		t.Fatal("current session failed", err)
	}
	o.Stop()
}

func TestOpponentErrorsAndRetry(t *testing.T) {
	cause := errors.New("search failed")
	for _, tc := range []struct {
		name        string
		move        chess.Move
		cause, want error
		change      bool
	}{
		{"search failure", chess.Move{}, cause, cause, false},
		{"illegal move", chess.Move{From: chess.E2, To: chess.E5}, nil, chess.ErrIllegalMove, false},
		{"position changed", chess.Move{From: chess.E2, To: chess.E4}, nil, application.ErrPositionChanged, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			calls := 0
			o := opponent(t, g, chooserFunc(func(context.Context, chess.Position) (chess.Move, error) {
				calls++
				if calls == 1 {
					return tc.move, tc.cause
				}
				return chess.Move{From: chess.E2, To: chess.E4}, nil
			}), 1)
			result := start(t, &o)()
			if tc.change {
				if _, err := g.Play(chess.Move{From: chess.D2, To: chess.D4}); err != nil {
					t.Fatal(err)
				}
			}
			before, history := g.Position(), g.HistoryLen()
			if accepted, err := o.Accept(result); !accepted || !errors.Is(err, tc.want) {
				t.Fatalf("got %v, %v", accepted, err)
			}
			if o.Thinking() || g.Position() != before || g.HistoryLen() != history {
				t.Fatal("failure changed game or left search pending")
			}
			if !tc.change {
				if accepted, err := o.Accept(start(t, &o)()); !accepted || err != nil {
					t.Fatal("retry failed", err)
				}
			}
		})
	}
}

func TestOpponentConfigurationAndZeroValue(t *testing.T) {
	chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) { return chess.Move{}, nil })
	for _, tc := range []struct {
		ctx     context.Context
		game    *game.Game
		chooser chess.MoveChooser
		color   chess.Color
	}{
		{nil, game.NewGame(), chooser, chess.White},
		{t.Context(), nil, chooser, chess.White},
		{t.Context(), game.NewGame(), nil, chess.White},
		{t.Context(), game.NewGame(), chooser, chess.NoColor},
	} {
		if _, err := application.NewOpponent(tc.ctx, tc.game, tc.chooser, tc.color, 1); !errors.Is(err, application.ErrInvalidAIConfig) {
			t.Fatal("invalid configuration accepted")
		}
	}
	var o application.Opponent
	o.Stop()
	if search, err := o.Start(); search != nil || err != nil || o.Enabled() || o.IsTurn() || o.Thinking() {
		t.Fatal("zero value is not inactive")
	}
}
