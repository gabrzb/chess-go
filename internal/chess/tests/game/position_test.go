package game_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestGamePositionPreservesFullStateAndIsolation(t *testing.T) {
	g := game.NewGame()
	want := chess.NewInitialPosition()
	for _, move := range []chess.Move{
		{From: chess.E2, To: chess.E4},
		{From: chess.E7, To: chess.E5},
		{From: chess.E1, To: chess.E2},
	} {
		var err error
		want, err = rules.PositionAfterMove(want, move)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
		if got := g.Position(); got != want {
			t.Errorf("Expected complete board, turn, castling, en passant and clocks after %v", move)
		}
	}
	copy := g.Position()
	if _, err := copy.Board.RemovePiece(chess.E2); err != nil {
		t.Fatal(err)
	}
	copy.CastlingRights = 0
	copy.EnPassant = chess.H3
	copy.SideToMove = chess.White
	if g.Position() != want || g.HistoryLen() != 4 {
		t.Error("Expected changing the search position to preserve the live game and history")
	}
}
