package ai_test

import (
	"context"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

const (
	oracleMateScore = 1_000_000
	oracleInfinity  = oracleMateScore + 1
)

func tacticalPosition(t *testing.T, pieces map[chess.Square]chess.Piece) chess.Position {
	t.Helper()
	var board chess.Board
	for square, piece := range pieces {
		if err := board.SetPiece(square, piece); err != nil {
			t.Fatal(err)
		}
	}
	p, err := chess.NewPosition(board, chess.White)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestQuiescenceSeesRecapture(t *testing.T) {
	p := tacticalPosition(t, map[chess.Square]chess.Piece{
		chess.A1: {Type: chess.King, Color: chess.White}, chess.H8: {Type: chess.King, Color: chess.Black},
		chess.D1: {Type: chess.Queen, Color: chess.White}, chess.D5: {Type: chess.Pawn, Color: chess.Black},
		chess.D8: {Type: chess.Rook, Color: chess.Black},
	})
	player, _ := ai.NewPlayer(ai.Easy, 1)
	move, err := player.ChooseMove(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if move == (chess.Move{From: chess.D1, To: chess.D5}) {
		t.Fatal("queen sacrificed for a defended pawn at the search horizon")
	}
}

// This oracle uses White's perspective and exhaustive minimax without pruning.
func referenceQuiet(t *testing.T, p chess.Position, remaining, ply int) int {
	t.Helper()
	status, err := rules.EvaluateStatus(p)
	if err != nil {
		t.Fatal(err)
	}
	if status == chess.StatusCheckmate {
		if p.SideToMove == chess.White {
			return -oracleMateScore + ply
		}
		return oracleMateScore - ply
	}
	if status != chess.StatusOngoing {
		return 0
	}
	best := ai.Evaluate(p.Board)
	if remaining == 0 {
		return best
	}
	check, err := rules.IsInCheck(p.Board, p.SideToMove)
	if err != nil {
		t.Fatal(err)
	}
	if check {
		best = oracleInfinity
		if p.SideToMove == chess.White {
			best = -oracleInfinity
		}
	}
	moves := allMoves(t, p)
	for _, move := range moves {
		victim, _ := p.Board.PieceAt(move.To)
		piece, _ := p.Board.PieceAt(move.From)
		quiet := victim.Type == chess.NoPiece && move.Promotion == chess.NoPiece && !(piece.Type == chess.Pawn && move.From.Column() != move.To.Column())
		if !check && quiet {
			continue
		}
		next, err := rules.PositionAfterMove(p, move)
		if err != nil {
			t.Fatal(err)
		}
		score := referenceQuiet(t, next, remaining-1, ply+1)
		if p.SideToMove == chess.White {
			best = max(best, score)
		} else {
			best = min(best, score)
		}
	}
	return best
}

func TestEasySearchMatchesExhaustiveTacticalSearch(t *testing.T) {
	for _, name := range []string{"en passant", "promotion", "check evasions"} {
		t.Run(name, func(t *testing.T) {
			pieces := map[chess.Square]chess.Piece{chess.E1: {Type: chess.King, Color: chess.White}, chess.H8: {Type: chess.King, Color: chess.Black}}
			switch name {
			case "en passant":
				pieces[chess.E5] = chess.Piece{Type: chess.Pawn, Color: chess.White}
				pieces[chess.D5] = chess.Piece{Type: chess.Pawn, Color: chess.Black}
			case "promotion":
				pieces[chess.A7] = chess.Piece{Type: chess.Pawn, Color: chess.White}
			case "check evasions":
				pieces[chess.E8] = chess.Piece{Type: chess.Rook, Color: chess.Black}
				pieces[chess.A1] = chess.Piece{Type: chess.Queen, Color: chess.White}
			}
			p := tacticalPosition(t, pieces)
			if name == "en passant" {
				p.EnPassant = chess.D6
			}
			for _, side := range []chess.Color{chess.White, chess.Black} {
				position := p
				if side == chess.Black {
					var board chess.Board
					for square := chess.A1; square <= chess.H8; square++ {
						piece, _ := p.Board.PieceAt(square)
						if piece.Type == chess.NoPiece {
							continue
						}
						if piece.Color == chess.White {
							piece.Color = chess.Black
						} else {
							piece.Color = chess.White
						}
						target, _ := chess.A8.Offset(square.Column(), -square.Row())
						_ = board.SetPiece(target, piece)
					}
					position, _ = chess.NewPosition(board, chess.Black)
					if p.EnPassant.Valid() {
						position.EnPassant, _ = chess.A8.Offset(p.EnPassant.Column(), -p.EnPassant.Row())
					}
				}
				moves := allMoves(t, position)
				want := oracleInfinity
				if side == chess.White {
					want = -oracleInfinity
				}
				for _, move := range moves {
					next, err := rules.PositionAfterMove(position, move)
					if err != nil {
						t.Fatal(err)
					}
					score := referenceQuiet(t, next, 8, 1)
					if side == chess.White {
						want = max(want, score)
					} else {
						want = min(want, score)
					}
				}
				chooser, err := ai.NewPlayer(ai.Easy, 1)
				if err != nil {
					t.Fatal(err)
				}
				chosen, err := chooser.ChooseMove(t.Context(), position)
				if err != nil {
					t.Fatal(err)
				}
				next, err := rules.PositionAfterMove(position, chosen)
				if err != nil {
					t.Fatal(err)
				}
				if got := referenceQuiet(t, next, 8, 1); got != want {
					t.Fatalf("side %v: got %d, want %d", side, got, want)
				}

			}
		})
	}
}
