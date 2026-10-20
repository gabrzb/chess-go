package game

import (
	"slices"
	"strconv"
	"strings"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
)

func SAN(before, next Snapshot, move chess.Move) string {
	piece, _ := before.Board.PieceAt(move.From)
	captured, _ := before.Board.PieceAt(move.To)
	letters := map[chess.PieceType]string{chess.Knight: "N", chess.Bishop: "B", chess.Rook: "R", chess.Queen: "Q", chess.King: "K"}
	var text strings.Builder
	if piece.Type == chess.King && abs(move.To.Column()-move.From.Column()) == 2 {
		text.WriteString("O-O")
		if move.To.Column() < move.From.Column() {
			text.WriteString("-O")
		}
	} else {
		capture := captured.Type != chess.NoPiece || (piece.Type == chess.Pawn && move.From.Column() != move.To.Column())
		text.WriteString(letters[piece.Type])
		if piece.Type != chess.Pawn {
			ambiguous, sameFile, sameRank := false, false, false
			for square := chess.A1; square <= chess.H8; square++ {
				other, _ := before.Board.PieceAt(square)
				if square == move.From || other != piece {
					continue
				}
				moves, err := rules.LegalMoves(before.Board, square)
				if err != nil || !slices.Contains(moves, chess.Move{From: square, To: move.To}) {
					continue
				}
				ambiguous = true
				sameFile = sameFile || square.Column() == move.From.Column()
				sameRank = sameRank || square.Row() == move.From.Row()
			}
			if ambiguous {
				if !sameFile || sameRank {
					text.WriteByte(move.From.String()[0])
				}
				if sameFile {
					text.WriteString(strconv.Itoa(move.From.Row() + 1))
				}
			}
		} else if capture {
			text.WriteByte(move.From.String()[0])
		}
		if capture {
			text.WriteByte('x')
		}
		text.WriteString(move.To.String())
		if move.Promotion != chess.NoPiece {
			text.WriteString("=" + letters[move.Promotion])
		}
	}
	if next.Status == chess.StatusCheckmate {
		text.WriteByte('#')
	} else if next.Check {
		text.WriteByte('+')
	}
	return text.String()
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
