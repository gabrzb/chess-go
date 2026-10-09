package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

// EvaluateStatus classifies the position for the side whose turn it is.
func EvaluateStatus(board chess.Board, sideToMove chess.Color) (chess.Status, error) {
	canMove, err := HasLegalMove(board, sideToMove)
	if err != nil {
		return chess.StatusOngoing, err
	}

	if canMove {
		return chess.StatusOngoing, nil
	}

	inCheck, err := IsInCheck(board, sideToMove)
	if err != nil {
		return chess.StatusOngoing, err
	}

	if inCheck {
		return chess.StatusCheckmate, nil
	}

	return chess.StatusStalemate, nil
}
