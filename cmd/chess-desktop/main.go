package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrzb/chess-go/internal/delivery/desktop"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/engine"
)

func main() {
	difficulty := flag.String("ai", "off", "AI opponent: off, easy, medium or hard")
	humanColor := flag.String("color", "white", "Your color against AI: white or black")
	seed := flag.Int64("seed", 1, "Seed for ties between equally evaluated AI moves")
	flag.Parse()

	options := desktop.Options{Difficulty: engine.Difficulty(*difficulty), HumanColor: chess.White, Seed: *seed}
	if *difficulty == "off" {
		options.Difficulty = ""
	}
	if *humanColor == "black" {
		options.HumanColor = chess.Black
	}
	var err error
	switch {
	case flag.NArg() > 0:
		err = fmt.Errorf("unexpected arguments: %v", flag.Args())
	case *humanColor != "white" && *humanColor != "black":
		err = fmt.Errorf("color must be white or black")
	default:
		err = desktop.Run(options)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to run desktop app:", err)
		os.Exit(1)
	}
}
