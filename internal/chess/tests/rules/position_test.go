package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestPositionAfterMoveState(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide

	cases := []struct {
		name             string
		color, nextColor chess.Color
		kingSide         bool
		move             chess.Move
		wantRights       chess.CastlingRights
		enPassant        chess.Square
		wantFullmove     uint16
	}{
		{"white h-rook", chess.White, chess.Black, true, chess.Move{From: 7, To: 15}, chess.WhiteQueenSide | blackRights, 42, 20},
		{"black h-rook", chess.Black, chess.White, true, chess.Move{From: 63, To: 55}, chess.BlackQueenSide | whiteRights, 18, 21},
		{"white a-rook", chess.White, chess.Black, false, chess.Move{From: 0, To: 8}, chess.WhiteKingSide | blackRights, 42, 20},
		{"black a-rook", chess.Black, chess.White, false, chess.Move{From: 56, To: 48}, chess.BlackKingSide | whiteRights, 18, 21},
		{"white king", chess.White, chess.Black, true, chess.Move{From: 4, To: 12}, blackRights, 42, 20},
		{"black king", chess.Black, chess.White, true, chess.Move{From: 60, To: 52}, whiteRights, 18, 21},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			position.CastlingRights = whiteRights | blackRights
			position.EnPassant = tc.enPassant
			position.HalfmoveClock = 7
			position.FullmoveNumber = 20

			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if next.SideToMove != tc.nextColor {
				t.Errorf("Expected next color %v, got %v", tc.nextColor, next.SideToMove)
			}
			if next.FullmoveNumber != tc.wantFullmove {
				t.Errorf("Expected fullmove number %d, got %d", tc.wantFullmove, next.FullmoveNumber)
			}
			if next.HalfmoveClock != 8 {
				t.Errorf("Expected halfmove clock 8, got %d", next.HalfmoveClock)
			}
			if next.EnPassant != chess.NoSquare {
				t.Errorf("Expected en passant target to expire, got %v", next.EnPassant)
			}
			if next.CastlingRights != tc.wantRights {
				t.Errorf("Expected rights %v, got %v", tc.wantRights, next.CastlingRights)
			}
		})
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

func TestPositionAfterMoveCaptures(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide

	cases := []struct {
		name       string
		color      chess.Color
		enemy      chess.Color
		pieceType  chess.PieceType
		move       chess.Move
		wantRights chess.CastlingRights
	}{
		{"white captures h8 rook", chess.White, chess.Black, chess.Bishop, chess.Move{From: 54, To: 63}, whiteRights | chess.BlackQueenSide},
		{"white captures a8 rook", chess.White, chess.Black, chess.Bishop, chess.Move{From: 49, To: 56}, whiteRights | chess.BlackKingSide},
		{"black captures h1 rook", chess.Black, chess.White, chess.Bishop, chess.Move{From: 14, To: 7}, blackRights | chess.WhiteQueenSide},
		{"black captures a1 rook", chess.Black, chess.White, chess.Bishop, chess.Move{From: 9, To: 0}, blackRights | chess.WhiteKingSide},
		{"white captures rook outside starting square", chess.White, chess.Black, chess.Rook, chess.Move{From: 7, To: 15}, chess.WhiteQueenSide | blackRights},
		{"black captures rook outside starting square", chess.Black, chess.White, chess.Rook, chess.Move{From: 63, To: 55}, chess.BlackQueenSide | whiteRights},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := kingSidePosition(t, tc.color)
			mover := chess.Piece{Type: tc.pieceType, Color: tc.color}
			captured := chess.Piece{Type: chess.Rook, Color: tc.enemy}
			if err := position.Board.SetPiece(tc.move.From, mover); err != nil {
				t.Fatalf("Failed to set moving piece: %v", err)
			}
			if err := position.Board.SetPiece(tc.move.To, captured); err != nil {
				t.Fatalf("Failed to set captured piece: %v", err)
			}
			position.CastlingRights = whiteRights | blackRights
			position.HalfmoveClock = 7
			before := position

			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			assertPositionPiece(t, next, tc.move.From, chess.Piece{})
			assertPositionPiece(t, next, tc.move.To, mover)
			if next.HalfmoveClock != 0 {
				t.Errorf("Expected capture to reset halfmove clock, got %d", next.HalfmoveClock)
			}
			if next.CastlingRights != tc.wantRights {
				t.Errorf("Expected rights %v, got %v", tc.wantRights, next.CastlingRights)
			}
			if position != before {
				t.Error("Expected original position to remain unchanged")
			}
		})
	}
}

func TestPositionAfterMovePawns(t *testing.T) {
	cases := []struct {
		name          string
		color, enemy  chess.Color
		move          chess.Move
		capture       bool
		wantEnPassant chess.Square
	}{
		{"white single advance", chess.White, chess.Black, chess.Move{From: 12, To: 20}, false, chess.NoSquare},
		{"white double advance", chess.White, chess.Black, chess.Move{From: 12, To: 28}, false, 20},
		{"white capture", chess.White, chess.Black, chess.Move{From: 12, To: 21}, true, chess.NoSquare},
		{"black single advance", chess.Black, chess.White, chess.Move{From: 52, To: 44}, false, chess.NoSquare},
		{"black double advance", chess.Black, chess.White, chess.Move{From: 52, To: 36}, false, 44},
		{"black capture", chess.Black, chess.White, chess.Move{From: 52, To: 45}, true, chess.NoSquare},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := kingSidePosition(t, tc.color)
			pawn := chess.Piece{Type: chess.Pawn, Color: tc.color}
			if err := position.Board.SetPiece(tc.move.From, pawn); err != nil {
				t.Fatalf("Failed to set pawn: %v", err)
			}
			if tc.capture {
				if err := position.Board.SetPiece(tc.move.To, chess.Piece{Type: chess.Knight, Color: tc.enemy}); err != nil {
					t.Fatalf("Failed to set captured piece: %v", err)
				}
			}
			position.EnPassant = 18
			position.HalfmoveClock = 7
			before := position

			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			assertPositionPiece(t, next, tc.move.From, chess.Piece{})
			assertPositionPiece(t, next, tc.move.To, pawn)
			if next.HalfmoveClock != 0 {
				t.Errorf("Expected pawn move to reset halfmove clock, got %d", next.HalfmoveClock)
			}
			if next.EnPassant != tc.wantEnPassant {
				t.Errorf("Expected en passant square %v, got %v", tc.wantEnPassant, next.EnPassant)
			}
			if next.CastlingRights != before.CastlingRights {
				t.Errorf("Expected pawn move to preserve rights %v, got %v", before.CastlingRights, next.CastlingRights)
			}
			if position != before {
				t.Error("Expected original position to remain unchanged")
			}
		})
	}
}

func TestPositionAfterMoveCastling(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide

	cases := []struct {
		name             string
		color, nextColor chess.Color
		kingSide         bool
		kingMove         chess.Move
		rookMove         chess.Move
		wantRights       chess.CastlingRights
		wantFullmove     uint16
	}{
		{"white king-side", chess.White, chess.Black, true, chess.Move{From: 4, To: 6}, chess.Move{From: 7, To: 5}, blackRights, 20},
		{"black king-side", chess.Black, chess.White, true, chess.Move{From: 60, To: 62}, chess.Move{From: 63, To: 61}, whiteRights, 21},
		{"white queen-side", chess.White, chess.Black, false, chess.Move{From: 4, To: 2}, chess.Move{From: 0, To: 3}, blackRights, 20},
		{"black queen-side", chess.Black, chess.White, false, chess.Move{From: 60, To: 58}, chess.Move{From: 56, To: 59}, whiteRights, 21},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := castlingPosition(t, tc.color, tc.kingSide)
			position.CastlingRights = whiteRights | blackRights
			position.EnPassant = 18
			position.HalfmoveClock = 7
			position.FullmoveNumber = 20
			before := position

			next, err := rules.PositionAfterMove(position, tc.kingMove)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			assertPositionPiece(t, next, tc.kingMove.From, chess.Piece{})
			assertPositionPiece(t, next, tc.rookMove.From, chess.Piece{})
			assertPositionPiece(t, next, tc.kingMove.To, chess.Piece{Type: chess.King, Color: tc.color})
			assertPositionPiece(t, next, tc.rookMove.To, chess.Piece{Type: chess.Rook, Color: tc.color})
			if next.SideToMove != tc.nextColor {
				t.Errorf("Expected next color %v, got %v", tc.nextColor, next.SideToMove)
			}
			if next.CastlingRights != tc.wantRights {
				t.Errorf("Expected rights %v, got %v", tc.wantRights, next.CastlingRights)
			}
			if next.EnPassant != chess.NoSquare {
				t.Errorf("Expected en passant target to expire, got %v", next.EnPassant)
			}
			if next.HalfmoveClock != 8 {
				t.Errorf("Expected halfmove clock 8, got %d", next.HalfmoveClock)
			}
			if next.FullmoveNumber != tc.wantFullmove {
				t.Errorf("Expected fullmove number %d, got %d", tc.wantFullmove, next.FullmoveNumber)
			}
			if position != before {
				t.Error("Expected original position to remain unchanged")
			}
		})
	}
}

func TestPositionAfterMoveRejected(t *testing.T) {
	colors := []struct {
		name         string
		color        chess.Color
		kingFrom     chess.Square
		rookFrom     chess.Square
		rookTo       chess.Square
		enemyKing    chess.Square
		ownKingSide  chess.Square
		ownQueenSide chess.Square
	}{
		{"white", chess.White, 4, 7, 15, 60, 6, 2},
		{"black", chess.Black, 60, 63, 55, 4, 62, 58},
	}

	for _, side := range colors {
		t.Run(side.name, func(t *testing.T) {
			cases := []struct {
				name          string
				move          chess.Move
				turn          chess.Color
				removeKing    bool
				withoutRights bool
				wantErr       error
			}{
				{"invalid source", chess.Move{From: chess.NoSquare, To: side.rookTo}, side.color, false, false, chess.ErrInvalidSquare},
				{"invalid destination", chess.Move{From: side.rookFrom, To: chess.Square(255)}, side.color, false, false, chess.ErrInvalidSquare},
				{"empty source", chess.Move{From: 27, To: 35}, side.color, false, false, chess.ErrEmptySource},
				{"missing turn color", chess.Move{From: side.rookFrom, To: side.rookTo}, chess.NoColor, false, false, chess.ErrInvalidColor},
				{"invalid turn color", chess.Move{From: side.rookFrom, To: side.rookTo}, chess.Color(255), false, false, chess.ErrInvalidColor},
				{"opposing piece", chess.Move{From: side.enemyKing, To: 27}, side.color, false, false, chess.ErrIllegalMove},
				{"missing king", chess.Move{From: side.rookFrom, To: side.rookTo}, side.color, true, false, chess.ErrKingNotFound},
				{"diagonal rook move", chess.Move{From: side.rookFrom, To: side.rookTo - 1}, side.color, false, false, chess.ErrIllegalMove},
				{"same source and destination", chess.Move{From: side.rookFrom, To: side.rookFrom}, side.color, false, false, chess.ErrIllegalMove},
				{"own piece at destination", chess.Move{From: side.rookFrom, To: side.kingFrom}, side.color, false, false, chess.ErrIllegalMove},
				{"promotion on rook move", chess.Move{From: side.rookFrom, To: side.rookTo, Promotion: chess.Queen}, side.color, false, false, chess.ErrIllegalMove},
				{"king-side castling without rights", chess.Move{From: side.kingFrom, To: side.ownKingSide}, side.color, false, true, chess.ErrIllegalMove},
				{"queen-side castling without rights", chess.Move{From: side.kingFrom, To: side.ownQueenSide}, side.color, false, true, chess.ErrIllegalMove},
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					position := castlingPosition(t, side.color, tc.move.To != side.ownQueenSide)
					position.SideToMove = tc.turn
					if tc.removeKing {
						if _, err := position.Board.RemovePiece(side.kingFrom); err != nil {
							t.Fatalf("Failed to remove king: %v", err)
						}
					}
					if tc.withoutRights {
						position.CastlingRights = 0
					}
					before := position

					_, err := rules.PositionAfterMove(position, tc.move)
					if !errors.Is(err, tc.wantErr) {
						t.Errorf("Expected error %v, got %v", tc.wantErr, err)
					}
					if position != before {
						t.Error("Expected rejected move to preserve original position")
					}
				})
			}
		})
	}
}

func TestPositionAfterMoveBasicMove(t *testing.T) {
	cases := []struct {
		name  string
		color chess.Color
		move  chess.Move
	}{
		{"white", chess.White, chess.Move{From: 7, To: 15}},
		{"black", chess.Black, chess.Move{From: 63, To: 55}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := kingSidePosition(t, tc.color)
			before := position

			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			source, _ := next.Board.PieceAt(tc.move.From)
			if source.Type != chess.NoPiece {
				t.Error("Expected source square to be empty")
			}

			target, _ := next.Board.PieceAt(tc.move.To)
			want := chess.Piece{Type: chess.Rook, Color: tc.color}
			if target != want {
				t.Errorf("Expected %v at destination, got %v", want, target)
			}

			if position != before {
				t.Error("Expected original position to remain unchanged")
			}
		})
	}
}
