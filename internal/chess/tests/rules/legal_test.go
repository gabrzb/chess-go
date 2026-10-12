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

	mustSetPiece(t, &board, chess.E1, chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.E2, chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.E8, chess.Piece{Type: chess.Rook, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.E2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{chess.E3, chess.E4, chess.E5, chess.E6, chess.E7, chess.E8}

	assertMoves(t, moves, chess.E2, destinations)
}

func TestLegalMovesKingEscapesCheck(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.E1, chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.E8, chess.Piece{Type: chess.Rook, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.E1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{chess.F1, chess.D1, chess.F2, chess.D2}

	assertMoves(t, moves, chess.E1, destinations)
}

func TestLegalMovesErrors(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.E2, chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.F2, chess.Piece{Type: chess.Rook, Color: chess.NoColor})

	cases := []struct {
		name    string
		from    chess.Square
		wantErr error
	}{
		{"invalid square", chess.NoSquare, chess.ErrInvalidSquare},
		{"empty square", chess.D4, chess.ErrEmptySource},
		{"invalid piece color", chess.F2, chess.ErrInvalidColor},
		{"missing king", chess.E2, chess.ErrKingNotFound},
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

	mustSetPiece(t, &board, chess.E1, chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.E2, chess.Piece{Type: chess.Rook, Color: chess.White})

	mustSetPiece(t, &board, chess.E8, chess.Piece{Type: chess.King, Color: chess.Black})

	moves, err := rules.LegalMoves(board, chess.E2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(moves) != 12 {
		t.Fatalf("Expected 12 moves, got %d", len(moves))
	}

	for _, move := range moves {
		if move.To == chess.E8 {
			t.Errorf("Expected not to be able to capture the king, but found move: %v", move)
		}
	}
}

func TestHasLegalMoveAvailable(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.E1, chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.E8, chess.Piece{Type: chess.King, Color: chess.Black})

	for _, color := range []chess.Color{chess.White, chess.Black} {
		position := chess.Position{Board: board, SideToMove: color, EnPassant: chess.NoSquare}
		got, err := rules.HasLegalMove(position)
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
			mustSetPiece(t, &board, chess.A8, chess.Piece{Type: chess.King, Color: chess.Black})
			mustSetPiece(t, &board, chess.C6, chess.Piece{Type: chess.King, Color: chess.White})
			mustSetPiece(t, &board, chess.B7, chess.Piece{Type: chess.Queen, Color: chess.White})
			if tc.addRook {
				mustSetPiece(t, &board, chess.B8, chess.Piece{Type: chess.Rook, Color: chess.Black})
			}
			position := chess.Position{Board: board, SideToMove: chess.Black, EnPassant: chess.NoSquare}
			got, err := rules.HasLegalMove(position)
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
		{"white", chess.White, chess.E1, chess.H1, chess.E8},
		{"black", chess.Black, chess.E8, chess.H8, chess.E1},
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
				{"empty square", chess.D4, side.color, false, chess.ErrEmptySource},
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
			position := chess.Position{Board: board, SideToMove: tc.color, EnPassant: chess.NoSquare}
			_, err := rules.HasLegalMove(position)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
