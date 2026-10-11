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

			movingPiece, _ := position.Board.PieceAt(tc.move.From)
			before := position
			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			assertPositionPiece(t, next, tc.move.From, chess.Piece{})
			assertPositionPiece(t, next, tc.move.To, movingPiece)
			assertPositionState(t, next, positionState{tc.nextColor, tc.wantRights, chess.NoSquare, 8, tc.wantFullmove})
			if position != before {
				t.Error("Expected original position to remain unchanged")
			}

		})
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
			mustSetPiece(t, &position.Board, tc.move.From, mover)
			mustSetPiece(t, &position.Board, tc.move.To, captured)
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
			mustSetPiece(t, &position.Board, tc.move.From, pawn)
			if tc.capture {
				mustSetPiece(t, &position.Board, tc.move.To, chess.Piece{Type: chess.Knight, Color: tc.enemy})
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
	for _, castle := range castlingCases {
		t.Run(castle.name, func(t *testing.T) {
			position := castlingPosition(t, castle.color, castle.kingSide)
			position.CastlingRights = whiteRights | blackRights
			position.EnPassant = 18
			position.HalfmoveClock = 7
			position.FullmoveNumber = 20
			before := position
			next, err := rules.PositionAfterMove(position, castle.kingMove)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			assertPositionPiece(t, next, castle.kingMove.From, chess.Piece{})
			assertPositionPiece(t, next, castle.rookMove.From, chess.Piece{})
			assertPositionPiece(t, next, castle.kingMove.To, chess.Piece{Type: chess.King, Color: castle.color})
			assertPositionPiece(t, next, castle.rookMove.To, chess.Piece{Type: chess.Rook, Color: castle.color})
			wantRights, wantFullmove := blackRights, uint16(20)
			if castle.color == chess.Black {
				wantRights, wantFullmove = whiteRights, 21
			}
			assertPositionState(t, next, positionState{castle.enemy, wantRights, chess.NoSquare, 8, wantFullmove})
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
