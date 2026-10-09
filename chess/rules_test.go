package chess

import (
	"errors"
	"testing"
)

func TestBoardAfterMove(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	beforeMove := board

	nextMove, err := boardAfterMove(board, Move{From: Square(27), To: Square(35)})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	source, _ := nextMove.PieceAt(Square(27))
	if source.Type != NoPiece {
		t.Errorf("Expected source square to be empty, got %v", source)
	}

	destination, _ := nextMove.PieceAt(Square(35))
	if destination != rook {
		t.Errorf("Expected destination square to have rook, got %v", destination)
	}

	if board != beforeMove {
		t.Errorf("Expected original board to remain unchanged, but it changed")
	}
}

func TestBoardAfterMoveCapture(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	if err := board.SetPiece(Square(35), pawn); err != nil {
		t.Fatalf("Failed to place pawn: %v", err)
	}

	beforeMove := board

	nextMove, err := boardAfterMove(board, Move{From: Square(27), To: Square(35)})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	source, _ := nextMove.PieceAt(Square(27))
	if source.Type != NoPiece {
		t.Errorf("Expected source square to be empty, got %v", source)
	}

	destination, _ := nextMove.PieceAt(Square(35))
	if destination != rook {
		t.Errorf("Expected destination square to have rook, got %v", destination)
	}

	if board != beforeMove {
		t.Errorf("Expected original board to remain unchanged, but it changed")
	}
}

func TestBoardAfterMoveErrors(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(27), Piece{Type: Rook, Color: White}); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	cases := []struct {
		name    string
		move    Move
		wantErr error
	}{
		{"invalid source", Move{From: NoSquare, To: Square(35)}, ErrInvalidSquare},
		{"invalid destination", Move{From: Square(27), To: NoSquare}, ErrInvalidSquare},
		{"empty source", Move{From: Square(28), To: Square(35)}, ErrEmptySource},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := boardAfterMove(board, tc.move)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
