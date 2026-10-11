package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func castlingPosition(t *testing.T, color chess.Color, kingSide bool) chess.Position {
	t.Helper()

	var board chess.Board
	rookSquare := chess.Square(0)
	right := chess.WhiteQueenSide
	if kingSide {
		rookSquare = 7
		right = chess.WhiteKingSide
	}

	if color == chess.Black {
		rookSquare += 56
		right = chess.BlackQueenSide
		if kingSide {
			right = chess.BlackKingSide
		}
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

func kingSidePosition(t *testing.T, color chess.Color) chess.Position {
	t.Helper()
	return castlingPosition(t, color, true)
}

func queenSidePosition(t *testing.T, color chess.Color) chess.Position {
	t.Helper()
	return castlingPosition(t, color, false)
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

func TestApplyCastlingPieces(t *testing.T) {
	cases := []struct {
		name             string
		color            chess.Color
		kingSide         bool
		apply            func(chess.Position) (chess.Position, error)
		kingFrom, kingTo chess.Square
		rookFrom, rookTo chess.Square
	}{
		{"white king-side", chess.White, true, rules.ApplyKingSideCastling, 4, 6, 7, 5},
		{"black king-side", chess.Black, true, rules.ApplyKingSideCastling, 60, 62, 63, 61},
		{"white queen-side", chess.White, false, rules.ApplyQueenSideCastling, 4, 2, 0, 3},
		{"black queen-side", chess.Black, false, rules.ApplyQueenSideCastling, 60, 58, 56, 59},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			before := position

			next, err := tc.apply(position)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			expected := []struct {
				square chess.Square
				piece  chess.Piece
			}{
				{tc.kingFrom, chess.Piece{}},
				{tc.rookFrom, chess.Piece{}},
				{tc.kingTo, chess.Piece{Type: chess.King, Color: tc.color}},
				{tc.rookTo, chess.Piece{Type: chess.Rook, Color: tc.color}},
			}

			for _, item := range expected {
				got, _ := next.Board.PieceAt(item.square)

				if got != item.piece {
					t.Errorf("At %v: expected %v, got %v", item.square, item.piece, got)
				}
			}

			if position != before {
				t.Error("Expected the original position to remain unchanged")
			}
		})
	}
}

func TestApplyCastlingState(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide

	cases := []struct {
		name             string
		color, nextColor chess.Color
		kingSide         bool
		apply            func(chess.Position) (chess.Position, error)
		wantRights       chess.CastlingRights
		enPassant        chess.Square
		wantFullmove     uint16
	}{
		{"white king-side", chess.White, chess.Black, true, rules.ApplyKingSideCastling, blackRights, 42, 20},
		{"black king-side", chess.Black, chess.White, true, rules.ApplyKingSideCastling, whiteRights, 18, 21},
		{"white queen-side", chess.White, chess.Black, false, rules.ApplyQueenSideCastling, blackRights, 42, 20},
		{"black queen-side", chess.Black, chess.White, false, rules.ApplyQueenSideCastling, whiteRights, 18, 21},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			position.CastlingRights = whiteRights | blackRights
			position.EnPassant = tc.enPassant
			position.HalfmoveClock = 7
			position.FullmoveNumber = 20

			next, err := tc.apply(position)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if next.SideToMove != tc.nextColor {
				t.Errorf("Expected next color %v, got %v", tc.nextColor, next.SideToMove)
			}

			if next.CastlingRights != tc.wantRights {
				t.Errorf("Expected rights %v, got %v", tc.wantRights, next.CastlingRights)
			}

			if next.EnPassant != chess.NoSquare {
				t.Error("Expected the en passant target to expire")
			}

			if next.HalfmoveClock != 8 {
				t.Errorf("Expected halfmove clock 8, got %d", next.HalfmoveClock)
			}

			if next.FullmoveNumber != tc.wantFullmove {
				t.Errorf("Expected fullmove number %d, got %d", tc.wantFullmove, next.FullmoveNumber)
			}
		})
	}
}

func TestApplyCastlingUnavailable(t *testing.T) {
	cases := []struct {
		name     string
		color    chess.Color
		kingSide bool
		apply    func(chess.Position) (chess.Position, error)
	}{
		{"white king-side", chess.White, true, rules.ApplyKingSideCastling},
		{"black king-side", chess.Black, true, rules.ApplyKingSideCastling},
		{"white queen-side", chess.White, false, rules.ApplyQueenSideCastling},
		{"black queen-side", chess.Black, false, rules.ApplyQueenSideCastling},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			position.CastlingRights = 0
			before := position

			_, err := tc.apply(position)

			if !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Expected ErrIllegalMove, got %v", err)
			}

			if position != before {
				t.Error("Expected the original position to remain unchanged after failed castling")
			}
		})
	}
}

func TestCanCastleQueenSideAAllowed(t *testing.T) {
	cases := []struct {
		name  string
		color chess.Color
	}{
		{"white", chess.White},
		{"black", chess.Black},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := queenSidePosition(t, tc.color)

			if !rules.CanCastleQueenSide(position) {
				t.Errorf("Expected %v queen-side castling to be allowed", tc.color)
			}
		})
	}
}

func TestCanCastleQueenSideWithoutRight(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		position := queenSidePosition(t, color)
		position.CastlingRights = 0

		if rules.CanCastleQueenSide(position) {
			t.Errorf("Expected %v queen-side castling to be disallowed without rights", color)
		}
	}
}

func TestCanCastleQueenSideBlocked(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		base := chess.Square(0)

		if color == chess.Black {
			base = 56
		}

		for _, offset := range []chess.Square{1, 2, 3} {
			position := queenSidePosition(t, color)
			square := base + offset
			blocker := chess.Piece{Type: chess.Knight, Color: color}

			if err := position.Board.SetPiece(square, blocker); err != nil {
				t.Fatalf("Failed to set blocker: %v", err)
			}

			if rules.CanCastleQueenSide(position) {
				t.Errorf("Expected %v queen-side castling to be disallowed when square %v is blocked", color, square)
			}
		}
	}
}

func TestCanCastleQueenSideAttackedSquares(t *testing.T) {
	cases := []struct {
		name       string
		color      chess.Color
		enemy      chess.Color
		attackkers []chess.Square
	}{
		{"white king in check", chess.White, chess.Black, []chess.Square{52, 51, 50}},
		{"black king in check", chess.Black, chess.White, []chess.Square{12, 11, 10}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, from := range tc.attackkers {
				position := queenSidePosition(t, tc.color)
				attacker := chess.Piece{Type: chess.Queen, Color: tc.enemy}

				if err := position.Board.SetPiece(from, attacker); err != nil {
					t.Fatalf("Failed to set attacker: %v", err)
				}

				if rules.CanCastleQueenSide(position) {
					t.Errorf("Expected %v queen-side castling to be disallowed when attacked from %v", tc.color, from)
				}
			}
		})
	}
}

func TestCanCastleQueenSideWithAttackedBSquare(t *testing.T) {
	cases := []struct {
		name  string
		color chess.Color
		enemy chess.Color
		from  chess.Square
	}{
		{"white", chess.White, chess.Black, 8},
		{"black", chess.Black, chess.White, 48},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := queenSidePosition(t, tc.color)
			attacker := chess.Piece{Type: chess.Bishop, Color: tc.enemy}

			if err := position.Board.SetPiece(tc.from, attacker); err != nil {
				t.Fatalf("Failed to set attacker: %v", err)
			}

			if !rules.CanCastleQueenSide(position) {
				t.Error("Expected attack on b-file square to allow castling")
			}
		})
	}
}

func TestCanCastleQueenSideMissingPiece(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		base := chess.Square(0)
		if color == chess.Black {
			base = 56
		}

		for _, offset := range []chess.Square{4, 0} {
			position := queenSidePosition(t, color)
			square := base + offset

			if _, err := position.Board.RemovePiece(square); err != nil {
				t.Fatalf("Failed to remove piece: %v", err)
			}

			if rules.CanCastleQueenSide(position) {
				t.Errorf("Expected missing piece at %v to prevent castling", square)
			}
		}
	}
}

func TestCastlingMove(t *testing.T) {
	cases := []struct {
		name     string
		color    chess.Color
		kingSide bool
		generate func(chess.Position) (chess.Move, bool)
		want     chess.Move
	}{
		{"white king-side", chess.White, true, rules.KingSideCastlingMove, chess.Move{From: 4, To: 6}},
		{"black king-side", chess.Black, true, rules.KingSideCastlingMove, chess.Move{From: 60, To: 62}},
		{"white queen-side", chess.White, false, rules.QueenSideCastlingMove, chess.Move{From: 4, To: 2}},
		{"black queen-side", chess.Black, false, rules.QueenSideCastlingMove, chess.Move{From: 60, To: 58}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)

			move, ok := tc.generate(position)
			if !ok {
				t.Fatal("Expected castling to be available")
			}

			if move != tc.want {
				t.Errorf("Expected move %v, got %v", tc.want, move)
			}
		})
	}
}

func TestCastlingMoveUnavailable(t *testing.T) {
	cases := []struct {
		name     string
		color    chess.Color
		kingSide bool
		generate func(chess.Position) (chess.Move, bool)
	}{
		{"white king-side", chess.White, true, rules.KingSideCastlingMove},
		{"black king-side", chess.Black, true, rules.KingSideCastlingMove},
		{"white queen-side", chess.White, false, rules.QueenSideCastlingMove},
		{"black queen-side", chess.Black, false, rules.QueenSideCastlingMove},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			position.CastlingRights = 0

			_, ok := tc.generate(position)
			if ok {
				t.Error("Expected castling to be unavailable without rights")
			}
		})
	}
}
