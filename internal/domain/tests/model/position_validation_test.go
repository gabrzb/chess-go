package model_test

import (
	"errors"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestPositionValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*chess.Position)
		want   error
	}{
		{"initial", func(*chess.Position) {}, nil},
		{"color", func(p *chess.Position) { p.SideToMove = chess.NoColor }, chess.ErrInvalidColor},
		{"clock", func(p *chess.Position) { p.FullmoveNumber = 0 }, chess.ErrInvalidPosition},
		{"flags", func(p *chess.Position) { p.CastlingRights = 128 }, chess.ErrInvalidPosition},
		{"missing king", func(p *chess.Position) { p.Board.RemovePiece(chess.E1) }, chess.ErrKingNotFound},
		{"duplicate king", func(p *chess.Position) { p.Board.SetPiece(chess.E3, chess.Piece{Type: chess.King, Color: chess.White}) }, chess.ErrInvalidPosition},
		{"piece type", func(p *chess.Position) { p.Board.SetPiece(chess.E3, chess.Piece{Type: 255, Color: chess.White}) }, chess.ErrInvalidPosition},
		{"piece color", func(p *chess.Position) { p.Board.SetPiece(chess.E3, chess.Piece{Type: chess.Pawn}) }, chess.ErrInvalidPosition},
		{"empty piece color", func(p *chess.Position) { p.Board.SetPiece(chess.E3, chess.Piece{Color: chess.White}) }, chess.ErrInvalidPosition},
		{"pawn rank", func(p *chess.Position) { p.Board.SetPiece(chess.A1, chess.Piece{Type: chess.Pawn, Color: chess.White}) }, chess.ErrInvalidPosition},
		{"castling rook", func(p *chess.Position) { p.Board.RemovePiece(chess.H1) }, chess.ErrInvalidPosition},
		{"castling king", func(p *chess.Position) {
			p.Board.RemovePiece(chess.E1)
			p.Board.SetPiece(chess.E3, chess.Piece{Type: chess.King, Color: chess.White})
		}, chess.ErrInvalidPosition},
		{"en passant absent pawn", func(p *chess.Position) { p.EnPassant = chess.E6 }, chess.ErrInvalidPosition},
		{"en passant bounds", func(p *chess.Position) { p.EnPassant = 255 }, chess.ErrInvalidPosition},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := chess.NewInitialPosition()
			tc.change(&p)
			before := p
			if err := p.Validate(); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
			if p != before {
				t.Fatal("validation changed the position")
			}
		})
	}
}

func TestPositionValidationEnPassant(t *testing.T) {
	for _, tc := range []struct {
		side             chess.Color
		from, to, target chess.Square
	}{
		{chess.Black, chess.E2, chess.E4, chess.E3},
		{chess.White, chess.D7, chess.D5, chess.D6},
	} {
		p := chess.NewInitialPosition()
		piece, _ := p.Board.RemovePiece(tc.from)
		p.Board.SetPiece(tc.to, piece)
		p.SideToMove, p.EnPassant = tc.side, tc.target
		if err := p.Validate(); err != nil {
			t.Fatal(err)
		}
		p.HalfmoveClock = 1
		if !errors.Is(p.Validate(), chess.ErrInvalidPosition) {
			t.Fatal("accepted en passant with a nonzero clock")
		}
	}
}
