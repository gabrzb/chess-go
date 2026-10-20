package engine_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
	"github.com/gabrzb/chess-go/internal/engine"
)

func positionWith(t *testing.T, side chess.Color, pieces map[chess.Square]chess.Piece) chess.Position {
	t.Helper()
	var board chess.Board
	for square, piece := range pieces {
		setPiece(t, &board, square, piece.Type, piece.Color)
	}
	position, err := chess.NewPosition(board, side)
	if err != nil {
		t.Fatal(err)
	}
	return position
}

func allMoves(t *testing.T, position chess.Position) []chess.Move {
	t.Helper()
	var moves []chess.Move
	for from := chess.A1; from <= chess.H8; from++ {
		piece, _ := position.Board.PieceAt(from)
		if piece.Type == chess.NoPiece || piece.Color != position.SideToMove {
			continue
		}
		fromMoves, err := rules.LegalMovesFrom(position, from)
		if err != nil {
			t.Fatal(err)
		}
		moves = append(moves, fromMoves...)
	}
	return moves
}

func player(t *testing.T, difficulty engine.Difficulty, seed int64) *engine.Player {
	t.Helper()
	p, err := engine.NewPlayer(difficulty, seed)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSearchFindsMate(t *testing.T) {
	for _, tc := range []struct {
		name                string
		color, enemy        chess.Color
		king, queen, target chess.Square
		move                chess.Move
	}{
		{"White", chess.White, chess.Black, chess.F6, chess.G6, chess.H8, chess.Move{From: chess.G6, To: chess.G7}},
		{"Black", chess.Black, chess.White, chess.F3, chess.G3, chess.H1, chess.Move{From: chess.G3, To: chess.G2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			position := positionWith(t, tc.color, map[chess.Square]chess.Piece{
				tc.king:   {Type: chess.King, Color: tc.color},
				tc.queen:  {Type: chess.Queen, Color: tc.color},
				tc.target: {Type: chess.King, Color: tc.enemy},
			})
			move, err := player(t, engine.Easy, 1).ChooseMove(t.Context(), position)
			if err != nil || move != tc.move {
				t.Fatalf("Expected mate %v, got %v: %v", tc.move, move, err)
			}
			next, err := rules.PositionAfterMove(position, move)
			if err != nil {
				t.Fatal(err)
			}
			if status, err := rules.EvaluateStatus(next); err != nil || status != chess.StatusCheckmate {
				t.Fatalf("Not checkmate: %v, %v", status, err)
			}
		})
	}
}

// The oracle uses full minimax with White maximizing and Black minimizing,
// independently of the production negamax bounds and move ordering.
func minimax(t *testing.T, position chess.Position, depth, ply int) int {
	t.Helper()
	moves := allMoves(t, position)
	if len(moves) == 0 {
		check, err := rules.IsInCheck(position.Board, position.SideToMove)
		if err != nil {
			t.Fatal(err)
		}
		if !check {
			return 0
		}
		const mate = 1_000_000
		score := mate - ply
		if position.SideToMove == chess.White {
			score = -score
		}
		return score
	}
	if depth == 0 {
		return engine.Evaluate(position.Board)
	}
	best := 0
	for index, move := range moves {
		next, err := rules.PositionAfterMove(position, move)
		if err != nil {
			t.Fatal(err)
		}
		score := minimax(t, next, depth-1, ply+1)
		if index == 0 || (position.SideToMove == chess.White && score > best) || (position.SideToMove == chess.Black && score < best) {
			best = score
		}
	}
	return best
}

func TestAlphaBetaMatchesFullMinimax(t *testing.T) {
	base := positionWith(t, chess.White, map[chess.Square]chess.Piece{
		chess.A1: {Type: chess.King, Color: chess.White},
		chess.H8: {Type: chess.King, Color: chess.Black},
		chess.D4: {Type: chess.Rook, Color: chess.White},
		chess.E6: {Type: chess.Rook, Color: chess.Black},
	})
	const difficulty = engine.Hard
	for _, side := range []chess.Color{chess.White, chess.Black} {
		name := map[chess.Color]string{chess.White: "White", chess.Black: "Black"}[side]
		t.Run(name, func(t *testing.T) {
			position := base
			position.SideToMove = side
			want := minimax(t, position, difficulty.Depth(), 0)
			for _, seed := range []int64{1, 42} {
				move, err := player(t, difficulty, seed).ChooseMove(context.Background(), position)
				if err != nil {
					t.Fatal(err)
				}
				next, err := rules.PositionAfterMove(position, move)
				if err != nil {
					t.Fatalf("Expected legal AI move: %v", err)
				}
				if got := minimax(t, next, difficulty.Depth()-1, 1); got != want {
					t.Errorf("Seed %d: expected optimal score %d, got %d from %v", seed, want, got, move)
				}
			}
		})
	}
}

func TestSearchLegalityAndReproducibleTies(t *testing.T) {
	position := chess.NewInitialPosition()
	before := position
	moves := allMoves(t, position)
	chosen := make(map[chess.Move]bool)
	for seed := int64(0); seed < 8; seed++ {
		p := player(t, engine.Easy, seed)
		first, err := p.ChooseMove(context.Background(), position)
		if err != nil {
			t.Fatal(err)
		}
		second, err := p.ChooseMove(context.Background(), position)
		if err != nil || first != second || !slices.Contains(moves, first) {
			t.Fatalf("Expected reproducible legal move for seed %d: %v, %v, %v", seed, first, second, err)
		}
		chosen[first] = true
	}
	if len(chosen) < 2 {
		t.Error("Expected seeds to vary between equally evaluated opening moves")
	}
	if position != before {
		t.Error("Expected search to preserve castling, clocks, en passant and board")
	}
}

func TestSearchCancellation(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	during, cancelDuring := context.WithCancel(context.Background())
	defer cancelDuring()
	for _, tc := range []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{"before search", canceled, context.Canceled},
		{"during search", during, context.Canceled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "during search" {
				timer := time.AfterFunc(time.Millisecond, cancelDuring)
				defer timer.Stop()
			}
			move, err := player(t, engine.Hard, 1).ChooseMove(tc.ctx, chess.NewInitialPosition())
			if !errors.Is(err, tc.err) || move != (chess.Move{}) {
				t.Errorf("Expected cancellation without a move, got %v: %v", move, err)
			}
		})
	}
}

func TestSearchWithoutLegalMoves(t *testing.T) {
	position := positionWith(t, chess.Black, map[chess.Square]chess.Piece{
		chess.F7: {Type: chess.King, Color: chess.White},
		chess.G6: {Type: chess.Queen, Color: chess.White},
		chess.H8: {Type: chess.King, Color: chess.Black},
	})
	move, err := player(t, engine.Medium, 1).ChooseMove(t.Context(), position)
	if !errors.Is(err, engine.ErrNoLegalMoves) || move != (chess.Move{}) {
		t.Fatalf("Expected no move after stalemate: %v, %v", move, err)
	}
}
