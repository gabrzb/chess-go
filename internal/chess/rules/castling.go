package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

type castlingPlan struct {
	right        chess.CastlingRights
	enemyColor   chess.Color
	kingMove     chess.Move
	rookMove     chess.Move
	emptySquares []chess.Square
	kingPath     []chess.Square
}

type castlingSide uint8

const (
	kingSide castlingSide = iota
	queenSide
)

func availableCastlingPlan(position chess.Position, side castlingSide) (castlingPlan, bool) {
	plan, ok := makeCastlingPlan(position.SideToMove, side)
	if !ok {
		return castlingPlan{}, false
	}

	if !position.CastlingRights.Has(plan.right) {
		return castlingPlan{}, false
	}

	king, _ := position.Board.PieceAt(plan.kingMove.From)
	if king.Type != chess.King || king.Color != position.SideToMove {
		return castlingPlan{}, false
	}

	rook, _ := position.Board.PieceAt(plan.rookMove.From)
	if rook.Type != chess.Rook || rook.Color != position.SideToMove {
		return castlingPlan{}, false
	}

	for _, square := range plan.emptySquares {
		piece, _ := position.Board.PieceAt(square)

		if piece.Type != chess.NoPiece {
			return castlingPlan{}, false
		}
	}

	for _, square := range plan.kingPath {
		if attacks.IsSquareAttacked(position.Board, square, plan.enemyColor) {
			return castlingPlan{}, false
		}
	}

	return plan, true
}

func CanCastleKingSide(position chess.Position) bool {
	_, ok := availableCastlingPlan(position, kingSide)
	return ok
}

func CanCastleQueenSide(position chess.Position) bool {
	_, ok := availableCastlingPlan(position, queenSide)
	return ok
}

func castlingMove(position chess.Position, side castlingSide) (chess.Move, bool) {
	plan, ok := availableCastlingPlan(position, side)
	if !ok {
		return chess.Move{}, false
	}

	return plan.kingMove, true
}

func KingSideCastlingMove(position chess.Position) (chess.Move, bool) {
	return castlingMove(position, kingSide)
}

func QueenSideCastlingMove(position chess.Position) (chess.Move, bool) {
	return castlingMove(position, queenSide)
}

func applyCastling(position chess.Position, side castlingSide) (chess.Position, error) {
	plan, ok := availableCastlingPlan(position, side)
	if !ok {
		return chess.Position{}, chess.ErrIllegalMove
	}

	for _, move := range []chess.Move{plan.kingMove, plan.rookMove} {
		board, err := BoardAfterMove(position.Board, move)
		if err != nil {
			return chess.Position{}, err
		}

		position.Board = board
	}

	rightsToRemove := chess.WhiteKingSide | chess.WhiteQueenSide

	if position.SideToMove == chess.Black {
		rightsToRemove = chess.BlackKingSide | chess.BlackQueenSide
		position.FullmoveNumber++
	}

	position.CastlingRights = position.CastlingRights.Without(rightsToRemove)
	position.EnPassant = chess.NoSquare
	position.HalfmoveClock++
	position.SideToMove = plan.enemyColor

	return position, nil
}

func ApplyKingSideCastling(position chess.Position) (chess.Position, error) {
	return applyCastling(position, kingSide)
}

func ApplyQueenSideCastling(position chess.Position) (chess.Position, error) {
	return applyCastling(position, queenSide)
}

func makeCastlingPlan(color chess.Color, side castlingSide) (castlingPlan, bool) {
	var base chess.Square
	var kingRight, queenRight chess.CastlingRights
	var enemyColor chess.Color

	switch color {
	case chess.White:
		base = 0
		kingRight = chess.WhiteKingSide
		queenRight = chess.WhiteQueenSide
		enemyColor = chess.Black
	case chess.Black:
		base = 56
		kingRight = chess.BlackKingSide
		queenRight = chess.BlackQueenSide
		enemyColor = chess.White
	default:
		return castlingPlan{}, false
	}

	plan := castlingPlan{
		enemyColor: enemyColor,
		kingMove:   chess.Move{From: base + 4},
	}

	switch side {
	case kingSide:
		plan.right = kingRight
		plan.kingMove.To = base + 6
		plan.rookMove = chess.Move{From: base + 7, To: base + 5}
		plan.emptySquares = []chess.Square{base + 5, base + 6}
		plan.kingPath = []chess.Square{base + 4, base + 5, base + 6}
	case queenSide:
		plan.right = queenRight
		plan.kingMove.To = base + 2
		plan.rookMove = chess.Move{From: base + 0, To: base + 3}
		plan.emptySquares = []chess.Square{base + 1, base + 2, base + 3}
		plan.kingPath = []chess.Square{base + 4, base + 3, base + 2}
	default:
		return castlingPlan{}, false
	}

	return plan, true
}
