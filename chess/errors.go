package chess

import "errors"

var (
	ErrInvalidSquare = errors.New("invalid square")
	ErrInvalidColor  = errors.New("invalid color")
	ErrKingNotFound  = errors.New("king not found")
)
