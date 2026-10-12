package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

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

	if move.Promotion != chess.NoPiece {
		piece.Type = move.Promotion
	}

	if err := board.SetPiece(move.To, piece); err != nil {
		return chess.Board{}, err
	}

	return board, nil
}

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

func LegalMovesFrom(position chess.Position, from chess.Square) ([]chess.Move, error) {
	piece, ok := position.Board.PieceAt(from)

	if !ok {
		return nil, chess.ErrInvalidSquare
	}

	if piece.Type == chess.NoPiece {
		return nil, chess.ErrEmptySource
	}

	if position.SideToMove != chess.White && position.SideToMove != chess.Black {
		return nil, chess.ErrInvalidColor
	}

	if piece.Color != position.SideToMove {
		return nil, chess.ErrIllegalMove
	}

	moves, err := LegalMoves(position.Board, from)
	if err != nil {
		return nil, err
	}

	if piece.Type == chess.Pawn {
		if move, ok := EnPassantMove(position, from); ok {
			moves = append(moves, move)
		}
	}

	if piece.Type != chess.King {
		return moves, nil
	}

	generators := []func(chess.Position) (chess.Move, bool){
		KingSideCastlingMove,
		QueenSideCastlingMove,
	}

	for _, generate := range generators {
		move, ok := generate(position)

		if ok && move.From == from {
			moves = append(moves, move)
		}
	}

	return moves, nil
}

func HasLegalMove(board chess.Board, color chess.Color) (bool, error) {
	if _, err := IsInCheck(board, color); err != nil {
		return false, err
	}

	for from := chess.A1; from < chess.NoSquare; from++ {
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
