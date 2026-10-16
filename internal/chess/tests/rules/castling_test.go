package rules_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

var castles = []struct {
	turn       string
	king, rook chess.Move
}{
	{"w", chess.Move{From: chess.E1, To: chess.G1}, chess.Move{From: chess.H1, To: chess.F1}},
	{"w", chess.Move{From: chess.E1, To: chess.C1}, chess.Move{From: chess.A1, To: chess.D1}},
	{"b", chess.Move{From: chess.E8, To: chess.G8}, chess.Move{From: chess.H8, To: chess.F8}},
	{"b", chess.Move{From: chess.E8, To: chess.C8}, chess.Move{From: chess.A8, To: chess.D8}},
}

func TestCastlingExecution(t *testing.T) {
	for _, tc := range castles {
		p := perftPositionFromFEN(t, fmt.Sprintf("r3k2r/8/8/8/8/8/8/R3K2R %s KQkq - 7 20", tc.turn))
		next := after(t, p, tc.king)
		assertPiece(t, next, tc.king.From, chess.Piece{})
		assertPiece(t, next, tc.rook.From, chess.Piece{})
		assertPiece(t, next, tc.king.To, chess.Piece{Type: chess.King, Color: p.SideToMove})
		assertPiece(t, next, tc.rook.To, chess.Piece{Type: chess.Rook, Color: p.SideToMove})
		wantRights, full := chess.BlackKingSide|chess.BlackQueenSide, p.FullmoveNumber
		if p.SideToMove == chess.Black {
			wantRights, full = chess.WhiteKingSide|chess.WhiteQueenSide, full+1
		}
		if next.CastlingRights != wantRights || next.SideToMove == p.SideToMove || next.EnPassant != chess.NoSquare || next.HalfmoveClock != p.HalfmoveClock+1 || next.FullmoveNumber != full {
			t.Errorf("Incorrect state after %v: %+v", tc.king, next)
		}
	}
}

func TestCastlingRestrictions(t *testing.T) {
	for _, tc := range castles {
		base := perftPositionFromFEN(t, fmt.Sprintf("r3k2r/8/8/8/8/8/8/R3K2R %s KQkq - 7 20", tc.turn))
		step := 1
		if tc.king.To.Column() < tc.king.From.Column() {
			step = -1
		}
		transit, _ := tc.king.From.Offset(step, 0)
		for _, condition := range []string{"no rights", "missing rook", "blocked", "in check", "through check", "into check", "queen-side b file blocked"} {
			p := base
			switch condition {
			case "no rights":
				p.CastlingRights = 0
			case "missing rook":
				mustSetPiece(t, &p.Board, tc.rook.From, chess.Piece{})
			case "blocked":
				mustSetPiece(t, &p.Board, transit, chess.Piece{Type: chess.Knight, Color: p.SideToMove})
			case "queen-side b file blocked":
				if step > 0 {
					continue
				}
				square, _ := tc.king.From.Offset(-3, 0)
				mustSetPiece(t, &p.Board, square, chess.Piece{Type: chess.Knight, Color: p.SideToMove})
			default:
				target := tc.king.From
				if condition == "through check" {
					target = transit
				}
				if condition == "into check" {
					target = tc.king.To
				}
				rowStep, enemy := 1, chess.Black
				if p.SideToMove == chess.Black {
					rowStep, enemy = -1, chess.White
				}
				attacker, _ := target.Offset(0, rowStep)
				mustSetPiece(t, &p.Board, attacker, chess.Piece{Type: chess.Rook, Color: enemy})
			}
			if _, err := rules.PositionAfterMove(p, tc.king); !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Allowed %v with %s: %v", tc.king, condition, err)
			}
		}
	}
}
