package game_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func snapshot(t *testing.T, g *game.Game) game.Snapshot {
	t.Helper()
	state, err := g.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func play(t *testing.T, g *game.Game, moves ...chess.Move) {
	t.Helper()
	for _, move := range moves {
		if _, err := g.Play(move); err != nil {
			t.Fatalf("Playing %v: %v", move, err)
		}
	}
}

func TestGamePlayAndHistory(t *testing.T) {
	g := game.NewGame()
	wantPosition := chess.NewInitialPosition()
	states := []game.Snapshot{snapshot(t, g)}
	for _, move := range []chess.Move{{From: chess.E2, To: chess.E4}, {From: chess.D7, To: chess.D5}, {From: chess.E4, To: chess.D5}} {
		result, err := g.Play(move)
		if err != nil {
			t.Fatal(err)
		}
		wantPosition, err = rules.PositionAfterMove(wantPosition, move)
		if err != nil {
			t.Fatal(err)
		}
		states = append(states, result.State)
		if result.Move != move || result.State.LastMove == nil || *result.State.LastMove != move || !reflect.DeepEqual(result.State, snapshot(t, g)) || g.Position() != wantPosition {
			t.Fatalf("Incorrect state after %v", move)
		}
	}
	if g.HistoryLen() != len(states) {
		t.Fatal("Incorrect history length")
	}
	before := snapshot(t, g)
	for index, want := range states {
		got, err := g.HistorySnapshot(index)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("History %d changed: %v", index, err)
		}
	}
	for _, index := range []int{-1, len(states)} {
		if _, err := g.HistorySnapshot(index); !errors.Is(err, game.ErrInvalidHistoryIndex) {
			t.Errorf("Invalid index %d accepted", index)
		}
	}
	if !reflect.DeepEqual(snapshot(t, g), before) {
		t.Error("History navigation changed the live game")
	}
	if moves, err := g.LegalMovesFrom(chess.E7); err != nil || len(moves) != 2 {
		t.Errorf("Incorrect move query: %v, %v", moves, err)
	}
	play(t, g, chess.Move{From: chess.E7, To: chess.E5})
	if g.HistoryLen() != len(states)+1 {
		t.Error("Continuing after review truncated history")
	}
}

func TestGameRejectionsAndCheckmate(t *testing.T) {
	g := game.NewGame()
	before := g.Position()
	for _, move := range []chess.Move{{From: chess.E2, To: chess.E5}, {From: chess.E7, To: chess.E5}, {From: chess.NoSquare, To: chess.E4}} {
		if _, err := g.Play(move); err == nil || g.Position() != before || g.HistoryLen() != 1 {
			t.Errorf("Rejected move %v changed the game or was accepted", move)
		}
	}
	play(t, g, chess.Move{From: chess.F2, To: chess.F3}, chess.Move{From: chess.E7, To: chess.E5}, chess.Move{From: chess.G2, To: chess.G4}, chess.Move{From: chess.D8, To: chess.H4})
	final := snapshot(t, g)
	if final.Status != chess.StatusCheckmate || !final.Check || final.SideToMove != chess.White {
		t.Fatal("Checkmate was not recorded")
	}
	if _, err := g.Play(chess.Move{From: chess.A2, To: chess.A3}); !errors.Is(err, game.ErrGameOver) || g.HistoryLen() != 5 || !reflect.DeepEqual(final, snapshot(t, g)) {
		t.Error("Finished game accepted a move or changed state")
	}
}

func TestGameSnapshotIsolation(t *testing.T) {
	g := game.NewGame()
	result, err := g.Play(chess.Move{From: chess.E2, To: chess.E4})
	if err != nil {
		t.Fatal(err)
	}
	play(t, g, chess.Move{From: chess.E7, To: chess.E5})
	before := g.Position()
	current := snapshot(t, g)
	saved, err := g.HistorySnapshot(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, copy := range []game.Snapshot{result.State, current, saved} {
		copy.LastMove.To = chess.H8
		if _, err := copy.Board.RemovePiece(chess.E1); err != nil {
			t.Fatal(err)
		}
	}
	copy := g.Position()
	copy.CastlingRights = 0
	copy.EnPassant = chess.H3
	if _, err := copy.Board.RemovePiece(chess.E8); err != nil {
		t.Fatal(err)
	}
	history, err := g.HistorySnapshot(1)
	if err != nil || history.LastMove == nil || *history.LastMove != result.Move || history.Board != result.State.Board || g.Position() != before || g.HistoryLen() != 3 {
		t.Error("Returned values alias the live game or stored history")
	}
	if latest := snapshot(t, g).LastMove; latest == nil || *latest != (chess.Move{From: chess.E7, To: chess.E5}) {
		t.Error("Changing a returned last move modified the live game")
	}
}
