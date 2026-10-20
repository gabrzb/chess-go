package game

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func cycle(t *testing.T, g *Game) {
	t.Helper()
	for _, move := range []chess.Move{{From: chess.G1, To: chess.F3}, {From: chess.G8, To: chess.F6}, {From: chess.F3, To: chess.G1}, {From: chess.F6, To: chess.G8}} {
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRepetitionClaimsAndAutomaticDraw(t *testing.T) {
	g := NewGame()
	if err := g.ClaimDraw(); !errors.Is(err, ErrDrawNotClaimable) {
		t.Fatal(err)
	}
	cycle(t, g)
	if g.ClaimableDraw() != NoReason {
		t.Fatal("draw available too early")
	}
	cycle(t, g)
	if g.ClaimableDraw() != Repetition || g.status != chess.StatusOngoing {
		t.Fatal("third repetition must allow a claim")
	}
	if err := g.ClaimDraw(); err != nil {
		t.Fatal(err)
	}
	if g.outcome.Reason != Repetition {
		t.Fatal("missing repetition outcome")
	}
	if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); !errors.Is(err, ErrGameOver) {
		t.Fatal(err)
	}

	g = NewGame()
	for range 4 {
		cycle(t, g)
	}
	if g.status != chess.StatusFivefoldRepetition {
		t.Fatal("fifth repetition did not end game")
	}
	s, err := g.HistorySnapshot(g.HistoryLen() - 1)
	if err != nil || s.Outcome.Reason != FivefoldRepetition {
		t.Fatal("history lost draw outcome", err)
	}
}

func TestMoveCountDraw(t *testing.T) {
	for _, tc := range []struct {
		clock uint16
		want  chess.Status
	}{{98, chess.StatusOngoing}, {99, chess.StatusOngoing}, {149, chess.StatusSeventyFiveMoves}} {
		g := NewGame()
		g.position.HalfmoveClock = tc.clock
		if _, err := g.Play(chess.Move{From: chess.G1, To: chess.F3}); err != nil {
			t.Fatal(err)
		}
		if g.status != tc.want {
			t.Fatalf("clock %d: status %v", tc.clock, g.status)
		}
		if tc.clock == 99 {
			if g.ClaimableDraw() != FiftyMoves {
				t.Fatal("missing 50-move claim")
			}
			if err := g.ClaimDraw(); err != nil {
				t.Fatal(err)
			}
			if g.outcome.Reason != FiftyMoves {
				t.Fatal("missing 50-move outcome")
			}
		}
	}
	g := NewGame()
	g.position.HalfmoveClock = 149
	if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); err != nil {
		t.Fatal(err)
	}
	if g.position.HalfmoveClock != 0 || g.status != chess.StatusOngoing {
		t.Fatal("pawn did not reset clock")
	}
}

func TestMateTakesPrecedenceOverMoveCount(t *testing.T) {
	g := NewGame()
	for _, move := range []chess.Move{{From: chess.F2, To: chess.F3}, {From: chess.E7, To: chess.E5}, {From: chess.G2, To: chess.G4}} {
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
	}
	g.position.HalfmoveClock = 149
	if _, err := g.Play(chess.Move{From: chess.D8, To: chess.H4}); err != nil {
		t.Fatal(err)
	}
	if g.status != chess.StatusCheckmate {
		t.Fatal("move count overrode checkmate")
	}
}
