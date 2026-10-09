package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func checkSlidingMoves(t *testing.T, dir chess.Direction, destinations []chess.Square) {
	t.Helper()

	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.SlidingMoves(
		board,
		chess.Square(27),
		dir,
	)

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

func TestSlidingMovesToRight(t *testing.T) {
	checkSlidingMoves(t,
		chess.Direction{ColStep: 1, RowStep: 0},
		[]chess.Square{28, 29, 30, 31},
	)
}

func TestSlidingMovesToLeft(t *testing.T) {
	checkSlidingMoves(t,
		chess.Direction{ColStep: -1, RowStep: 0},
		[]chess.Square{26, 25, 24},
	)
}

func TestSlidingMovesUp(t *testing.T) {
	checkSlidingMoves(t,
		chess.Direction{ColStep: 0, RowStep: 1},
		[]chess.Square{35, 43, 51, 59},
	)
}

func TestSlidingMovesDown(t *testing.T) {
	checkSlidingMoves(t,
		chess.Direction{ColStep: 0, RowStep: -1},
		[]chess.Square{19, 11, 3},
	)
}

func TestCanOccupySquareEmpty(t *testing.T) {
	var board chess.Board

	if !movegen.CanOccupySquare(board, chess.Square(27), chess.White) {
		t.Errorf("Expected canOccupySquare to return true for empty square")
	}
}

func TestCanOccupySquareInvalid(t *testing.T) {
	var board chess.Board

	if movegen.CanOccupySquare(board, chess.NoSquare, chess.White) {
		t.Errorf("Expected canOccupySquare to return false for NoSquare")
	}
}

func TestCanOccupySquareOwnPiece(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if movegen.CanOccupySquare(board, chess.Square(27), chess.White) {
		t.Errorf("Expected canOccupySquare to return false for own piece")
	}
}

func TestCanOccupySquareEnemyPiece(t *testing.T) {
	var board chess.Board
	pawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if !movegen.CanOccupySquare(board, chess.Square(27), chess.White) {
		t.Errorf("Expected canOccupySquare to return true for enemy piece")
	}
}

func TestJumpMoves(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := movegen.JumpMoves(
		board,
		chess.Square(27),
		chess.Direction{ColStep: 2, RowStep: 1},
	)

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := chess.Move{From: chess.Square(27), To: chess.Square(37)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestPseudoLegalMovesFromEmptySquare(t *testing.T) {
	var board chess.Board

	moves := movegen.PseudoLegalMoves(board, chess.Square(27))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 moves, got %d", len(moves))
	}
}

func TestPseudoLegalMovesFromInvalidSquare(t *testing.T) {
	var board chess.Board

	moves := movegen.PseudoLegalMoves(board, chess.NoSquare)

	if len(moves) != 0 {
		t.Fatalf("Expected 0 moves, got %d", len(moves))
	}
}

func TestPseudoLegalMovesByPiece(t *testing.T) {
	cases := []struct {
		name         string
		pieceType    chess.PieceType
		destinations []chess.Square
	}{
		{"pawn", chess.Pawn, []chess.Square{35}},
		{"knight", chess.Knight, []chess.Square{37, 21, 33, 17, 44, 42, 12, 10}},
		{"bishop", chess.Bishop, []chess.Square{36, 45, 54, 63, 34, 41, 48, 20, 13, 6, 18, 9, 0}},
		{"rook", chess.Rook, []chess.Square{28, 29, 30, 31, 26, 25, 24, 35, 43, 51, 59, 19, 11, 3}},
		{"queen", chess.Queen, []chess.Square{
			28, 29, 30, 31, 26, 25, 24,
			35, 43, 51, 59, 19, 11, 3,
			36, 45, 54, 63, 34, 41, 48,
			20, 13, 6, 18, 9, 0,
		}},
		{"king", chess.King, []chess.Square{28, 26, 35, 19, 36, 34, 20, 18}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			piece := chess.Piece{Type: tc.pieceType, Color: chess.White}

			if err := board.SetPiece(chess.Square(27), piece); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			moves := movegen.PseudoLegalMoves(board, chess.Square(27))

			if len(moves) != len(tc.destinations) {
				t.Fatalf("Expected %d moves, got %d", len(tc.destinations), len(moves))
			}

			for i, destination := range tc.destinations {
				expectedMove := chess.Move{From: chess.Square(27), To: destination}

				if moves[i] != expectedMove {
					t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
				}
			}
		})
	}
}
