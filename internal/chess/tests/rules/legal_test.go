package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestLegalMovesPinnedRook(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := rules.LegalMoves(board, chess.Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{20, 28, 36, 44, 52, 60}

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

func TestLegalMovesKingEscapesCheck(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := rules.LegalMoves(board, chess.Square(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	destinations := []chess.Square{5, 3, 13, 11}

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := chess.Move{From: chess.Square(4), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestLegalMovesErrors(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(13), chess.Piece{Type: chess.Rook, Color: chess.NoColor}); err != nil {
		t.Fatalf("Failed to set no color rook: %v", err)
	}

	cases := []struct {
		name    string
		from    chess.Square
		wantErr error
	}{
		{"invalid square", chess.NoSquare, chess.ErrInvalidSquare},
		{"empty square", chess.Square(27), chess.ErrEmptySource},
		{"invalid piece color", chess.Square(13), chess.ErrInvalidColor},
		{"missing king", chess.Square(12), chess.ErrKingNotFound},
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

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(12), chess.Piece{Type: chess.Rook, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	moves, err := rules.LegalMoves(board, chess.Square(12))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(moves) != 12 {
		t.Fatalf("Expected 12 moves, got %d", len(moves))
	}

	for _, move := range moves {
		if move.To == chess.Square(60) {
			t.Errorf("Expected not to be able to capture the king, but found move: %v", move)
		}
	}
}

func TestHasLegalMoveAvailable(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	for _, color := range []chess.Color{chess.White, chess.Black} {
		got, err := rules.HasLegalMove(board, color)
		if err != nil {
			t.Fatalf("Unexpected error for color %v: %v", color, err)
		}

		if !got {
			t.Errorf("Expected color %v to have legal moves, but got false", color)
		}
	}
}

func TestHasLegalMoveUnavailable(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(56), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(chess.Square(42), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(49), chess.Piece{Type: chess.Queen, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := rules.HasLegalMove(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got {
		t.Error("Expected color Black to have no legal moves, but got true")
	}
}

func TestHasLegalMoveContinuesSearching(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(56), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(chess.Square(42), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(49), chess.Piece{Type: chess.Queen, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	if err := board.SetPiece(chess.Square(57), chess.Piece{Type: chess.Rook, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	got, err := rules.HasLegalMove(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if !got {
		t.Error("Expected color Black to have a legal move, but got false")
	}
}

func TestLegalMovesFromIncludesCastling(t *testing.T) {
	cases := []struct {
		name     string
		color    chess.Color
		kingSide bool
		want     chess.Move
	}{
		{"white king-side", chess.White, true, chess.Move{From: 4, To: 6}},
		{"black king-side", chess.Black, true, chess.Move{From: 60, To: 62}},
		{"white queen-side", chess.White, false, chess.Move{From: 4, To: 2}},
		{"black queen-side", chess.Black, false, chess.Move{From: 60, To: 58}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)

			moves, err := rules.LegalMovesFrom(position, tc.want.From)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			found := false
			for _, move := range moves {
				if move == tc.want {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected castling move %v in %v", tc.want, moves)
			}
		})
	}
}

func TestLegalMovesFromErrors(t *testing.T) {
	colors := []struct {
		name     string
		color    chess.Color
		kingFrom chess.Square
		rookFrom chess.Square
	}{
		{"white", chess.White, 4, 7},
		{"black", chess.Black, 60, 63},
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
				{"empty square", 27, side.color, false, chess.ErrEmptySource},
				{"missing turn color", side.kingFrom, chess.NoColor, false, chess.ErrInvalidColor},
				{"invalid turn color", side.kingFrom, chess.Color(255), false, chess.ErrInvalidColor},
				{"missing king", side.rookFrom, side.color, true, chess.ErrKingNotFound},
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

func TestLegalMovesFromWrongTurn(t *testing.T) {
	cases := []struct {
		name string
		turn chess.Color
		from chess.Square
	}{
		{"white turn selecting black piece", chess.White, 60},
		{"black turn selecting white piece", chess.Black, 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := kingSidePosition(t, tc.turn)

			moves, err := rules.LegalMovesFrom(position, tc.from)
			if !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Expected ErrIllegalMove, got %v", err)
			}
			if len(moves) != 0 {
				t.Errorf("Expected no moves for the opposing side, got %v", moves)
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
			_, err := rules.HasLegalMove(board, tc.color)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
