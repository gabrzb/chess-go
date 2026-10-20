package rules

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func EvaluateStatus(position chess.Position) (chess.Status, error) {
	canMove, err := HasLegalMove(position)
	if err != nil {
		return chess.StatusOngoing, err
	}

	if canMove {
		if InsufficientMaterial(position.Board) {
			return chess.StatusInsufficientMaterial, nil
		}
		if position.HalfmoveClock >= 150 {
			return chess.StatusSeventyFiveMoves, nil
		}
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
