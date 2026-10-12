package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestEvaluateStatus(t *testing.T) {
	cases := []struct {
		name  string
		queen chess.Square
		want  chess.Status
	}{
		{"checkmate", chess.B7, chess.StatusCheckmate},
		{"stalemate", chess.B6, chess.StatusStalemate},
		{"ongoing", chess.A6, chess.StatusOngoing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, chess.A8, chess.Piece{Type: chess.King, Color: chess.Black})
			mustSetPiece(t, &board, chess.C6, chess.Piece{Type: chess.King, Color: chess.White})
			mustSetPiece(t, &board, tc.queen, chess.Piece{Type: chess.Queen, Color: chess.White})
			position := chess.Position{Board: board, SideToMove: chess.Black, EnPassant: chess.NoSquare}
			got, err := rules.EvaluateStatus(position)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Expected status %v, got %v", tc.want, got)
			}
		})
	}
}

func TestEvaluateStatusEnPassantOnlyMove(t *testing.T) {
	sides := []struct {
		name         string
		color, enemy chess.Color
		king, pawn   chess.Square
		enemyKing    chess.Square
		enemyBishop  chess.Square
		enemyPawns   []chess.Square
		enPassant    chess.Square
	}{
		{"white", chess.White, chess.Black, chess.H5, chess.E5, chess.F5, chess.G3,
			[]chess.Square{chess.H6, chess.G7, chess.E6, chess.D5}, chess.D6},
		{"black", chess.Black, chess.White, chess.H4, chess.E4, chess.F4, chess.G6,
			[]chess.Square{chess.H3, chess.G2, chess.E3, chess.D4}, chess.D3},
	}
	for _, side := range sides {
		t.Run(side.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, side.king, chess.Piece{Type: chess.King, Color: side.color})
			mustSetPiece(t, &board, side.pawn, chess.Piece{Type: chess.Pawn, Color: side.color})
			mustSetPiece(t, &board, side.enemyKing, chess.Piece{Type: chess.King, Color: side.enemy})
			mustSetPiece(t, &board, side.enemyBishop, chess.Piece{Type: chess.Bishop, Color: side.enemy})
			for _, square := range side.enemyPawns {
				mustSetPiece(t, &board, square, chess.Piece{Type: chess.Pawn, Color: side.enemy})
			}
			if check, err := rules.IsInCheck(board, side.color); err != nil || check {
				t.Fatalf("Expected setup with safe king, got check=%v, error=%v", check, err)
			}
			for _, from := range []chess.Square{side.king, side.pawn} {
				moves, err := rules.LegalMoves(board, from)
				if err != nil || len(moves) != 0 {
					t.Fatalf("Expected no basic legal moves from %v, got %v, %v", from, moves, err)
				}
			}
			cases := []struct {
				name      string
				enPassant chess.Square
				canMove   bool
				status    chess.Status
			}{
				{"capture available", side.enPassant, true, chess.StatusOngoing},
				{"without opportunity", chess.NoSquare, false, chess.StatusStalemate},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					position, err := chess.NewPosition(board, side.color)
					if err != nil {
						t.Fatalf("Failed to create position: %v", err)
					}
					position.EnPassant = tc.enPassant
					before := position
					moves, err := rules.LegalMovesFrom(position, side.pawn)
					if err != nil {
						t.Fatalf("Unexpected generation error: %v", err)
					}
					var destinations []chess.Square
					if tc.canMove {
						destinations = []chess.Square{side.enPassant}
					}
					assertMoves(t, moves, side.pawn, destinations)
					canMove, err := rules.HasLegalMove(position)
					if err != nil || canMove != tc.canMove {
						t.Errorf("Expected move availability %v, got %v, %v", tc.canMove, canMove, err)
					}
					status, err := rules.EvaluateStatus(position)
					if err != nil || status != tc.status {
						t.Errorf("Expected status %v, got %v, %v", tc.status, status, err)
					}
					if position != before {
						t.Error("Expected evaluation to preserve original position")
					}
				})
			}
		})
	}
}

func TestEvaluateStatusErrors(t *testing.T) {
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
			position := chess.Position{Board: board, SideToMove: tc.color, EnPassant: chess.NoSquare}
			_, err := rules.EvaluateStatus(position)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
