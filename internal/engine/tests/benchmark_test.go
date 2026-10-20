package engine_test

import (
	"context"
	"github.com/gabrzb/chess-go/internal/engine"
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func BenchmarkChooseMove(b *testing.B) {
	for _, difficulty := range []engine.Difficulty{engine.Medium, engine.Hard} {
		b.Run(string(difficulty), func(b *testing.B) {
			player, err := engine.NewPlayer(difficulty, 1)
			if err != nil {
				b.Fatal(err)
			}
			position := chess.NewInitialPosition()
			b.ReportAllocs()
			for b.Loop() {
				if _, err := player.ChooseMove(context.Background(), position); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
