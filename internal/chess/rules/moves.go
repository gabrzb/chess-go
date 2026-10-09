package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

// BoardAfterMove applies a generated basic move to a copy of the board.
func BoardAfterMove(board chess.Board, move chess.Move) (chess.Board, error) {
	if !move.From.Valid() || !move.To.Valid() {
		return chess.Board{}, chess.ErrInvalidSquare
	}

	piece, _ := board.PieceAt(move.From)
	if piece.Type == chess.NoPiece {
		return chess.Board{}, chess.ErrEmptySource
	}

	if _, err := board.RemovePiece(move.From); err != nil {
		return chess.Board{}, err
	}
	if err := board.SetPiece(move.To, piece); err != nil {
		return chess.Board{}, err
	}

	return board, nil
}

// LegalMoves filters basic moves that leave the moving side's king in check.
func LegalMoves(board chess.Board, from chess.Square) ([]chess.Move, error) {
	piece, ok := board.PieceAt(from)
	if !ok {
		return nil, chess.ErrInvalidSquare
	}

	if piece.Type == chess.NoPiece {
		return nil, chess.ErrEmptySource
	}

	if _, err := IsInCheck(board, piece.Color); err != nil {
		return nil, err
	}

	var moves []chess.Move

	for _, move := range movegen.PseudoLegalMoves(board, from) {
		target, _ := board.PieceAt(move.To)
		if target.Type == chess.King {
			continue
		}

		nextMove, err := BoardAfterMove(board, move)
		if err != nil {
			return nil, err
		}

		inCheck, err := IsInCheck(nextMove, piece.Color)
		if err != nil {
			return nil, err
		}

		if !inCheck {
			moves = append(moves, move)
		}
	}

	return moves, nil
}

// HasLegalMove reports whether any piece of the given color has a legal move.
func HasLegalMove(board chess.Board, color chess.Color) (bool, error) {
	if _, err := IsInCheck(board, color); err != nil {
		return false, err
	}

	for from := chess.Square(0); from < chess.NoSquare; from++ {
		piece, _ := board.PieceAt(from)

		if piece.Type == chess.NoPiece || piece.Color != color {
			continue
		}

		moves, err := LegalMoves(board, from)
		if err != nil {
			return false, err
		}

		if len(moves) > 0 {
			return true, nil
		}
	}

	return false, nil
}
