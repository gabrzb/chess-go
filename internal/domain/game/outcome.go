package game

import (
	"errors"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

type EndReason uint8

const (
	NoReason EndReason = iota
	Checkmate
	Stalemate
	Resignation
	Agreement
	InsufficientMaterial
	Repetition
	FiftyMoves
	FivefoldRepetition
	SeventyFiveMoves
)

type Outcome struct {
	Winner chess.Color
	Reason EndReason
}

var ErrDrawOffer = errors.New("invalid or pending draw offer")

func other(color chess.Color) chess.Color {
	if color == chess.White {
		return chess.Black
	}
	return chess.White
}

func outcomeFor(status chess.Status, side chess.Color) Outcome {
	switch status {
	case chess.StatusCheckmate:
		return Outcome{Winner: other(side), Reason: Checkmate}
	case chess.StatusStalemate:
		return Outcome{Reason: Stalemate}
	case chess.StatusInsufficientMaterial:
		return Outcome{Reason: InsufficientMaterial}
	case chess.StatusRepetition:
		return Outcome{Reason: Repetition}
	case chess.StatusFiftyMoves:
		return Outcome{Reason: FiftyMoves}
	case chess.StatusFivefoldRepetition:
		return Outcome{Reason: FivefoldRepetition}
	case chess.StatusSeventyFiveMoves:
		return Outcome{Reason: SeventyFiveMoves}
	default:
		return Outcome{}
	}
}

func (g *Game) Resign(color chess.Color) error {
	if g.status != chess.StatusOngoing {
		return ErrGameOver
	}
	if color != chess.White && color != chess.Black {
		return chess.ErrInvalidColor
	}
	g.status = chess.StatusResigned
	g.outcome = Outcome{Winner: other(color), Reason: Resignation}
	g.drawOffer = chess.NoColor
	return nil
}

func (g *Game) OfferDraw(color chess.Color) error {
	if g.status != chess.StatusOngoing {
		return ErrGameOver
	}
	if color != g.position.SideToMove || g.drawOffer != chess.NoColor {
		return ErrDrawOffer
	}
	g.drawOffer = color
	return nil
}

func (g *Game) RespondDraw(color chess.Color, accept bool) error {
	if g.status != chess.StatusOngoing {
		return ErrGameOver
	}
	if g.drawOffer == chess.NoColor || color != other(g.drawOffer) {
		return ErrDrawOffer
	}
	g.drawOffer = chess.NoColor
	if accept {
		g.status = chess.StatusDrawAgreed
		g.outcome = Outcome{Reason: Agreement}
	}
	return nil
}
