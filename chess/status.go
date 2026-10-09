package chess

type Status uint8

const (
	StatusOngoing Status = iota
	StatusCheckmate
	StatusStalemate
)

func EvaluateStatus(board Board, sideToMove Color) (Status, error) {
	canMove, err := hasLegalMove(board, sideToMove)
	if err != nil {
		return StatusOngoing, err
	}

	if canMove {
		return StatusOngoing, nil
	}

	inCheck, err := IsInCheck(board, sideToMove)
	if err != nil {
		return StatusOngoing, err
	}

	if inCheck {
		return StatusCheckmate, nil
	}

	return StatusStalemate, nil
}
