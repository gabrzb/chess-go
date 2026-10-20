package ai

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestMoveBudgetReturnsLegalFallback(t *testing.T) {
	p, _ := NewPlayer(Hard, 1)
	p.moveTime = time.Nanosecond
	position := chess.NewInitialPosition()
	moves, err := legalMoves(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}
	move, err := p.ChooseMove(context.Background(), position)
	if err != nil || !slices.Contains(moves, move) {
		t.Fatalf("timeout did not return a legal move: %v, %v", move, err)
	}
}

func TestParentDeadlineDoesNotReturnFallback(t *testing.T) {
	p, _ := NewPlayer(Hard, 1)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	move, err := p.ChooseMove(ctx, chess.NewInitialPosition())
	if !errors.Is(err, context.DeadlineExceeded) || move != (chess.Move{}) {
		t.Fatal("parent cancellation returned a move", move, err)
	}
}
