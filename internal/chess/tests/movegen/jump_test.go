package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestKnightMoves(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KnightMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKnightMovesFromCorner(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}

	if err := board.SetPiece(chess.Square(0), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KnightMoves(board, chess.Square(0))

	destinationSquares := []chess.Square{10, 17} // c2, b3

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(0), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKnightMovesJumpOverPieces(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	for _, square := range []chess.Square{18, 19, 20, 26, 28, 34, 35, 36} {
		if err := board.SetPiece(square, pawn); err != nil {
			t.Fatalf("Failed to set piece: %v", err)
		}
	}

	moves := movegen.KnightMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKnightMovesBlockedByOwnPiece(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(37), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KnightMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		21,     // f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKnightMovesCapturesEnemy(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(37), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KnightMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMoves(t *testing.T) {
	var board chess.Board
	king := chess.Piece{Type: chess.King, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KingMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28, 26, // e4, c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesFromCorner(t *testing.T) {
	var board chess.Board
	king := chess.Piece{Type: chess.King, Color: chess.White}

	if err := board.SetPiece(chess.Square(0), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KingMoves(board, chess.Square(0))

	destinationSquares := []chess.Square{1, 8, 9} // b1, a2, b2

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(0), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesBlockedByOwnPiece(t *testing.T) {
	var board chess.Board
	king := chess.Piece{Type: chess.King, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(28), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KingMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		26,     // c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesCapturesEnemy(t *testing.T) {
	var board chess.Board
	king := chess.Piece{Type: chess.King, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(28), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.KingMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28, 26, // e4, c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := chess.Move{From: chess.Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}
