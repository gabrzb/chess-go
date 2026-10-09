package chess

import (
	"errors"
	"testing"
)

func TestLegalMovesPinnedRook(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(12), Piece{Type: Rook, Color: White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: Rook, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := LegalMoves(board, Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []Square{20, 28, 36, 44, 52, 60}

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := Move{From: Square(12), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestLegalMovesKingEscapesCheck(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: Rook, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := LegalMoves(board, Square(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []Square{5, 3, 13, 11}

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := Move{From: Square(4), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestLegalMovesErrors(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(12), Piece{Type: Rook, Color: White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(Square(13), Piece{Type: Rook, Color: NoColor}); err != nil {
		t.Fatalf("Failed to set no color rook: %v", err)
	}

	cases := []struct {
		name    string
		from    Square
		wantErr error
	}{
		{"invalid square", NoSquare, ErrInvalidSquare},
		{"empty square", Square(27), ErrEmptySource},
		{"invalid piece color", Square(13), ErrInvalidColor},
		{"missing king", Square(12), ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LegalMoves(board, tc.from)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestLegalMovesCannotCaptureKing(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(12), Piece{Type: Rook, Color: White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := LegalMoves(board, Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(moves) != 12 {
		t.Fatalf("Expected 12 moves, got %d", len(moves))
	}

	for _, move := range moves {
		if move.To == Square(60) {
			t.Errorf("Expected not to be able to capture the king, but found move: %v", move)
		}
	}
}
