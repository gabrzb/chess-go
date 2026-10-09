package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestRookMovesToRightBlockedByOwnPiece(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.SlidingMoves(
		board,
		chess.Square(27),
		chess.Direction{ColStep: 1, RowStep: 0},
	)

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(27), To: chess.Square(28)}
	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestRookMovesToRightCapturesEnemy(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.SlidingMoves(
		board,
		chess.Square(27),
		chess.Direction{ColStep: 1, RowStep: 0},
	)

	expectedMoves := []chess.Move{
		{From: chess.Square(27), To: chess.Square(28)},
		{From: chess.Square(27), To: chess.Square(29)},
	}

	if len(moves) != len(expectedMoves) {
		t.Fatalf("Expected %d moves, got %d", len(expectedMoves), len(moves))
	}

	for i, wantedMove := range expectedMoves {
		if moves[i] != wantedMove {
			t.Errorf("Expected move %v, got %v", wantedMove, moves[i])
		}
	}
}

func TestRookMovesFromEmptySquare(t *testing.T) {
	var board chess.Board
	moves := movegen.RookMoves(board, chess.Square(27))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 moves, got %d", len(moves))
	}
}

func TestRookMoves(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.RookMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28, 29, 30, 31, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
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

func TestBishopMoves(t *testing.T) {
	var board chess.Board
	bishop := chess.Piece{Type: chess.Bishop, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.BishopMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		36, 45, 54, 63, // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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

func TestBishopMovesBlockedByOwnPiece(t *testing.T) {
	var board chess.Board
	bishop := chess.Piece{Type: chess.Bishop, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(45), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.BishopMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		36,         // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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

func TestBishopCapturesEnemy(t *testing.T) {
	var board chess.Board
	bishop := chess.Piece{Type: chess.Bishop, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(45), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.BishopMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		36, 45, // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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

func TestQueenMoves(t *testing.T) {
	var board chess.Board
	queen := chess.Piece{Type: chess.Queen, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.QueenMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28, 29, 30, 31, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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

func TestQueenMovesBlockedByOwnPiece(t *testing.T) {
	var board chess.Board
	queen := chess.Piece{Type: chess.Queen, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.QueenMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28,         // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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

func TestQueenMovesCapturesEnemy(t *testing.T) {
	var board chess.Board
	queen := chess.Piece{Type: chess.Queen, Color: chess.White}
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.QueenMoves(board, chess.Square(27))

	destinationSquares := []chess.Square{
		28, 29, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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
