package chess

type Square uint8

const NoSquare Square = SquareCount

func (s Square) String() string {
	if s >= NoSquare {
		return "-"
	}

	col := s.Column()
	row := s.Row()

	letters := "abcdefgh"
	digits := "12345678"

	return string([]byte{letters[col], digits[row]})
}

func (s Square) Valid() bool {
	return s < NoSquare
}

// Column returns the zero-based column, or -1 for an invalid square.
func (s Square) Column() int {
	if !s.Valid() {
		return -1
	}
	return int(s) % BoardSize
}

// Row returns the zero-based row, or -1 for an invalid square.
func (s Square) Row() int {
	if !s.Valid() {
		return -1
	}
	return int(s) / BoardSize
}

// Offset moves by columns and rows, returning false when it leaves the board.
func (s Square) Offset(colStep, rowStep int) (Square, bool) {
	if !s.Valid() {
		return NoSquare, false
	}
	col := s.Column() + colStep
	row := s.Row() + rowStep
	if col < 0 || col >= BoardSize || row < 0 || row >= BoardSize {
		return NoSquare, false
	}
	return Square(row*BoardSize + col), true
}

func ParseSquare(text string) (Square, error) {
	var col, row int

	if len(text) != 2 {
		return NoSquare, ErrInvalidSquare
	}

	if text[0] < 'a' || text[0] > 'h' || text[1] < '1' || text[1] > '8' {
		return NoSquare, ErrInvalidSquare
	}

	col = int(text[0] - 'a')
	row = int(text[1] - '1')

	return Square(row*BoardSize + col), nil
}
