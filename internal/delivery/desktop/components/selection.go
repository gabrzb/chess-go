package components

import (
	"slices"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

func (a *App) clickSquare(square chess.Square) {
	if a.review >= 0 || a.opponent.IsTurn() || a.opponent.Thinking() {
		return
	}
	var movesToSquare []chess.Move
	for _, m := range a.targets {
		if m.To == square {
			movesToSquare = append(movesToSquare, m)
		}
	}
	a.selected, a.targets = chess.NoSquare, nil
	switch {
	case len(movesToSquare) == 1:
		a.play(movesToSquare[0])
		return
	case len(movesToSquare) > 1:
		a.promoting, a.choosing = movesToSquare, true
		return
	}
	snapshot, _ := a.game.Snapshot()
	piece, _ := snapshot.Board.PieceAt(square)
	if !a.canMove(snapshot, piece) {
		return
	}
	if moves, err := a.game.LegalMovesFrom(square); err == nil && len(moves) > 0 {
		a.selected, a.targets = square, moves
	}
}

func (a *App) canMove(s game.Snapshot, piece chess.Piece) bool {
	return a.review < 0 && !a.opponent.IsTurn() && !a.opponent.Thinking() &&
		s.Status == chess.StatusOngoing && piece.Type != chess.NoPiece && piece.Color == s.SideToMove
}

func (a *App) isTarget(sq chess.Square) bool {
	return slices.ContainsFunc(a.targets, func(m chess.Move) bool { return m.To == sq })
}

func (a *App) grab(from chess.Square) {
	if a.selected != from {
		a.selected, a.targets = chess.NoSquare, nil
		a.clickSquare(from)
	}
}

func (a *App) dropPiece(from, to chess.Square) {
	a.grab(from)
	if a.isTarget(to) {
		a.clickSquare(to)
	}
}

func (a *App) promote(kind chess.PieceType) {
	moves := a.promoting
	a.promoting, a.choosing = nil, false
	for _, m := range moves {
		if m.Promotion == kind {
			a.play(m)
			return
		}
	}
}
