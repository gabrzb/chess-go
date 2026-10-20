package ai_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
)

func TestMoveBudgetReturnsLegalMove(t *testing.T) {
	board := chess.NewInitialPosition().Board
	// A crowded tactical position exercises the budget without changing private fields.
	for file := range chess.BoardSize {
		white, _ := chess.A2.Offset(file, 0)
		black, _ := chess.A7.Offset(file, 0)
		setPiece(t, &board, white, chess.Queen, chess.White)
		setPiece(t, &board, black, chess.Queen, chess.Black)
	}
	position, err := chess.NewPosition(board, chess.White)
	if err != nil {
		t.Fatal(err)
	}
	moves := allMoves(t, position)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	move, err := player(t, ai.Hard, 1).ChooseMove(ctx, position)
	if err != nil || !slices.Contains(moves, move) {
		t.Fatalf("budget did not return a legal move: %v, %v", move, err)
	}
}

func TestParentDeadlineDoesNotReturnFallback(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	defer cancel()
	move, err := player(t, ai.Hard, 1).ChooseMove(ctx, chess.NewInitialPosition())
	if !errors.Is(err, context.DeadlineExceeded) || move != (chess.Move{}) {
		t.Fatal("parent cancellation returned a move", move, err)
	}
}
