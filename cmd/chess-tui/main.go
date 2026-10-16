package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/tui"
)

func main() {
	difficulty := flag.String("ai", "off", "AI opponent: off, easy, medium or hard")
	humanColor := flag.String("color", "white", "Your color against AI: white or black")
	seed := flag.Int64("seed", 1, "Seed for ties between equally evaluated AI moves")
	noAnimations := flag.Bool("no-animation", false, "Disable UI animations")
	flag.Parse()
	direct := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "ai" {
			direct = true
		}
	})

	human := chess.White
	var err error
	if flag.NArg() > 0 {
		err = fmt.Errorf("unexpected arguments: %v", flag.Args())
	} else if *humanColor != "white" && *humanColor != "black" {
		err = fmt.Errorf("color must be white or black")
	} else {
		if *humanColor == "black" {
			human = chess.Black
		}
		level := ai.Difficulty(*difficulty)
		if *difficulty == "off" {
			level = ""
		}
		err = tui.RunWithOptions(tui.Options{Difficulty: level, HumanColor: human, Seed: *seed, StartGame: direct, NoAnimations: *noAnimations})
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to run TUI:", err)
		os.Exit(1)
	}
}
