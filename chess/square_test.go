package chess

import (
	"errors"
	"testing"
)

func checkParseSquare(t *testing.T, input string, expected Square) {
	got, err := ParseSquare(input)
	if err != nil {
		t.Fatal(err)
	}

	if got != expected {
		t.Errorf("ParseSquare(%s) expects %d, got %d instead", input, expected, got)
	}
}

func checkParseSquareInvalid(t *testing.T, input string) {
	got, err := ParseSquare(input)

	if !errors.Is(err, ErrInvalidSquare) {
		t.Fatalf("ParseSquare(%s) expects ErrInvalidSquare, got %v instead", input, err)
	}

	if got != NoSquare {
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
	if got := Square(28).String(); got != "e4" {
		t.Errorf("Square(28).String() expects e4, got %s instead", got)
	}

	if got := NoSquare.String(); got != "-" {
		t.Errorf("NoSquare.String() expects -, got %s instead", got)
	}
}

func TestSquareValid(t *testing.T) {
	testA := Square(0).Valid()
	testB := Square(63).Valid()
	testC := NoSquare.Valid()

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
