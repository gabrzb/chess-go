package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestBoardAfterMove(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	beforeMove := board

	nextMove, err := rules.BoardAfterMove(board, chess.Move{From: chess.Square(27), To: chess.Square(35)})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	source, _ := nextMove.PieceAt(chess.Square(27))
	if source.Type != chess.NoPiece {
		t.Errorf("Expected source square to be empty, got %v", source)
	}

	destination, _ := nextMove.PieceAt(chess.Square(35))
	if destination != rook {
		t.Errorf("Expected destination square to have rook, got %v", destination)
	}

	if board != beforeMove {
		t.Errorf("Expected original board to remain unchanged, but it changed")
	}
}

func TestBoardAfterMoveCapture(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(35), pawn); err != nil {
		t.Fatalf("Failed to place pawn: %v", err)
	}

	beforeMove := board

	nextMove, err := rules.BoardAfterMove(board, chess.Move{From: chess.Square(27), To: chess.Square(35)})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	source, _ := nextMove.PieceAt(chess.Square(27))
	if source.Type != chess.NoPiece {
		t.Errorf("Expected source square to be empty, got %v", source)
	}

	destination, _ := nextMove.PieceAt(chess.Square(35))
	if destination != rook {
		t.Errorf("Expected destination square to have rook, got %v", destination)
	}

	if board != beforeMove {
		t.Errorf("Expected original board to remain unchanged, but it changed")
	}
}

func TestBoardAfterMoveErrors(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(27), chess.Piece{Type: chess.Rook, Color: chess.White}); err != nil {
		t.Fatalf("Failed to place rook: %v", err)
	}

	cases := []struct {
		name    string
		move    chess.Move
		wantErr error
	}{
		{"invalid source", chess.Move{From: chess.NoSquare, To: chess.Square(35)}, chess.ErrInvalidSquare},
		{"invalid destination", chess.Move{From: chess.Square(27), To: chess.NoSquare}, chess.ErrInvalidSquare},
		{"empty source", chess.Move{From: chess.Square(28), To: chess.Square(35)}, chess.ErrEmptySource},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rules.BoardAfterMove(board, tc.move)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
