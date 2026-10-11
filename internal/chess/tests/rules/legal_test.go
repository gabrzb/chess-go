package rules_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestLegalMovesPinnedRook(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{20, 28, 36, 44, 52, 60}

	assertMoves(t, moves, 12, destinations)
}

func TestLegalMovesKingEscapesCheck(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.Square(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{5, 3, 13, 11}

	assertMoves(t, moves, 4, destinations)
}

func TestLegalMovesErrors(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(13), chess.Piece{Type: chess.Rook, Color: chess.NoColor})

	cases := []struct {
		name    string
		from    chess.Square
		wantErr error
	}{
		{"invalid square", chess.NoSquare, chess.ErrInvalidSquare},
		{"empty square", chess.Square(27), chess.ErrEmptySource},
		{"invalid piece color", chess.Square(13), chess.ErrInvalidColor},
		{"missing king", chess.Square(12), chess.ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rules.LegalMoves(board, tc.from)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestLegalMovesCannotCaptureKing(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(moves) != 12 {
		t.Fatalf("Expected 12 moves, got %d", len(moves))
	}

	for _, move := range moves {
		if move.To == chess.Square(60) {
			t.Errorf("Expected not to be able to capture the king, but found move: %v", move)
		}
	}
}

func TestHasLegalMoveAvailable(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black})

	for _, color := range []chess.Color{chess.White, chess.Black} {
		got, err := rules.HasLegalMove(board, color)
		if err != nil {
			t.Fatalf("Unexpected error for color %v: %v", color, err)
		}

		if !got {
			t.Errorf("Expected color %v to have legal moves, but got false", color)
		}
	}
}

func TestHasLegalMoveSearch(t *testing.T) {
	for _, tc := range []struct {
		name    string
		addRook bool
		want    bool
	}{
		{"no legal moves", false, false},
		{"continue to next piece", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, 56, chess.Piece{Type: chess.King, Color: chess.Black})
			mustSetPiece(t, &board, 42, chess.Piece{Type: chess.King, Color: chess.White})
			mustSetPiece(t, &board, 49, chess.Piece{Type: chess.Queen, Color: chess.White})
			if tc.addRook {
				mustSetPiece(t, &board, 57, chess.Piece{Type: chess.Rook, Color: chess.Black})
			}
			got, err := rules.HasLegalMove(board, chess.Black)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Expected legal-move availability %v, got %v", tc.want, got)
			}
		})
	}
}

func TestLegalMovesFromIncludesCastling(t *testing.T) {
	for _, castle := range castlingCases {
		t.Run(castle.name, func(t *testing.T) {
			position := castlingPosition(t, castle.color, castle.kingSide)
			moves, err := rules.LegalMovesFrom(position, castle.kingMove.From)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !slices.Contains(moves, castle.kingMove) {
				t.Errorf("Expected castling move %v in %v", castle.kingMove, moves)
			}
		})
	}
}

func TestLegalMovesFromErrors(t *testing.T) {
	colors := []struct {
		name      string
		color     chess.Color
		kingFrom  chess.Square
		rookFrom  chess.Square
		enemyKing chess.Square
	}{
		{"white", chess.White, 4, 7, 60},
		{"black", chess.Black, 60, 63, 4},
	}

	for _, side := range colors {
		t.Run(side.name, func(t *testing.T) {
			cases := []struct {
				name       string
				from       chess.Square
				turn       chess.Color
				removeKing bool
				wantErr    error
			}{
				{"invalid square", chess.NoSquare, side.color, false, chess.ErrInvalidSquare},
				{"out of range square", chess.Square(255), side.color, false, chess.ErrInvalidSquare},
				{"empty square", 27, side.color, false, chess.ErrEmptySource},
				{"missing turn color", side.kingFrom, chess.NoColor, false, chess.ErrInvalidColor},
				{"invalid turn color", side.kingFrom, chess.Color(255), false, chess.ErrInvalidColor},
				{"missing king", side.rookFrom, side.color, true, chess.ErrKingNotFound},
				{"opposing piece", side.enemyKing, side.color, false, chess.ErrIllegalMove},
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					position := kingSidePosition(t, side.color)
					position.SideToMove = tc.turn
					if tc.removeKing {
						if _, err := position.Board.RemovePiece(side.kingFrom); err != nil {
							t.Fatalf("Failed to remove king: %v", err)
						}
					}

					moves, err := rules.LegalMovesFrom(position, tc.from)
					if !errors.Is(err, tc.wantErr) {
						t.Errorf("Expected error %v, got %v", tc.wantErr, err)
					}
					if len(moves) != 0 {
						t.Errorf("Expected no moves on error, got %v", moves)
					}
				})
			}
		})
	}
}

func TestHasLegalMoveErrors(t *testing.T) {
	var board chess.Board

	cases := []struct {
		name    string
		color   chess.Color
		wantErr error
	}{
		{"invalid color", chess.NoColor, chess.ErrInvalidColor},
		{"missing king", chess.White, chess.ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rules.HasLegalMove(board, tc.color)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
