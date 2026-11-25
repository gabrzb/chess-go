package components

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

const (
	boardCornerRadius = 6
	raisedPieceLayer  = 10
	aiSlideDuration   = 320 * time.Millisecond
)

var pieceGlyphs = map[chess.PieceType]string{
	chess.King: "♚", chess.Queen: "♛", chess.Rook: "♜",
	chess.Bishop: "♝", chess.Knight: "♞", chess.Pawn: "♟︎",
}

func (a *App) Board(c *ui.Context, size float32) {
	position := a.displayedPosition()
	const lastRank = chess.BoardSize - 1
	ui.Column(c).Radius(boardCornerRadius).Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.28)).Children(func() {
		for row := range chess.BoardSize {
			rank := ui.Row(c.Key(row))
			rank.Children(func() {
				for col := range chess.BoardSize {
					y, x := lastRank-row, col
					if a.Flipped {
						y, x = row, lastRank-col
					}
					square := chess.Square(y*chess.BoardSize + x)
					if a.drawSquare(c, position, square, size, col == 0, row == lastRank) {
						// Raise the whole row so the piece paints over other ranks.
						rank.ZIndex(raisedPieceLayer)
					}
				}
			})
		}
	}).ClipInset(0, 0, 0, 0, boardCornerRadius)
}

func (a *App) drawSquare(c *ui.Context, position game.Snapshot, sq chess.Square, size float32, rankEdge, fileEdge bool) bool {
	base, ink := lightSquare, darkSquare
	if (sq.Row()+sq.Column())%2 == 0 {
		base, ink = darkSquare, lightSquare
	}
	piece, _ := position.Board.PieceAt(sq)
	movable := a.canMove(position, piece)
	cell := ui.Box(c.Key(sq)).Size(size, size).Background(base).Cursor(ui.CursorPointer)
	if !movable {
		cell.Label(sq.String())
	}
	if cell.Clicked() {
		a.clickSquare(sq)
	}
	target := a.isTarget(sq)
	if target && sq == a.dragTarget {
		cell.InsetShadow(0, 0, 0, size/16, ui.RGBA(255, 255, 255, 0.65))
	}
	fill := func() ui.Element { return ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).PassThrough() }
	var raised bool
	cell.Children(func() {
		if position.LastMove != nil && (position.LastMove.From == sq || position.LastMove.To == sq) {
			fill().Background(lastMoveTint)
		}
		if sq == a.selected {
			fill().Background(selectedTint)
		}
		if position.Check && piece.Type == chess.King && piece.Color == position.SideToMove {
			fill().Radius(size).Background(checkTint).Blur(size / 8)
		}
		label := max(10, size/6)
		if rankEdge {
			ui.Text(c, sq.String()[1:]).FontSize(label).Bold().TextColor(ink).Absolute().Top(2).Left(4).PassThrough()
		}
		if fileEdge {
			ui.Text(c, sq.String()[:1]).FontSize(label).Bold().TextColor(ink).Absolute().Bottom(1).Right(4).PassThrough()
		}
		if piece.Type != chess.NoPiece {
			raised = a.drawPiece(c, sq, piece, size, movable)
		}
		switch {
		case target && piece.Type != chess.NoPiece:
			fill().Margin(size/24).Radius(size).Border(size/12, hintColor)
		case target:
			fill().Center().Children(func() { ui.Box(c).Size(size*0.3, size*0.3).Radius(size).Background(hintColor) })
		}
	})
	if raised {
		cell.ZIndex(raisedPieceLayer)
	}
	return raised
}

func (a *App) drawPiece(c *ui.Context, sq chess.Square, piece chess.Piece, size float32, movable bool) bool {
	text, shade := whiteInk, ui.RGBA(0, 0, 0, 0.75)
	if piece.Color == chess.Black {
		text, shade = blackInk, ui.RGBA(255, 255, 255, 0.35)
	}
	var holder ui.Element
	var raised bool
	switch {
	case movable:
		holder, raised = a.draggablePiece(c, sq, size)
	case a.sliding != nil && sq == a.sliding.To && a.review == livePosition:
		holder, raised = a.slidingPiece(c, size), true
	default:
		holder = ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).PassThrough().Center()
	}
	holder.Children(func() {
		ui.Text(c, pieceGlyphs[piece.Type]).FontSize(size*0.75).TextColor(text).TextShadow(0, 1, 2, shade)
	})
	return raised
}

func (a *App) draggablePiece(c *ui.Context, sq chess.Square, size float32) (ui.Element, bool) {
	piece := ui.Box(c.Key("piece")).Absolute().Top(0).Left(0).Right(0).Bottom(0).Center().
		Label(sq.String()).Cursor(ui.CursorGrab)
	offset := ui.Local(piece, "offset", func() ui.Point { return ui.Point{} })
	piece.Draggable(offset, ui.DraggableOptions{Animation: ui.Snappy})
	if piece.Clicked() {
		a.clickSquare(sq)
	}
	if _, _, dragging := piece.Dragged(); dragging || piece.Pressed() && *offset != (ui.Point{}) {
		a.grab(sq)
		a.dragTarget = a.squareAt(sq, *offset, size)
		piece.Cursor(ui.CursorGrabbing)
	}
	if _, _, ended := piece.DragEnded(); ended {
		to := a.squareAt(sq, *offset, size)
		*offset, a.dragTarget = ui.Point{}, chess.NoSquare
		a.dropPiece(sq, to)
	}
	return piece, *offset != (ui.Point{})
}

func (a *App) squareAt(sq chess.Square, offset ui.Point, size float32) chess.Square {
	cols, rows := int(math.Round(float64(offset.X/size))), -int(math.Round(float64(offset.Y/size)))
	if a.Flipped {
		cols, rows = -cols, -rows
	}
	to, ok := sq.Offset(cols, rows)
	if !ok {
		return chess.NoSquare
	}
	return to
}

func (a *App) slidingPiece(c *ui.Context, size float32) ui.Element {
	from, to := a.sliding.From, a.sliding.To
	dx, dy := float32(from.Column()-to.Column())*size, float32(to.Row()-from.Row())*size
	if a.Flipped {
		dx, dy = -dx, -dy
	}
	return ui.Box(c.Key(len(a.moveHistory))).Absolute().Top(0).Left(0).Right(0).Bottom(0).PassThrough().Center().
		Transition(ui.ElementTransition{Enter: &ui.Motion{X: dx, Y: dy, Opacity: 1}, Duration: aiSlideDuration, Ease: ui.EaseInOut})
}
