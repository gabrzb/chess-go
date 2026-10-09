package chess

type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

func (m Move) String() string {
	squares := m.From.String() + m.To.String()

	switch m.Promotion {
	case Knight:
		return squares + "n"
	case Bishop:
		return squares + "b"
	case Rook:
		return squares + "r"
	case Queen:
		return squares + "q"
	}

	return squares
}
