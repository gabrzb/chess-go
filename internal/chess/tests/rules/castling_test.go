package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func kingSidePosition(t *testing.T, color chess.Color) chess.Position {
	t.Helper()

	var board chess.Board
	rookSquare := chess.Square(7)
	right := chess.WhiteKingSide
	if color == chess.Black {
		rookSquare = 63
		right = chess.BlackKingSide
	}

	pieces := []struct {
		square chess.Square
		piece  chess.Piece
	}{
		{4, chess.Piece{Type: chess.King, Color: chess.White}},
		{rookSquare, chess.Piece{Type: chess.Rook, Color: color}},
		{60, chess.Piece{Type: chess.King, Color: chess.Black}},
	}

	for _, item := range pieces {
		if err := board.SetPiece(item.square, item.piece); err != nil {
			t.Fatalf("Failed to set piece on board: %v", err)
		}
	}

	position, err := chess.NewPosition(board, color)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	position.CastlingRights = right

	return position
}

func whiteKingSidePosition(t *testing.T) chess.Position {
	t.Helper()
	return kingSidePosition(t, chess.White)
}

func blackKingSidePosition(t *testing.T) chess.Position {
	t.Helper()
	return kingSidePosition(t, chess.Black)
}

func TestCanCastleKingSideWhiteAllowed(t *testing.T) {
	position := whiteKingSidePosition(t)

	if !rules.CanCastleKingSide(position) {
		t.Error("Expected white king-side castling to be allowed")
	}
}

func TestCanCastleKingSideWhiteWithoutRight(t *testing.T) {
	for _, rights := range []chess.CastlingRights{0, chess.WhiteQueenSide} {
		position := whiteKingSidePosition(t)
		position.CastlingRights = rights

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected white king-side castling to be disallowed with rights %v", rights)
		}
	}
}

func TestCanCastleKingSideWhiteBlocked(t *testing.T) {
	for _, square := range []chess.Square{5, 6} {
		position := whiteKingSidePosition(t)
		blocker := chess.Piece{Type: chess.Knight, Color: chess.White}

		if err := position.Board.SetPiece(square, blocker); err != nil {
			t.Fatalf("Failed to set piece on board: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected white king-side castling to be disallowed when square %v is blocked", square)
		}
	}
}

func TestCanCastleKingSideWhiteAttackedSquares(t *testing.T) {
	cases := []struct {
		name string
		from chess.Square
	}{
		{"king in check", 52},
		{"crossing check", 53},
		{"destination check", 54},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := whiteKingSidePosition(t)
			attacker := chess.Piece{Type: chess.Queen, Color: chess.Black}

			if err := position.Board.SetPiece(tc.from, attacker); err != nil {
				t.Fatalf("Failed to set piece on board: %v", err)
			}

			if rules.CanCastleKingSide(position) {
				t.Errorf("Expected the attacker at %v to prevent castling", tc.from)
			}
		})
	}
}

func TestCanCastleKingSideWhiteMissingPiece(t *testing.T) {
	for _, square := range []chess.Square{4, 7} {
		position := whiteKingSidePosition(t)

		if _, err := position.Board.RemovePiece(square); err != nil {
			t.Fatalf("Failed to remove piece from board: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected missing piece at %v to prevent castling", square)
		}
	}
}

func TestCanCastleKingSideBlackAllowed(t *testing.T) {
	position := blackKingSidePosition(t)

	if !rules.CanCastleKingSide(position) {
		t.Error("Expected black king-side castling to be allowed")
	}
}

func TestCanCastleKingSideBlackWithoutRight(t *testing.T) {
	for _, rights := range []chess.CastlingRights{0, chess.BlackQueenSide, chess.WhiteKingSide} {
		position := blackKingSidePosition(t)
		position.CastlingRights = rights

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected black king-side castling to be disallowed with rights %v", rights)
		}
	}
}

func TestCanCastleKingSideBlackBlocked(t *testing.T) {
	for _, square := range []chess.Square{61, 62} {
		position := blackKingSidePosition(t)
		blocker := chess.Piece{Type: chess.Knight, Color: chess.Black}

		if err := position.Board.SetPiece(square, blocker); err != nil {
			t.Fatalf("Failed to set blocker: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected the piece at %v to block black castling", square)
		}
	}
}

func TestCanCastleKingSideBlackAttackedSquares(t *testing.T) {
	cases := []struct {
		name string
		from chess.Square
	}{
		{"king in check", 12},     // e2 attacks e8
		{"crossing check", 13},    // f2 attacks f8
		{"destination check", 14}, // g2 attacks g8
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := blackKingSidePosition(t)
			attacker := chess.Piece{Type: chess.Queen, Color: chess.White}

			if err := position.Board.SetPiece(tc.from, attacker); err != nil {
				t.Fatalf("Failed to set attacker: %v", err)
			}

			if rules.CanCastleKingSide(position) {
				t.Errorf("Expected the attacker at %v to prevent black castling", tc.from)
			}
		})
	}
}

func TestCanCastleKingSideBlackMissingPiece(t *testing.T) {
	for _, square := range []chess.Square{60, 63} {
		position := blackKingSidePosition(t)

		if _, err := position.Board.RemovePiece(square); err != nil {
			t.Fatalf("Failed to remove piece: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected missing piece at %v to prevent black castling", square)
		}
	}
}

func TestCanCastleKingSideInvalidColor(t *testing.T) {
	for _, color := range []chess.Color{chess.NoColor, chess.Color(255)} {
		position := whiteKingSidePosition(t)
		position.SideToMove = color

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected invalid color %v to prevent castling", color)
		}
	}
}

func TestKingSideCastlingMove(t *testing.T) {
	cases := []struct {
		name  string
		color chess.Color
		want  chess.Move
	}{
		{"white king-side", chess.White, chess.Move{From: 4, To: 6}},
		{"black king-side", chess.Black, chess.Move{From: 60, To: 62}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := kingSidePosition(t, tc.color)

			move, ok := rules.KingSideCastlingMove(position)
			if !ok {
				t.Fatal("Expected king-side castling to be available")
			}

			if move != tc.want {
				t.Errorf("Expected move %v, got %v", tc.want, move)
			}
		})
	}
}

func TestKingSideCastlingMoveUnavailable(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		position := kingSidePosition(t, color)
		position.CastlingRights = 0

		_, ok := rules.KingSideCastlingMove(position)

		if ok {
			t.Errorf("Expected king-side castling to be unavailable for color %v", color)
		}
	}
}
