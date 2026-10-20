package rules

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/attacks"
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
	home := chess.A1
	var kingRight, queenRight chess.CastlingRights
	var enemyColor chess.Color
	switch color {
	case chess.White:
		kingRight, queenRight = chess.WhiteKingSide, chess.WhiteQueenSide
		enemyColor = chess.Black
	case chess.Black:
		home = chess.A8
		kingRight, queenRight = chess.BlackKingSide, chess.BlackQueenSide
		enemyColor = chess.White
	default:
		return castlingPlan{}, false
	}
	// First-row templates let both colors share the same castling layout.
	onHomeRow := func(square chess.Square) chess.Square {
		mapped, _ := square.Offset(0, home.Row())
		return mapped
	}
	plan := castlingPlan{enemyColor: enemyColor, kingMove: chess.Move{From: onHomeRow(chess.E1)}}
	switch side {
	case kingSide:
		plan.right = kingRight
		plan.kingMove.To = onHomeRow(chess.G1)
		plan.rookMove = chess.Move{From: onHomeRow(chess.H1), To: onHomeRow(chess.F1)}
		plan.emptySquares = []chess.Square{onHomeRow(chess.F1), onHomeRow(chess.G1)}
		plan.kingPath = []chess.Square{onHomeRow(chess.E1), onHomeRow(chess.F1), onHomeRow(chess.G1)}
	case queenSide:
		plan.right = queenRight
		plan.kingMove.To = onHomeRow(chess.C1)
		plan.rookMove = chess.Move{From: onHomeRow(chess.A1), To: onHomeRow(chess.D1)}
		plan.emptySquares = []chess.Square{onHomeRow(chess.B1), onHomeRow(chess.C1), onHomeRow(chess.D1)}
		plan.kingPath = []chess.Square{onHomeRow(chess.E1), onHomeRow(chess.D1), onHomeRow(chess.C1)}
	default:
		return castlingPlan{}, false
	}
	return plan, true
}

func castlingRightsWithoutPiece(
	rights chess.CastlingRights,
	piece chess.Piece,
	square chess.Square,
) chess.CastlingRights {
	for _, side := range []castlingSide{kingSide, queenSide} {
		plan, ok := makeCastlingPlan(piece.Color, side)
		if !ok {
			continue
		}

		kingLeft := piece.Type == chess.King
		rookLeft := piece.Type == chess.Rook && plan.rookMove.From == square

		if kingLeft || rookLeft {
			rights = rights.Without(plan.right)
		}
	}

	return rights
}
