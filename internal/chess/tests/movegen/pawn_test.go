package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestPawnMoves(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(27))

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(27), To: chess.Square(35)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestPawnMovesBlackForward(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(27))

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(27), To: chess.Square(19)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestPawnMovesBlockedInFront(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	if err := board.SetPiece(chess.Square(35), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 35: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(27))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 move, got %d", len(moves))
	}
}

func TestPawnMovesWhiteInitialAdvance(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(12), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 12: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(12))
	destinations := []chess.Square{20, 28} // e3, e4

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(12), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestPawnMovesWhiteDoubleAdvanceBlocked(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(12), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 12: %v", err)
	}

	if err := board.SetPiece(chess.Square(28), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 28: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(12))
	destinations := []chess.Square{20} // e3

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(12), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestPawnMovesInitialAdvanceBlockedInFront(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(12), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 12: %v", err)
	}

	if err := board.SetPiece(chess.Square(20), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 20: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(12))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 move, got %d", len(moves))
	}
}

func TestPawnMovesBlackInitialAdvance(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(52), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 12: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(52))
	destinations := []chess.Square{44, 36} // e6, e5

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(52), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestPawnMovesWhiteCapturesDiagonally(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	if err := board.SetPiece(chess.Square(35), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 35: %v", err)
	}

	for _, square := range []chess.Square{34, 36} { // c5, e5
		if err := board.SetPiece(square, blockerPawn); err != nil {
			t.Fatalf("Failed to set piece at square %d: %v", square, err)
		}
	}

	moves := movegen.PawnMoves(board, chess.Square(27))
	destinations := []chess.Square{34, 36} // c5, e5

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestPawnMovesBlackCapturesDiagonally(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	if err := board.SetPiece(chess.Square(19), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 19: %v", err)
	}

	for _, square := range []chess.Square{18, 20} { // c3, e3
		if err := board.SetPiece(square, blockerPawn); err != nil {
			t.Fatalf("Failed to set piece at square %d: %v", square, err)
		}
	}

	moves := movegen.PawnMoves(board, chess.Square(27))
	destinations := []chess.Square{18, 20} // c3, e3

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestPawnMovesCannotCaptureOwnPieces(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blockerPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 27: %v", err)
	}

	if err := board.SetPiece(chess.Square(35), blockerPawn); err != nil {
		t.Fatalf("Failed to set piece at square 35: %v", err)
	}

	for _, square := range []chess.Square{34, 36} { // c5, e5
		if err := board.SetPiece(square, blockerPawn); err != nil {
			t.Fatalf("Failed to set piece at square %d: %v", square, err)
		}
	}

	moves := movegen.PawnMoves(board, chess.Square(27))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 moves, got %d", len(moves))
	}
}

func TestPawnMovesFromLeftEdge(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemy := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(24), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 24: %v", err)
	}

	if err := board.SetPiece(chess.Square(31), enemy); err != nil {
		t.Fatalf("Failed to set piece at square 31: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(24))

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(24), To: chess.Square(32)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestPawnMovesFromRightEdge(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemy := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(31), pawn); err != nil {
		t.Fatalf("Failed to set piece at square 31: %v", err)
	}

	if err := board.SetPiece(chess.Square(40), enemy); err != nil {
		t.Fatalf("Failed to set piece at square 40: %v", err)
	}

	moves := movegen.PawnMoves(board, chess.Square(31))

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(31), To: chess.Square(39)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}
