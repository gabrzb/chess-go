package rules

import (
	"slices"

	"github.com/gabrzb/chess-go/internal/chess"
)

func PositionAfterMove(position chess.Position, move chess.Move) (chess.Position, error) {
	if !move.From.Valid() || !move.To.Valid() {
		return chess.Position{}, chess.ErrInvalidSquare
	}

	moves, err := LegalMovesFrom(position, move.From)
	if err != nil {
		return chess.Position{}, err
	}

	if !slices.Contains(moves, move) {
		return chess.Position{}, chess.ErrIllegalMove
	}

	piece, _ := position.Board.PieceAt(move.From)

	if piece.Type == chess.King {
		for _, side := range []castlingSide{kingSide, queenSide} {
			castleMove, ok := castlingMove(position, side)

			if ok && castleMove == move {
				return applyCastling(position, side)
			}
		}
	}

	capturedPiece, _ := position.Board.PieceAt(move.To)
	var board chess.Board

	plan, ok := makeEnPassantPlan(position, move.From)
	if ok && plan.move == move {
		capturedPiece, _ = position.Board.PieceAt(plan.capturedSquare)
		board, err = boardAfterEnPassant(position.Board, plan)
	} else {
		board, err = BoardAfterMove(position.Board, move)
	}

	if err != nil {
		return chess.Position{}, err
	}

	position.Board = board
	position.CastlingRights = castlingRightsWithoutPiece(
		position.CastlingRights,
		piece,
		move.From,
	)
	position.CastlingRights = castlingRightsWithoutPiece(
		position.CastlingRights,
		capturedPiece,
		move.To,
	)

	if piece.Type == chess.Pawn || capturedPiece.Type != chess.NoPiece {
		position.HalfmoveClock = 0
	} else {
		position.HalfmoveClock++
	}

	position.EnPassant = chess.NoSquare

	if piece.Type == chess.Pawn {
		rowDistance := move.To.Row() - move.From.Row()
		if rowDistance == 2 || rowDistance == -2 {
			position.EnPassant, _ = move.From.Offset(0, rowDistance/2)
		}
	}

	if position.SideToMove == chess.Black {
		position.FullmoveNumber++
		position.SideToMove = chess.White
	} else {
		position.SideToMove = chess.Black
	}

	return position, nil
}
