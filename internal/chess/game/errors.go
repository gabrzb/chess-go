package game

import "errors"

var (
	ErrGameOver            = errors.New("game is over")
	ErrInvalidHistoryIndex = errors.New("invalid history index")
)
