package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func mustSetPiece(t *testing.T, board *chess.Board, square chess.Square, piece chess.Piece) {
	t.Helper()
	if err := board.SetPiece(square, piece); err != nil {
		t.Fatalf("Failed to set piece at %v: %v", square, err)
	}
}

func assertMoves(t *testing.T, moves []chess.Move, from chess.Square, destinations []chess.Square) {
	t.Helper()
	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}
	for i, to := range destinations {
		want := chess.Move{From: from, To: to}
		if moves[i] != want {
			t.Errorf("Move %d: expected %v, got %v", i, want, moves[i])
		}
	}
}

type positionState struct {
	turn               chess.Color
	rights             chess.CastlingRights
	enPassant          chess.Square
	halfmove, fullmove uint16
}

func assertPositionState(t *testing.T, position chess.Position, want positionState) {
	t.Helper()
	got := positionState{position.SideToMove, position.CastlingRights, position.EnPassant, position.HalfmoveClock, position.FullmoveNumber}
	if got != want {
		t.Errorf("Expected state %+v, got %+v", want, got)
	}
}

func assertPositionPiece(t *testing.T, position chess.Position, square chess.Square, want chess.Piece) {
	t.Helper()
	got, ok := position.Board.PieceAt(square)
	if !ok {
		t.Fatalf("Invalid square in test: %v", square)
	}
	if got != want {
		t.Errorf("At %v: expected %v, got %v", square, want, got)
	}
}

func castlingPosition(t *testing.T, color chess.Color, kingSide bool) chess.Position {
	t.Helper()

	var board chess.Board
	rookSquare := chess.A1
	right := chess.WhiteQueenSide
	if kingSide {
		rookSquare = chess.H1
		right = chess.WhiteKingSide
	}

	if color == chess.Black {
		rookSquare = chess.A8
		right = chess.BlackQueenSide
		if kingSide {
			right = chess.BlackKingSide
			rookSquare = chess.H8
		}
	}

	pieces := []struct {
		square chess.Square
		piece  chess.Piece
	}{
		{chess.E1, chess.Piece{Type: chess.King, Color: chess.White}},
		{rookSquare, chess.Piece{Type: chess.Rook, Color: color}},
		{chess.E8, chess.Piece{Type: chess.King, Color: chess.Black}},
	}

	for _, item := range pieces {
		mustSetPiece(t, &board, item.square, item.piece)
	}

	position, err := chess.NewPosition(board, color)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	position.CastlingRights = right

	return position
}

func kingSidePosition(t *testing.T, color chess.Color) chess.Position {
	t.Helper()
	return castlingPosition(t, color, true)
}
