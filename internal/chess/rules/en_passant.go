package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

type enPassantPlan struct {
	move           chess.Move
	capturedSquare chess.Square
}

func makeEnPassantPlan(position chess.Position, from chess.Square) (enPassantPlan, bool) {
	if !position.EnPassant.Valid() {
		return enPassantPlan{}, false
	}

	if position.SideToMove != chess.White && position.SideToMove != chess.Black {
		return enPassantPlan{}, false
	}

	pawn, ok := position.Board.PieceAt(from)
	if !ok || pawn.Type != chess.Pawn || pawn.Color != position.SideToMove {
		return enPassantPlan{}, false
	}

	captureRow := chess.A5.Row()
	if pawn.Color == chess.Black {
		captureRow = chess.A4.Row()
	}

	if from.Row() != captureRow {
		return enPassantPlan{}, false
	}

	to := position.EnPassant

	if !attacks.PawnAttacksSquare(from, to, pawn.Color) {
		return enPassantPlan{}, false
	}

	target, _ := position.Board.PieceAt(to)
	if target.Type != chess.NoPiece {
		return enPassantPlan{}, false
	}

	capturedSquare, ok := from.Offset(to.Column()-from.Column(), 0)
	if !ok {
		return enPassantPlan{}, false
	}

	enemyColor := chess.Black
	if pawn.Color == chess.Black {
		enemyColor = chess.White
	}

	captured, _ := position.Board.PieceAt(capturedSquare)
	if captured.Type != chess.Pawn || captured.Color != enemyColor {
		return enPassantPlan{}, false
	}

	return enPassantPlan{
		move:           chess.Move{From: from, To: to},
		capturedSquare: capturedSquare,
	}, true
}

func boardAfterEnPassant(board chess.Board, plan enPassantPlan) (chess.Board, error) {
	if _, err := board.RemovePiece(plan.capturedSquare); err != nil {
		return chess.Board{}, err
	}

	return BoardAfterMove(board, plan.move)
}

func EnPassantMove(position chess.Position, from chess.Square) (chess.Move, bool) {
	plan, ok := makeEnPassantPlan(position, from)
	if !ok {
		return chess.Move{}, false
	}

	board, err := boardAfterEnPassant(position.Board, plan)
	if err != nil {
		return chess.Move{}, false
	}

	inCheck, err := IsInCheck(board, position.SideToMove)
	if err != nil || inCheck {
		return chess.Move{}, false
	}

	return plan.move, true
}
