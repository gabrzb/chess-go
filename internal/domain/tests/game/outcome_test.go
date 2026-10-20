package game_test

import (
	"errors"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

func TestExplicitGameEndings(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		g := game.NewGame()
		before := g.Position()
		if err := g.Resign(color); err != nil {
			t.Fatal(err)
		}
		s := snapshot(t, g)
		if s.Outcome.Winner == color || s.Outcome.Winner == chess.NoColor || s.Outcome.Reason != game.Resignation || g.Position() != before || g.HistoryLen() != 1 {
			t.Fatal("Resignation changed the position or assigned the wrong winner")
		}
		if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); !errors.Is(err, game.ErrGameOver) {
			t.Fatal("Game continued after resignation")
		}
	}
	g := game.NewGame()
	for _, accept := range []bool{false, true} {
		if err := g.OfferDraw(chess.White); err != nil {
			t.Fatal(err)
		}
		if err := g.RespondDraw(chess.White, true); !errors.Is(err, game.ErrDrawOffer) {
			t.Fatal("Offering player accepted their own offer")
		}
		if err := g.RespondDraw(chess.Black, accept); err != nil {
			t.Fatal(err)
		}
		s := snapshot(t, g)
		if s.DrawOffer != chess.NoColor || (s.Status == chess.StatusDrawAgreed) != accept || s.Outcome.Winner != chess.NoColor {
			t.Fatal("Incorrect draw response")
		}
	}
	if snapshot(t, g).Outcome.Reason != game.Agreement {
		t.Fatal("Missing draw reason")
	}
	if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); !errors.Is(err, game.ErrGameOver) {
		t.Fatal("Game continued after agreement")
	}
}
