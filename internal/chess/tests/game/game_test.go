package game_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func TestNewGameSnapshot(t *testing.T) {
	g := game.NewGame()
	snapshot, err := g.Snapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %v", err)
	}

	initial := chess.NewInitialPosition()
	if snapshot.Board != initial.Board {
		t.Error("Expected initial board in new game snapshot")
	}
	if snapshot.SideToMove != chess.White {
		t.Errorf("Expected White to move, got %v", snapshot.SideToMove)
	}
	if snapshot.Status != chess.StatusOngoing {
		t.Errorf("Expected ongoing game, got %v", snapshot.Status)
	}
	if snapshot.Check {
		t.Error("Expected no check in initial position")
	}
	if snapshot.LastMove != nil {
		t.Errorf("Expected no last move in new game, got %v", *snapshot.LastMove)
	}
}

func TestGameLegalMovesFrom(t *testing.T) {
	cases := []struct {
		name    string
		from    chess.Square
		want    []chess.Move
		wantErr error
	}{
		{
			name: "pawn initial moves",
			from: chess.E2,
			want: []chess.Move{
				{From: chess.E2, To: chess.E3},
				{From: chess.E2, To: chess.E4},
			},
		},
		{
			name: "knight jumps",
			from: chess.B1,
			want: []chess.Move{
				{From: chess.B1, To: chess.A3},
				{From: chess.B1, To: chess.C3},
			},
		},
		{name: "blocked rook", from: chess.A1},
		{name: "empty square", from: chess.E4, wantErr: chess.ErrEmptySource},
		{name: "opponent piece", from: chess.E7, wantErr: chess.ErrIllegalMove},
		{name: "invalid square", from: chess.NoSquare, wantErr: chess.ErrInvalidSquare},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			before, err := g.Snapshot()
			if err != nil {
				t.Fatalf("Failed to get snapshot: %v", err)
			}

			moves, err := g.LegalMovesFrom(tc.from)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Expected error %v, got %v", tc.wantErr, err)
			}
			if len(moves) != len(tc.want) {
				t.Errorf("Expected %d moves, got %v", len(tc.want), moves)
			}
			for _, move := range tc.want {
				if !slices.Contains(moves, move) {
					t.Errorf("Expected legal move %v, got %v", move, moves)
				}
			}

			after, err := g.Snapshot()
			if err != nil {
				t.Fatalf("Failed to get snapshot after query: %v", err)
			}
			if after != before {
				t.Error("Expected move query to preserve game state")
			}
		})
	}
}

func TestGameSnapshotBoardIsIndependent(t *testing.T) {
	g := game.NewGame()
	before, err := g.Snapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %v", err)
	}

	snapshot := before
	if _, err := snapshot.Board.RemovePiece(chess.E2); err != nil {
		t.Fatalf("Failed to change snapshot board: %v", err)
	}

	after, err := g.Snapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot after change: %v", err)
	}
	if after != before {
		t.Error("Expected changing snapshot board to preserve game state")
	}
}
