package ai_test

import (
	"context"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func BenchmarkChooseMove(b *testing.B) {
	for _, difficulty := range []ai.Difficulty{ai.Medium, ai.Hard} {
		b.Run(string(difficulty), func(b *testing.B) {
			player, err := ai.NewPlayer(difficulty, 1)
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
