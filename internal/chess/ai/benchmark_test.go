package ai

import (
	"context"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func BenchmarkChooseMove(b *testing.B) {
	for _, difficulty := range []Difficulty{Medium, Hard} {
		b.Run(string(difficulty), func(b *testing.B) {
			player, err := NewPlayer(difficulty, 1)
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
