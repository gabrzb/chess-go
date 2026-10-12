package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

func EvaluateStatus(position chess.Position) (chess.Status, error) {
	canMove, err := HasLegalMove(position)
	if err != nil {
		return chess.StatusOngoing, err
	}

	if canMove {
		return chess.StatusOngoing, nil
	}

	inCheck, err := IsInCheck(position.Board, position.SideToMove)
	if err != nil {
		return chess.StatusOngoing, err
	}

	if inCheck {
		return chess.StatusCheckmate, nil
	}

	return chess.StatusStalemate, nil
}
