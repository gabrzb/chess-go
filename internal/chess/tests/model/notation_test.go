package model_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func checkParseSquare(t *testing.T, input string, expected chess.Square) {
	got, err := chess.ParseSquare(input)
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("ParseSquare(%s) expects %d, got %d instead", input, expected, got)
	}
}

func checkParseSquareInvalid(t *testing.T, input string) {
	got, err := chess.ParseSquare(input)

	if !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("ParseSquare(%s) expects ErrInvalidSquare, got %v instead", input, err)
	}

	if got != chess.NoSquare {
		t.Fatalf("ParseSquare(%s) expects NoSquare, got %d instead", input, got)
	}
}

func TestParseSquare(t *testing.T) {
	checkParseSquare(t, "a1", 0)
	checkParseSquare(t, "h8", 63)
}

func TestParseSquareInvalid(t *testing.T) {
	checkParseSquareInvalid(t, "a0")
	checkParseSquareInvalid(t, "h9")
}

func TestSquareString(t *testing.T) {
	if got := chess.Square(28).String(); got != "e4" {
		t.Errorf("Square(28).String() expects e4, got %s instead", got)
	}

	if got := chess.NoSquare.String(); got != "-" {
		t.Errorf("NoSquare.String() expects -, got %s instead", got)
	}
}

func TestSquareValid(t *testing.T) {
	testA := chess.Square(0).Valid()
	testB := chess.Square(63).Valid()
	testC := chess.NoSquare.Valid()

	if !testA {
		t.Errorf("Square(0).Valid() expects true, got false instead")
	}

	if !testB {
		t.Errorf("Square(63).Valid() expects true, got false instead")
	}

	if testC {
		t.Errorf("NoSquare.Valid() expects false, got true instead")
	}
}

func checkMoveString(t *testing.T, move chess.Move, expected string) {
	got := move.String()

	if got != expected {
		t.Errorf("Move.String() = %s; expected %s", got, expected)
	}
}

func checkMovePromotionString(t *testing.T, move chess.Move, expected string) {
	got := move.String()

	if got != expected {
		t.Errorf("Move.String() with promotion = %s; expected %s", got, expected)
	}
}

func TestMoveString(t *testing.T) {
	checkMoveString(
		t,
		chess.Move{From: chess.Square(12), To: chess.Square(28)},
		"e2e4",
	)
	checkMovePromotionString(
		t,
		chess.Move{
			From: chess.Square(52), To: chess.Square(60), Promotion: chess.Knight,
		},
		"e7e8n",
	)
	checkMovePromotionString(
		t,
		chess.Move{
			From: chess.Square(52), To: chess.Square(60), Promotion: chess.Bishop,
		},
		"e7e8b",
	)
	checkMovePromotionString(
		t,
		chess.Move{
			From: chess.Square(52), To: chess.Square(60), Promotion: chess.Rook,
		},
		"e7e8r",
	)
	checkMovePromotionString(
		t,
		chess.Move{
			From: chess.Square(52), To: chess.Square(60), Promotion: chess.Queen,
		},
		"e7e8q",
	)
}
