package ai_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/rules"
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

func player(t *testing.T, difficulty ai.Difficulty, seed int64) *ai.Player {
	t.Helper()
	p, err := ai.NewPlayer(difficulty, seed)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSearchTacticsForBothColors(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		for _, mate := range []bool{false, true} {
			name := "White"
			if color == chess.Black {
				name = "Black"
			}
			if mate {
				name += "/mate in one"
			} else {
				name += "/free queen"
			}
			t.Run(name, func(t *testing.T) {
				pieces := map[chess.Square]chess.Piece{
					chess.A1: {Type: chess.King, Color: chess.White},
					chess.H8: {Type: chess.King, Color: chess.Black},
					chess.D4: {Type: chess.Rook, Color: chess.White},
					chess.D5: {Type: chess.Queen, Color: chess.Black},
				}
				want := chess.Move{From: chess.D4, To: chess.D5}
				if mate {
					pieces = map[chess.Square]chess.Piece{
						chess.F6: {Type: chess.King, Color: chess.White},
						chess.G6: {Type: chess.Queen, Color: chess.White},
						chess.H8: {Type: chess.King, Color: chess.Black},
					}
					want = chess.Move{From: chess.G6, To: chess.G7}
				}
				if color == chess.Black {
					mirrored := make(map[chess.Square]chess.Piece)
					for square, piece := range pieces {
						square, _ = chess.A8.Offset(square.Column(), -square.Row())
						if piece.Color == chess.White {
							piece.Color = chess.Black
						} else {
							piece.Color = chess.White
						}
						mirrored[square] = piece
					}
					pieces = mirrored
					want.From, _ = chess.A8.Offset(want.From.Column(), -want.From.Row())
					want.To, _ = chess.A8.Offset(want.To.Column(), -want.To.Row())
				}
				position := positionWith(t, color, pieces)
				before := position
				move, err := player(t, ai.Easy, 7).ChooseMove(context.Background(), position)
				if err != nil || move != want {
					t.Fatalf("Expected %v, got %v: %v", want, move, err)
				}
				if position != before {
					t.Error("Expected search to preserve its input position")
				}
				if mate {
					next, err := rules.PositionAfterMove(position, move)
					if err != nil {
						t.Fatal(err)
					}
					status, err := rules.EvaluateStatus(next)
					if err != nil || status != chess.StatusCheckmate {
						t.Fatalf("Expected checkmate, got %v: %v", status, err)
					}
				}
			})
		}
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
		return ai.Evaluate(position.Board)
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
	for _, side := range []chess.Color{chess.White, chess.Black} {
		for _, difficulty := range []ai.Difficulty{ai.Easy, ai.Medium, ai.Hard} {
			name := map[chess.Color]string{chess.White: "White", chess.Black: "Black"}[side]
			t.Run(name+"/"+string(difficulty), func(t *testing.T) {
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
}

func TestSearchLegalityAndReproducibleTies(t *testing.T) {
	position := chess.NewInitialPosition()
	before := position
	moves := allMoves(t, position)
	chosen := make(map[chess.Move]bool)
	for seed := int64(0); seed < 8; seed++ {
		p := player(t, ai.Easy, seed)
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

func TestSearchInitialDifficultyLevels(t *testing.T) {
	position := chess.NewInitialPosition()
	for _, difficulty := range []ai.Difficulty{ai.Easy, ai.Medium, ai.Hard} {
		t.Run(string(difficulty), func(t *testing.T) {
			move, err := player(t, difficulty, 1).ChooseMove(t.Context(), position)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := rules.PositionAfterMove(position, move); err != nil {
				t.Fatalf("Expected legal opening move at difficulty %s: %v", difficulty, err)
			}
		})
	}
}

type cancelAfterChecks struct {
	context.Context
	cancel    context.CancelFunc
	remaining int
}

func (c *cancelAfterChecks) Err() error {
	c.remaining--
	if c.remaining <= 0 {
		c.cancel()
	}
	return c.Context.Err()
}

func TestSearchCancellation(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	during, cancelDuring := context.WithCancel(context.Background())
	defer cancelDuring()
	for _, tc := range []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{"before search", canceled, context.Canceled},
		{"expired deadline", expired, context.DeadlineExceeded},
		{"during search", &cancelAfterChecks{during, cancelDuring, 200}, context.Canceled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			move, err := player(t, ai.Hard, 1).ChooseMove(tc.ctx, chess.NewInitialPosition())
			if !errors.Is(err, tc.err) || move != (chess.Move{}) {
				t.Errorf("Expected cancellation without a move, got %v: %v", move, err)
			}
		})
	}
}

func TestSearchTerminalAndInvalidPositions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		position chess.Position
		err      error
	}{
		{"invalid turn", chess.Position{}, chess.ErrInvalidColor},
		{"missing king", chess.Position{SideToMove: chess.White}, chess.ErrKingNotFound},
		{"checkmate", positionWith(t, chess.Black, map[chess.Square]chess.Piece{
			chess.F6: {Type: chess.King, Color: chess.White},
			chess.G7: {Type: chess.Queen, Color: chess.White},
			chess.H8: {Type: chess.King, Color: chess.Black},
		}), ai.ErrNoLegalMoves},
		{"stalemate", positionWith(t, chess.Black, map[chess.Square]chess.Piece{
			chess.F7: {Type: chess.King, Color: chess.White},
			chess.G6: {Type: chess.Queen, Color: chess.White},
			chess.H8: {Type: chess.King, Color: chess.Black},
		}), ai.ErrNoLegalMoves},
	} {
		t.Run(tc.name, func(t *testing.T) {
			move, err := player(t, ai.Medium, 1).ChooseMove(context.Background(), tc.position)
			if !errors.Is(err, tc.err) || move != (chess.Move{}) {
				t.Errorf("Expected error %v without a move, got %v: %v", tc.err, move, err)
			}
		})
	}
	if _, err := ai.NewPlayer(ai.Difficulty("unknown"), 1); !errors.Is(err, ai.ErrInvalidDifficulty) {
		t.Errorf("Expected invalid difficulty error, got %v", err)
	}
}

func TestSearchPromotionAndEnPassant(t *testing.T) {
	promotion := positionWith(t, chess.White, map[chess.Square]chess.Piece{
		chess.A1: {Type: chess.King, Color: chess.White},
		chess.H6: {Type: chess.King, Color: chess.Black},
		chess.B7: {Type: chess.Pawn, Color: chess.White},
	})
	enPassant := positionWith(t, chess.White, map[chess.Square]chess.Piece{
		chess.A1: {Type: chess.King, Color: chess.White},
		chess.H8: {Type: chess.King, Color: chess.Black},
		chess.E5: {Type: chess.Pawn, Color: chess.White},
		chess.D5: {Type: chess.Pawn, Color: chess.Black},
	})
	enPassant.EnPassant = chess.D6
	for _, tc := range []struct {
		name     string
		position chess.Position
		want     chess.Move
	}{
		{"promotion", promotion, chess.Move{From: chess.B7, To: chess.B8, Promotion: chess.Queen}},
		{"en passant", enPassant, chess.Move{From: chess.E5, To: chess.D6}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			move, err := player(t, ai.Easy, 1).ChooseMove(context.Background(), tc.position)
			if err != nil || move != tc.want {
				t.Fatalf("Expected special move %v, got %v: %v", tc.want, move, err)
			}
			if _, err := rules.PositionAfterMove(tc.position, move); err != nil {
				t.Fatalf("Expected legal special move: %v", err)
			}
		})
	}
}
