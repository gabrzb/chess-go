package game_test

import (
	"errors"
	"github.com/gabrzb/chess-go/internal/domain/game"
	"github.com/gabrzb/chess-go/internal/domain/rules"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func cycle(t *testing.T, g *game.Game) {
	t.Helper()
	for _, move := range []chess.Move{{From: chess.G1, To: chess.F3}, {From: chess.G8, To: chess.F6}, {From: chess.F3, To: chess.G1}, {From: chess.F6, To: chess.G8}} {
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRepetitionClaimsAndAutomaticDraw(t *testing.T) {
	g := game.NewGame()
	if err := g.ClaimDraw(); !errors.Is(err, game.ErrDrawNotClaimable) {
		t.Fatal(err)
	}
	cycle(t, g)
	if g.ClaimableDraw() != game.NoReason {
		t.Fatal("draw available too early")
	}
	cycle(t, g)
	if g.ClaimableDraw() != game.Repetition || snapshot(t, g).Status != chess.StatusOngoing {
		t.Fatal("third repetition must allow a claim")
	}
	if err := g.ClaimDraw(); err != nil {
		t.Fatal(err)
	}
	if snapshot(t, g).Outcome.Reason != game.Repetition {
		t.Fatal("missing repetition outcome")
	}
	if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); !errors.Is(err, game.ErrGameOver) {
		t.Fatal(err)
	}

	g = game.NewGame()
	for range 4 {
		cycle(t, g)
	}
	if snapshot(t, g).Status != chess.StatusFivefoldRepetition {
		t.Fatal("fifth repetition did not end game")
	}
	s, err := g.HistorySnapshot(g.HistoryLen() - 1)
	if err != nil || s.Outcome.Reason != game.FivefoldRepetition {
		t.Fatal("history lost draw outcome", err)
	}
}

func quietKnightMoves(t *testing.T, g *game.Game, count int) {
	t.Helper()
	visits := map[chess.Position]int{rules.RepetitionKey(g.Position()): 1}
	for range count {
		position := g.Position()
		var best chess.Move
		fewest := 5
		for from := chess.A1; from <= chess.H8; from++ {
			piece, _ := position.Board.PieceAt(from)
			if piece.Type != chess.Knight || piece.Color != position.SideToMove {
				continue
			}
			moves, err := g.LegalMovesFrom(from)
			if err != nil {
				t.Fatal(err)
			}
			for _, move := range moves {
				victim, _ := position.Board.PieceAt(move.To)
				if victim.Type != chess.NoPiece {
					continue
				}
				next, err := rules.PositionAfterMove(position, move)
				if err != nil {
					t.Fatal(err)
				}
				check, err := rules.IsInCheck(next.Board, next.SideToMove)
				if err != nil {
					t.Fatal(err)
				}
				if check {
					continue
				}
				seen := visits[rules.RepetitionKey(next)]
				if seen < fewest {
					best, fewest = move, seen
				}
			}
		}
		if fewest >= 4 {
			t.Fatal("quiet setup would cause an automatic repetition draw")
		}
		if _, err := g.Play(best); err != nil {
			t.Fatal(err)
		}
		visits[rules.RepetitionKey(g.Position())]++
	}
}

func TestMoveCountDraw(t *testing.T) {
	for _, count := range []int{99, 100, 150} {
		g := game.NewGame()
		quietKnightMoves(t, g, count)
		state := snapshot(t, g)
		if count == 150 {
			if state.Status != chess.StatusSeventyFiveMoves || state.Outcome.Reason != game.SeventyFiveMoves {
				t.Fatal("missing automatic 75-move draw")
			}
			continue
		}
		if state.Status != chess.StatusOngoing {
			t.Fatal("game ended before the automatic threshold")
		}
		if count == 99 {
			if g.ClaimableDraw() != game.NoReason {
				t.Fatal("draw claim available too early")
			}
			continue
		}
		if g.ClaimableDraw() != game.FiftyMoves {
			t.Fatal("missing 50-move claim")
		}
		if err := g.ClaimDraw(); err != nil {
			t.Fatal(err)
		}
		if snapshot(t, g).Outcome.Reason != game.FiftyMoves {
			t.Fatal("missing 50-move outcome")
		}
	}
	g := game.NewGame()
	quietKnightMoves(t, g, 149)
	position := g.Position()
	for from := chess.A1; from <= chess.H8; from++ {
		piece, _ := position.Board.PieceAt(from)
		if piece.Type != chess.Pawn || piece.Color != position.SideToMove {
			continue
		}
		moves, err := g.LegalMovesFrom(from)
		if err != nil {
			t.Fatal(err)
		}
		if len(moves) == 0 {
			continue
		}
		if _, err := g.Play(moves[0]); err != nil {
			t.Fatal(err)
		}
		if g.Position().HalfmoveClock != 0 || snapshot(t, g).Status != chess.StatusOngoing {
			t.Fatal("pawn did not reset clock")
		}
		return
	}
	t.Fatal("setup has no pawn move")
}

func TestMateTakesPrecedenceOverMoveCount(t *testing.T) {
	g := game.NewGame()
	play(t, g, chess.Move{From: chess.F2, To: chess.F3}, chess.Move{From: chess.E7, To: chess.E5}, chess.Move{From: chess.G2, To: chess.G4})
	position := g.Position()
	position.HalfmoveClock = 149
	next, err := rules.PositionAfterMove(position, chess.Move{From: chess.D8, To: chess.H4})
	if err != nil {
		t.Fatal(err)
	}
	status, err := rules.EvaluateStatus(next)
	if err != nil || status != chess.StatusCheckmate {
		t.Fatal("move count overrode checkmate", err)
	}
}
