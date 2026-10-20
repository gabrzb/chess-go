package rules_test

import (
	"github.com/gabrzb/chess-go/internal/domain/rules"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestInsufficientMaterial(t *testing.T) {
	for _, tc := range []struct {
		name   string
		pieces map[chess.Square]chess.Piece
		want   bool
	}{
		{"kings", nil, true},
		{"one knight", map[chess.Square]chess.Piece{chess.B1: {Type: chess.Knight, Color: chess.White}}, true},
		{"same color bishops", map[chess.Square]chess.Piece{chess.C1: {Type: chess.Bishop, Color: chess.White}, chess.F8: {Type: chess.Bishop, Color: chess.Black}}, true},
		{"opposite color bishops", map[chess.Square]chess.Piece{chess.C1: {Type: chess.Bishop, Color: chess.White}, chess.C8: {Type: chess.Bishop, Color: chess.Black}}, false},
		{"two knights", map[chess.Square]chess.Piece{chess.B1: {Type: chess.Knight, Color: chess.White}, chess.G1: {Type: chess.Knight, Color: chess.White}}, false},
		{"pawn", map[chess.Square]chess.Piece{chess.A2: {Type: chess.Pawn, Color: chess.White}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			_ = board.SetPiece(chess.E1, chess.Piece{Type: chess.King, Color: chess.White})
			_ = board.SetPiece(chess.E8, chess.Piece{Type: chess.King, Color: chess.Black})
			for square, piece := range tc.pieces {
				_ = board.SetPiece(square, piece)
			}
			if got := rules.InsufficientMaterial(board); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRepetitionKey(t *testing.T) {
	p := chess.NewInitialPosition()
	p.EnPassant = chess.E3
	if rules.RepetitionKey(p).EnPassant != chess.NoSquare {
		t.Fatal("irrelevant en passant changes repetition")
	}
	p.EnPassant = chess.NoSquare
	q := p
	q.HalfmoveClock, q.FullmoveNumber = 20, 30
	if rules.RepetitionKey(p) != rules.RepetitionKey(q) {
		t.Fatal("clocks change repetition")
	}
	q.CastlingRights = 0
	if rules.RepetitionKey(p) == rules.RepetitionKey(q) {
		t.Fatal("castling rights ignored")
	}

	var board chess.Board
	for square, piece := range map[chess.Square]chess.Piece{
		chess.E1: {Type: chess.King, Color: chess.White}, chess.A8: {Type: chess.King, Color: chess.Black},
		chess.E5: {Type: chess.Pawn, Color: chess.White}, chess.D5: {Type: chess.Pawn, Color: chess.Black},
	} {
		_ = board.SetPiece(square, piece)
	}
	p, _ = chess.NewPosition(board, chess.White)
	p.EnPassant = chess.D6
	if rules.RepetitionKey(p).EnPassant != chess.D6 {
		t.Fatal("legal en passant ignored")
	}
	_ = p.Board.SetPiece(chess.E8, chess.Piece{Type: chess.Rook, Color: chess.Black})
	if rules.RepetitionKey(p).EnPassant != chess.NoSquare {
		t.Fatal("pinned en passant changes repetition")
	}
}
