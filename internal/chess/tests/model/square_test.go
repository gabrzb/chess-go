package model_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestSquareCoordinates(t *testing.T) {
	cases := []struct {
		name        string
		square      chess.Square
		column, row int
	}{
		{"a1", chess.A1, 0, 0},
		{"h1", chess.H1, chess.BoardSize - 1, 0},
		{"e4", chess.E4, 4, 3},
		{"a8", chess.A8, 0, chess.BoardSize - 1},
		{"h8", chess.H8, chess.BoardSize - 1, chess.BoardSize - 1},
		{"no square", chess.NoSquare, -1, -1},
		{"out of range", chess.Square(255), -1, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if col, row := tc.square.Column(), tc.square.Row(); col != tc.column || row != tc.row {
				t.Errorf("Expected coordinates (%d, %d), got (%d, %d)", tc.column, tc.row, col, row)
			}
		})
	}
}

func TestSquareOffset(t *testing.T) {
	lastIndex := chess.BoardSize - 1
	cases := []struct {
		name string
		from chess.Square
		step chess.Direction
		want chess.Square
		ok   bool
	}{
		{"diagonal", chess.D4, chess.Direction{ColStep: 1, RowStep: 1}, chess.E5, true},
		{"knight jump", chess.D4, chess.Direction{ColStep: 2, RowStep: 1}, chess.F5, true},
		{"reverse across board", chess.H8, chess.Direction{ColStep: -lastIndex, RowStep: -lastIndex}, chess.A1, true},
		{"zero offset", chess.D4, chess.Direction{}, chess.D4, true},
		{"left edge", chess.A1, chess.Direction{ColStep: -1}, chess.NoSquare, false},
		{"right edge", chess.H1, chess.Direction{ColStep: 1}, chess.NoSquare, false},
		{"bottom edge", chess.A1, chess.Direction{RowStep: -1}, chess.NoSquare, false},
		{"top edge", chess.A8, chess.Direction{RowStep: 1}, chess.NoSquare, false},
		{"large offset", chess.A1, chess.Direction{ColStep: chess.BoardSize}, chess.NoSquare, false},
		{"invalid origin", chess.NoSquare, chess.Direction{}, chess.NoSquare, false},
		{"out of range origin", chess.Square(255), chess.Direction{}, chess.NoSquare, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.from.Offset(tc.step.ColStep, tc.step.RowStep)
			if got != tc.want || ok != tc.ok {
				t.Errorf("Expected (%v, %v), got (%v, %v)", tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestNamedSquares(t *testing.T) {
	rows := [][]chess.Square{
		{chess.A1, chess.B1, chess.C1, chess.D1, chess.E1, chess.F1, chess.G1, chess.H1},
		{chess.A2, chess.B2, chess.C2, chess.D2, chess.E2, chess.F2, chess.G2, chess.H2},
		{chess.A3, chess.B3, chess.C3, chess.D3, chess.E3, chess.F3, chess.G3, chess.H3},
		{chess.A4, chess.B4, chess.C4, chess.D4, chess.E4, chess.F4, chess.G4, chess.H4},
		{chess.A5, chess.B5, chess.C5, chess.D5, chess.E5, chess.F5, chess.G5, chess.H5},
		{chess.A6, chess.B6, chess.C6, chess.D6, chess.E6, chess.F6, chess.G6, chess.H6},
		{chess.A7, chess.B7, chess.C7, chess.D7, chess.E7, chess.F7, chess.G7, chess.H7},
		{chess.A8, chess.B8, chess.C8, chess.D8, chess.E8, chess.F8, chess.G8, chess.H8},
	}
	files, ranks := "abcdefgh", "12345678"
	for row, squares := range rows {
		for col, square := range squares {
			want := string([]byte{files[col], ranks[row]})
			t.Run(want, func(t *testing.T) {
				if got := square.String(); got != want {
					t.Errorf("Expected named square %s, got %s", want, got)
				}
				got, err := chess.ParseSquare(want)
				if err != nil || got != square {
					t.Errorf("ParseSquare(%s) = %v, %v; expected %v", want, got, err, square)
				}
			})
		}
	}
}
