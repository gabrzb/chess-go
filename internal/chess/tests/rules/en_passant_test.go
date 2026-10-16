package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestEnPassantExecutionAndExpiration(t *testing.T) {
	for _, tc := range []struct {
		fen              string
		advance, capture chess.Move
		captured         chess.Square
		decline, reply   chess.Move
	}{
		{"4k3/3p4/8/4P3/8/8/8/4K3 b - - 0 1", chess.Move{From: chess.D7, To: chess.D5}, chess.Move{From: chess.E5, To: chess.D6}, chess.D5, chess.Move{From: chess.E1, To: chess.F1}, chess.Move{From: chess.E8, To: chess.F8}},
		{"4k3/8/8/8/4p3/8/3P4/4K3 w - - 0 1", chess.Move{From: chess.D2, To: chess.D4}, chess.Move{From: chess.E4, To: chess.D3}, chess.D4, chess.Move{From: chess.E8, To: chess.F8}, chess.Move{From: chess.E1, To: chess.F1}},
	} {
		p := after(t, perftPositionFromFEN(t, tc.fen), tc.advance)
		if p.EnPassant != tc.capture.To {
			t.Fatal("Double advance did not record en passant")
		}
		next := after(t, p, tc.capture)
		assertPiece(t, next, tc.capture.From, chess.Piece{})
		assertPiece(t, next, tc.captured, chess.Piece{})
		assertPiece(t, next, tc.capture.To, chess.Piece{Type: chess.Pawn, Color: p.SideToMove})
		if next.EnPassant != chess.NoSquare || next.HalfmoveClock != 0 {
			t.Error("En passant did not clear the opportunity and clock")
		}
		expired := after(t, after(t, p, tc.decline), tc.reply)
		if _, err := rules.PositionAfterMove(expired, tc.capture); !errors.Is(err, chess.ErrIllegalMove) {
			t.Error("Expired en passant was accepted")
		}
	}
}

func TestEnPassantDiscoveredCheck(t *testing.T) {
	for _, tc := range []struct {
		fen  string
		move chess.Move
	}{
		{"4k3/8/8/K2pP2r/8/8/8/8 w - d6 0 1", chess.Move{From: chess.E5, To: chess.D6}},
		{"8/8/8/8/k2Pp2R/8/8/4K3 b - d3 0 1", chess.Move{From: chess.E4, To: chess.D3}},
	} {
		p := perftPositionFromFEN(t, tc.fen)
		if _, err := rules.PositionAfterMove(p, tc.move); !errors.Is(err, chess.ErrIllegalMove) {
			t.Errorf("En passant exposed the king: %v", err)
		}
	}
}

func TestEnPassantPreventsStalemate(t *testing.T) {
	for _, fen := range []string{
		"8/6p1/4p2p/3pPk1K/8/6b1/8/8 w - d6 0 1",
		"8/8/6B1/8/3PpK1k/4P2P/6P1/8 b - d3 0 1",
	} {
		p := perftPositionFromFEN(t, fen)
		if status, err := rules.EvaluateStatus(p); err != nil || status != chess.StatusOngoing {
			t.Fatalf("En passant must keep the game ongoing: %v, %v", status, err)
		}
		p.EnPassant = chess.NoSquare
		if status, err := rules.EvaluateStatus(p); err != nil || status != chess.StatusStalemate {
			t.Fatalf("Expected stalemate without en passant: %v, %v", status, err)
		}
	}
}
