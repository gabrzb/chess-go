package game_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func TestGameHistoryNavigation(t *testing.T) {
	cases := []struct {
		name  string
		moves []chess.Move
	}{
		{name: "initial position"},
		{
			name: "both colors and capture",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.E4, To: chess.D5},
			},
		},
		{
			name: "check followed by defense",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.F1, To: chess.B5},
				{From: chess.C7, To: chess.C6},
			},
		},
		{
			name: "finished game",
			moves: []chess.Move{
				{From: chess.F2, To: chess.F3},
				{From: chess.E7, To: chess.E5},
				{From: chess.G2, To: chess.G4},
				{From: chess.D8, To: chess.H4},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			want := []game.Snapshot{mustGameSnapshot(t, g)}
			if got := g.HistoryLen(); got != len(want) {
				t.Fatalf("Expected one initial history position, got %d", got)
			}
			for _, move := range tc.moves {
				result, err := g.Play(move)
				if err != nil {
					t.Fatalf("Failed to play setup move %v: %v", move, err)
				}
				want = append(want, result.State)
				if got := g.HistoryLen(); got != len(want) {
					t.Fatalf("Expected %d history positions, got %d", len(want), got)
				}
			}

			current := mustGameSnapshot(t, g)
			var indices []int
			for index := len(want) - 1; index >= 0; index-- {
				indices = append(indices, index)
			}
			for index := range want {
				indices = append(indices, index)
			}
			for _, index := range indices {
				got, err := g.HistorySnapshot(index)
				if err != nil {
					t.Fatalf("Failed to get history index %d: %v", index, err)
				}
				if !reflect.DeepEqual(got, want[index]) {
					t.Errorf("Expected recorded snapshot at history index %d, got %+v", index, got)
				}
				if g.HistoryLen() != len(want) || !reflect.DeepEqual(mustGameSnapshot(t, g), current) {
					t.Fatal("Expected navigation to preserve history length and current game state")
				}
			}

			wantErr := chess.ErrInvalidSquare
			if current.Status != chess.StatusOngoing {
				wantErr = game.ErrGameOver
			}
			if _, err := g.Play(chess.Move{From: chess.NoSquare, To: chess.NoSquare}); !errors.Is(err, wantErr) {
				t.Errorf("Expected rejected move error %v, got %v", wantErr, err)
			}
			if g.HistoryLen() != len(want) || !reflect.DeepEqual(mustGameSnapshot(t, g), current) {
				t.Fatal("Expected rejected move to preserve history length and current game state")
			}

			if current.Status == chess.StatusOngoing {
				move := chess.Move{From: chess.A2, To: chess.A3}
				if current.SideToMove == chess.Black {
					move = chess.Move{From: chess.A7, To: chess.A6}
				}
				result, err := g.Play(move)
				if err != nil {
					t.Fatalf("Failed to continue current game after review: %v", err)
				}
				want = append(want, result.State)
				if g.HistoryLen() != len(want) {
					t.Fatalf("Expected %d positions after continuing game, got %d", len(want), g.HistoryLen())
				}
			}
			for index, expected := range want {
				got, err := g.HistorySnapshot(index)
				if err != nil || !reflect.DeepEqual(got, expected) {
					t.Errorf("Expected preserved history position %d, got %+v and error %v", index, got, err)
				}
			}
		})
	}
}

func TestGameHistoryInvalidIndices(t *testing.T) {
	for _, played := range []bool{false, true} {
		name := "initial game"
		if played {
			name = "after move"
		}
		t.Run(name, func(t *testing.T) {
			g := game.NewGame()
			if played {
				if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); err != nil {
					t.Fatalf("Failed to play setup move: %v", err)
				}
			}
			before := mustGameSnapshot(t, g)
			length := g.HistoryLen()
			for _, index := range []int{-1, length, length + 1} {
				got, err := g.HistorySnapshot(index)
				if !errors.Is(err, game.ErrInvalidHistoryIndex) {
					t.Errorf("At index %d: expected invalid history index error, got %v", index, err)
				}
				if got != (game.Snapshot{}) {
					t.Errorf("At index %d: expected empty snapshot, got %+v", index, got)
				}
				if g.HistoryLen() != length || !reflect.DeepEqual(mustGameSnapshot(t, g), before) {
					t.Fatal("Expected invalid query to preserve history length and current game state")
				}
			}
		})
	}
}

func TestGameHistorySnapshotIndependence(t *testing.T) {
	g := game.NewGame()
	for _, move := range []chess.Move{
		{From: chess.E2, To: chess.E4},
		{From: chess.E7, To: chess.E5},
	} {
		if _, err := g.Play(move); err != nil {
			t.Fatalf("Failed to play setup move %v: %v", move, err)
		}
	}
	current := mustGameSnapshot(t, g)
	var want []game.Snapshot
	for index := 0; index < g.HistoryLen(); index++ {
		snapshot, err := g.HistorySnapshot(index)
		if err != nil {
			t.Fatalf("Failed to get history index %d: %v", index, err)
		}
		want = append(want, snapshot)
	}
	for index := range want {
		snapshot, err := g.HistorySnapshot(index)
		if err != nil {
			t.Fatalf("Failed to get history index %d: %v", index, err)
		}
		if _, err := snapshot.Board.RemovePiece(chess.E1); err != nil {
			t.Fatalf("Failed to change returned history board: %v", err)
		}
		if snapshot.LastMove != nil {
			snapshot.LastMove.To = chess.H8
		}
		for savedIndex, expected := range want {
			got, err := g.HistorySnapshot(savedIndex)
			if err != nil || !reflect.DeepEqual(got, expected) {
				t.Errorf("Changing snapshot %d affected history position %d: snapshot %+v, error %v", index, savedIndex, got, err)
			}
		}
		if g.HistoryLen() != len(want) || !reflect.DeepEqual(mustGameSnapshot(t, g), current) {
			t.Fatal("Expected changing returned snapshot to preserve history length and current game state")
		}
	}
}
