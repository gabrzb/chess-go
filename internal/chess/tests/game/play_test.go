package game_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func mustGameSnapshot(t *testing.T, g *game.Game) game.Snapshot {
	t.Helper()
	snapshot, err := g.Snapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %v", err)
	}
	return snapshot
}

func TestGamePlay(t *testing.T) {
	cases := []struct {
		name   string
		moves  []chess.Move
		turn   chess.Color
		status chess.Status
		check  bool
		pieces map[chess.Square]chess.Piece
	}{
		{
			name: "both colors and capture",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.E4, To: chess.D5},
			},
			turn: chess.Black, status: chess.StatusOngoing,
			pieces: map[chess.Square]chess.Piece{
				chess.E2: {}, chess.E4: {}, chess.D7: {},
				chess.D5: {Type: chess.Pawn, Color: chess.White},
			},
		},
		{
			name: "check with game ongoing",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.F1, To: chess.B5},
			},
			turn: chess.Black, status: chess.StatusOngoing, check: true,
			pieces: map[chess.Square]chess.Piece{
				chess.F1: {},
				chess.B5: {Type: chess.Bishop, Color: chess.White},
			},
		},
		{
			name: "black delivers checkmate",
			moves: []chess.Move{
				{From: chess.F2, To: chess.F3},
				{From: chess.E7, To: chess.E5},
				{From: chess.G2, To: chess.G4},
				{From: chess.D8, To: chess.H4},
			},
			turn: chess.White, status: chess.StatusCheckmate, check: true,
			pieces: map[chess.Square]chess.Piece{
				chess.D8: {},
				chess.H4: {Type: chess.Queen, Color: chess.Black},
			},
		},
		{
			name: "white delivers checkmate",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.F7, To: chess.F6},
				{From: chess.D2, To: chess.D4},
				{From: chess.G7, To: chess.G5},
				{From: chess.D1, To: chess.H5},
			},
			turn: chess.Black, status: chess.StatusCheckmate, check: true,
			pieces: map[chess.Square]chess.Piece{
				chess.D1: {},
				chess.H5: {Type: chess.Queen, Color: chess.White},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			for _, move := range tc.moves {
				before := mustGameSnapshot(t, g)
				result, err := g.Play(move)
				if err != nil {
					t.Fatalf("Failed to play %v: %v", move, err)
				}
				if result.Move != move {
					t.Errorf("Expected played move %v, got %v", move, result.Move)
				}
				if result.State.LastMove == nil || *result.State.LastMove != move {
					t.Errorf("Expected last move %v, got %v", move, result.State.LastMove)
				}
				if result.State.SideToMove == before.SideToMove {
					t.Error("Expected side to move to change after playing")
				}
				if after := mustGameSnapshot(t, g); !reflect.DeepEqual(after, result.State) {
					t.Error("Expected returned state to match current game state")
				}
			}
			snapshot := mustGameSnapshot(t, g)
			if snapshot.SideToMove != tc.turn || snapshot.Status != tc.status || snapshot.Check != tc.check {
				t.Errorf("Expected turn/status/check %v/%v/%v, got %v/%v/%v",
					tc.turn, tc.status, tc.check, snapshot.SideToMove, snapshot.Status, snapshot.Check)
			}
			for square, want := range tc.pieces {
				if got, valid := snapshot.Board.PieceAt(square); !valid || got != want {
					t.Errorf("At %v: expected %v, got %v", square, want, got)
				}
			}
			if tc.status == chess.StatusCheckmate {
				result, err := g.Play(chess.Move{From: chess.A2, To: chess.A3})
				if !errors.Is(err, game.ErrGameOver) {
					t.Errorf("Expected game over error, got %v", err)
				}
				if result != (game.MoveResult{}) {
					t.Error("Expected empty result after rejected move")
				}
				if after := mustGameSnapshot(t, g); !reflect.DeepEqual(after, snapshot) {
					t.Error("Expected rejected move to preserve finished game")
				}
			}
		})
	}
}

func TestGamePlayErrorsPreserveState(t *testing.T) {
	cases := []struct {
		name string
		move chess.Move
		want error
	}{
		{"invalid source", chess.Move{From: chess.NoSquare, To: chess.E5}, chess.ErrInvalidSquare},
		{"invalid destination", chess.Move{From: chess.E7, To: chess.NoSquare}, chess.ErrInvalidSquare},
		{"empty source", chess.Move{From: chess.E6, To: chess.E5}, chess.ErrEmptySource},
		{"opponent piece", chess.Move{From: chess.D2, To: chess.D4}, chess.ErrIllegalMove},
		{"illegal destination", chess.Move{From: chess.E7, To: chess.E4}, chess.ErrIllegalMove},
		{"premature promotion", chess.Move{From: chess.E7, To: chess.E5, Promotion: chess.Queen}, chess.ErrIllegalMove},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); err != nil {
				t.Fatalf("Failed to play setup move: %v", err)
			}
			before := mustGameSnapshot(t, g)
			result, err := g.Play(tc.move)
			if !errors.Is(err, tc.want) {
				t.Errorf("Expected error %v, got %v", tc.want, err)
			}
			if result != (game.MoveResult{}) {
				t.Error("Expected empty result after rejected move")
			}
			if after := mustGameSnapshot(t, g); !reflect.DeepEqual(after, before) {
				t.Error("Expected rejected move to preserve game state")
			}
		})
	}
}

func TestGamePlaySnapshotIndependence(t *testing.T) {
	g := game.NewGame()
	result, err := g.Play(chess.Move{From: chess.E2, To: chess.E4})
	if err != nil {
		t.Fatalf("Failed to play setup move: %v", err)
	}
	before := mustGameSnapshot(t, g)
	for name, snapshot := range map[string]game.Snapshot{
		"play result":     result.State,
		"snapshot method": mustGameSnapshot(t, g),
	} {
		t.Run(name, func(t *testing.T) {
			if snapshot.LastMove == nil {
				t.Fatal("Expected last move after playing")
			}
			snapshot.LastMove.To = chess.H8
			if _, err := snapshot.Board.RemovePiece(chess.E4); err != nil {
				t.Fatalf("Failed to change snapshot board: %v", err)
			}
			if after := mustGameSnapshot(t, g); !reflect.DeepEqual(after, before) {
				t.Error("Expected changing returned snapshot to preserve game state")
			}
		})
	}
}
