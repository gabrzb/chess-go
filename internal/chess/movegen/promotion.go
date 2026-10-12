package movegen

import "github.com/gabrzb/chess-go/internal/chess"

func pawnMovesTo(from, to chess.Square) []chess.Move {
	row := to.Row()
	if row != chess.A1.Row() && row != chess.A8.Row() {
		return []chess.Move{{From: from, To: to}}
	}

	var moves []chess.Move
	for _, promotion := range []chess.PieceType{
		chess.Queen,
		chess.Rook,
		chess.Bishop,
		chess.Knight,
	} {
		moves = append(moves, chess.Move{
			From:      from,
			To:        to,
			Promotion: promotion,
		})
	}

	return moves
}
