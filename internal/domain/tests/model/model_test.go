package model_test

import (
	"errors"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestSquareNotationAndBounds(t *testing.T) {
	for square := chess.A1; square <= chess.H8; square++ {
		parsed, err := chess.ParseSquare(square.String())
		if err != nil || parsed != square || !square.Valid() {
			t.Fatalf("Square round trip %v: %v, %v", square, parsed, err)
		}
	}
	for _, input := range []string{"", "a0", "h9", "i1", "a10"} {
		if square, err := chess.ParseSquare(input); square != chess.NoSquare || !errors.Is(err, chess.ErrInvalidSquare) {
			t.Errorf("Invalid square %q: %v, %v", input, square, err)
		}
	}
	for _, tc := range []struct {
		from   chess.Square
		dx, dy int
		want   chess.Square
	}{
		{chess.D4, 1, 1, chess.E5}, {chess.H8, 1, 0, chess.NoSquare},
		{chess.A1, -1, 0, chess.NoSquare}, {chess.H8, 0, 1, chess.NoSquare},
		{chess.A1, 0, -1, chess.NoSquare}, {chess.NoSquare, 0, 0, chess.NoSquare},
	} {
		if got, ok := tc.from.Offset(tc.dx, tc.dy); got != tc.want || ok != tc.want.Valid() {
			t.Errorf("Offset(%v, %d, %d) = %v, %v", tc.from, tc.dx, tc.dy, got, ok)
		}
	}
}

func TestPositionInitialization(t *testing.T) {
	p := chess.NewInitialPosition()
	all := chess.WhiteKingSide | chess.WhiteQueenSide | chess.BlackKingSide | chess.BlackQueenSide
	if p.SideToMove != chess.White || p.CastlingRights != all || p.EnPassant != chess.NoSquare || p.HalfmoveClock != 0 || p.FullmoveNumber != 1 {
		t.Fatal("Incorrect initial position metadata")
	}
	backRank := []chess.PieceType{chess.Rook, chess.Knight, chess.Bishop, chess.Queen, chess.King, chess.Bishop, chess.Knight, chess.Rook}
	for square := chess.A1; square <= chess.H8; square++ {
		want := chess.Piece{}
		if square.Row() == chess.A1.Row() || square.Row() == chess.A8.Row() {
			want.Type = backRank[square.Column()]
		} else if square.Row() == chess.A2.Row() || square.Row() == chess.A7.Row() {
			want.Type = chess.Pawn
		}
		if want.Type != chess.NoPiece {
			want.Color = chess.White
			if square.Row() >= chess.A7.Row() {
				want.Color = chess.Black
			}
		}
		if got, _ := p.Board.PieceAt(square); got != want {
			t.Errorf("At %v: got %v, want %v", square, got, want)
		}
	}
	custom, err := chess.NewPosition(p.Board, chess.Black)
	if err != nil || custom.Board != p.Board || custom.SideToMove != chess.Black || custom.CastlingRights != 0 || custom.EnPassant != chess.NoSquare || custom.FullmoveNumber != 1 {
		t.Fatalf("Incorrect custom position: %+v, %v", custom, err)
	}
	if _, err := chess.NewPosition(p.Board, chess.NoColor); !errors.Is(err, chess.ErrInvalidColor) {
		t.Error("Expected invalid color")
	}
}
