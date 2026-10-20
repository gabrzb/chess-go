package domain

type Position struct {
	Board          Board
	SideToMove     Color
	CastlingRights CastlingRights
	EnPassant      Square
	HalfmoveClock  uint16
	FullmoveNumber uint16
}

func NewPosition(board Board, sideToMove Color) (Position, error) {
	if sideToMove != White && sideToMove != Black {
		return Position{}, ErrInvalidColor
	}

	return Position{
		Board:          board,
		SideToMove:     sideToMove,
		EnPassant:      NoSquare,
		FullmoveNumber: 1,
	}, nil
}
