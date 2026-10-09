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

func TestHasLegalMoveAvailable(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	for _, color := range []Color{White, Black} {
		got, err := hasLegalMove(board, color)
		if err != nil {
			t.Fatalf("Unexpected error for color %v: %v", color, err)
		}

		if !got {
			t.Errorf("Expected color %v to have legal moves, but got false", color)
		}
	}
}

func TestHasLegalMoveUnavailable(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(56), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(Square(42), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(49), Piece{Type: Queen, Color: White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := hasLegalMove(board, Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got {
		t.Error("Expected color Black to have no legal moves, but got true")
	}
}

func TestHasLegalMoveContinuesSearching(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(56), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(Square(42), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(49), Piece{Type: Queen, Color: White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	if err := board.SetPiece(Square(57), Piece{Type: Rook, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	got, err := hasLegalMove(board, Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if !got {
		t.Error("Expected color Black to have a legal move, but got false")
	}
}

func TestHasLegalMoveErrors(t *testing.T) {
	var board Board

	cases := []struct {
		name    string
		color   Color
		wantErr error
	}{
		{"invalid color", NoColor, ErrInvalidColor},
		{"missing king", White, ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := hasLegalMove(board, tc.color)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
