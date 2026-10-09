package movegen

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

// SlidingMoves generates moves along one unit direction until a piece blocks the path.
func SlidingMoves(board chess.Board, from chess.Square, dir chess.Direction) []chess.Move {
	if !from.Valid() {
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8
	var moves []chess.Move
	piece, _ := board.PieceAt(from)

	if piece.Type == chess.NoPiece {
		return nil
	}

	nextCol := col + dir.ColStep
	nextRow := row + dir.RowStep

	for nextCol >= 0 && nextCol < 8 && nextRow >= 0 && nextRow < 8 {
		destination := chess.Square(nextRow*8 + nextCol)
		occupant, _ := board.PieceAt(destination)

		if occupant.Type != chess.NoPiece && occupant.Color == piece.Color {
			break
		}

		moves = append(moves, chess.Move{From: from, To: destination})

		if occupant.Type != chess.NoPiece {
			break
		}

		nextCol += dir.ColStep
		nextRow += dir.RowStep
	}

	return moves
}

// CanOccupySquare checks whether a valid square is empty or holds an opposing piece.
func CanOccupySquare(board chess.Board, to chess.Square, color chess.Color) bool {
	occupant, ok := board.PieceAt(to)

	if !ok {
		return false
	}

	return occupant.Type == chess.NoPiece || occupant.Color != color
}

// JumpMoves generates a move to a single offset without checking intermediate squares.
func JumpMoves(board chess.Board, from chess.Square, dir chess.Direction) []chess.Move {
	piece, ok := board.PieceAt(from)

	if !ok || piece.Type == chess.NoPiece {
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8

	nextCol := col + dir.ColStep
	nextRow := row + dir.RowStep

	if nextCol < 0 || nextCol >= 8 || nextRow < 0 || nextRow >= 8 {
		return nil
	}

	destination := chess.Square(nextRow*8 + nextCol)

	if !CanOccupySquare(board, destination, piece.Color) {
		return nil
	}

	return []chess.Move{{From: from, To: destination}}
}

// PseudoLegalMoves selects a basic move generator before filtering king safety.
func PseudoLegalMoves(board chess.Board, from chess.Square) []chess.Move {
	piece, ok := board.PieceAt(from)

	if !ok {
		return nil
	}

	switch piece.Type {
	case chess.Pawn:
		return PawnMoves(board, from)
	case chess.Knight:
		return KnightMoves(board, from)
	case chess.Bishop:
		return BishopMoves(board, from)
	case chess.Rook:
		return RookMoves(board, from)
	case chess.Queen:
		return QueenMoves(board, from)
	case chess.King:
		return KingMoves(board, from)
	default:
		return nil
	}
}
