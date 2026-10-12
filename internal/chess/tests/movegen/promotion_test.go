package movegen_test

import (
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

type promotionCase struct {
	name                   string
	color, enemy           chess.Color
	from                   chess.Square
	destinations, captures []chess.Square
}

var promotionCases = []promotionCase{
	{"white advance", chess.White, chess.Black, chess.A7, []chess.Square{chess.A8}, nil},
	{"black advance", chess.Black, chess.White, chess.A2, []chess.Square{chess.A1}, nil},
	{"white captures a8", chess.White, chess.Black, chess.B7, []chess.Square{chess.B8, chess.A8}, []chess.Square{chess.A8}},
	{"white captures h8", chess.White, chess.Black, chess.G7, []chess.Square{chess.G8, chess.H8}, []chess.Square{chess.H8}},
	{"black captures a1", chess.Black, chess.White, chess.B2, []chess.Square{chess.B1, chess.A1}, []chess.Square{chess.A1}},
	{"black captures h1", chess.Black, chess.White, chess.G2, []chess.Square{chess.G1, chess.H1}, []chess.Square{chess.H1}},
}

func expectedPromotions(from, to chess.Square) []chess.Move {
	return []chess.Move{
		{From: from, To: to, Promotion: chess.Queen},
		{From: from, To: to, Promotion: chess.Rook},
		{From: from, To: to, Promotion: chess.Bishop},
		{From: from, To: to, Promotion: chess.Knight},
	}
}

func TestPawnPromotions(t *testing.T) {
	for _, tc := range promotionCases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, tc.from, chess.Piece{Type: chess.Pawn, Color: tc.color})
			for _, square := range tc.captures {
				mustSetPiece(t, &board, square, chess.Piece{Type: chess.Rook, Color: tc.enemy})
			}
			before := board
			var want []chess.Move
			for _, to := range tc.destinations {
				want = append(want, expectedPromotions(tc.from, to)...)
			}
			for _, generate := range []func(chess.Board, chess.Square) []chess.Move{movegen.PawnMoves, movegen.PseudoLegalMoves} {
				if got := generate(board, tc.from); !slices.Equal(got, want) {
					t.Errorf("Expected promotions %v, got %v", want, got)
				}
			}
			if board != before {
				t.Error("Expected generation to preserve original board")
			}
		})
	}
}

func TestPawnPromotionBlocked(t *testing.T) {
	sides := []struct {
		name         string
		color, enemy chess.Color
		from, to     chess.Square
	}{
		{"white", chess.White, chess.Black, chess.A7, chess.A8},
		{"black", chess.Black, chess.White, chess.A2, chess.A1},
	}
	for _, side := range sides {
		t.Run(side.name, func(t *testing.T) {
			cases := []struct {
				name         string
				occupant     chess.Color
				alreadyAtEnd bool
			}{
				{"own piece in front", side.color, false},
				{"enemy piece in front", side.enemy, false},
				{"pawn already at last row", chess.NoColor, true},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					var board chess.Board
					from := side.from
					if tc.alreadyAtEnd {
						from = side.to
					} else {
						mustSetPiece(t, &board, side.to, chess.Piece{Type: chess.Rook, Color: tc.occupant})
					}
					mustSetPiece(t, &board, from, chess.Piece{Type: chess.Pawn, Color: side.color})
					if got := movegen.PawnMoves(board, from); len(got) != 0 {
						t.Errorf("Expected no promotion moves, got %v", got)
					}
				})
			}
		})
	}
}

func TestPawnPromotionCannotCaptureOwnPiece(t *testing.T) {
	for _, tc := range promotionCases {
		if len(tc.captures) == 0 {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, tc.from, chess.Piece{Type: chess.Pawn, Color: tc.color})
			for _, square := range tc.captures {
				mustSetPiece(t, &board, square, chess.Piece{Type: chess.Rook, Color: tc.color})
			}
			want := expectedPromotions(tc.from, tc.destinations[0])
			if got := movegen.PawnMoves(board, tc.from); !slices.Equal(got, want) {
				t.Errorf("Expected only forward promotions %v, got %v", want, got)
			}
		})
	}
}
