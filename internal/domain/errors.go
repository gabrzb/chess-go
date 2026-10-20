package domain

import (
	"errors"
)

var (
	ErrInvalidPosition = errors.New("invalid position")
	ErrInvalidSquare   = errors.New("invalid square")
	ErrInvalidColor    = errors.New("invalid color")
	ErrKingNotFound    = errors.New("king not found")
	ErrEmptySource     = errors.New("empty source square")
	ErrIllegalMove     = errors.New("illegal move")
)
