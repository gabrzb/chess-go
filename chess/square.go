package chess

import (
	"errors"
)

type Square uint8

const NoSquare Square = 64

func (s Square) String() string {
	if s >= NoSquare {
		return "-"
	}

	col := s % 8
	row := s / 8

	letters := "abcdefgh"
	digits := "12345678"

	return string([]byte{letters[col], digits[row]})
}

func (s Square) Valid() bool {
	return s < NoSquare
}

func ParseSquare(text string) (Square, error) {
	var col, row int

	if len(text) != 2 {
		return NoSquare, errors.New("invalid square format")
	}

	if text[0] < 'a' || text[0] > 'h' || text[1] < '1' || text[1] > '8' {
		return NoSquare, errors.New("invalid square format")
	}

	col = int(text[0] - 'a')
	row = int(text[1] - '1')

	return Square(row*8 + col), nil
}
