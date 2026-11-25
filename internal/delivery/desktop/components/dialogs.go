package components

import (
	"github.com/egoist/mygo/ui"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

var promotionOptions = [...]struct {
	piece chess.PieceType
	label string
}{
	{chess.Queen, "Queen"},
	{chess.Rook, "Rook"},
	{chess.Bishop, "Bishop"},
	{chess.Knight, "Knight"},
}

const (
	cancelAlert = 0
	acceptAlert = 1
)

func (a *App) Dialogs(c *ui.Context) {
	t := c.Theme()
	menuButton := func(label string, primary bool) bool {
		build := ui.Button
		if primary {
			build = ui.PrimaryButton
		}
		return build(c, label).Height(40).FillWidth().Clicked()
	}

	if ui.Modal(c, &a.paused, func() {
		ui.Column(c).Width(280).Gap(8).Children(func() {
			ui.Text(c, "Game menu").FontSize(20).Bold().Margin(0, 0, 8, 0)
			if menuButton("Resume", true) {
				a.resume()
			}
			if !a.vsAI && menuButton("Offer draw", false) {
				a.offerDraw()
			}
			if a.game.ClaimableDraw() != 0 && menuButton("Claim draw", false) {
				a.paused = false
				a.claimDraw()
			}
			if menuButton("Resign", false) {
				a.paused, a.confirmResign = false, true
			}
			if menuButton("Return to main menu", false) {
				a.paused, a.confirmLeave = false, true
			}
		})
	}).Dismissed() {
		a.resume()
	}

	switch ui.AlertDialog(c, &a.confirmResign, "Resign this game?", a.playerName(other(a.humanColor()))+" will win.", "Keep playing", "Resign") {
	case acceptAlert:
		a.resign()
	case cancelAlert, ui.AlertDismissed:
		a.StartAI()
	}
	switch ui.AlertDialog(c, &a.confirmLeave, "End game?", "The current game will be discarded.", "Keep playing", "End game") {
	case acceptAlert:
		a.toMenu()
	case cancelAlert, ui.AlertDismissed:
		a.StartAI()
	}
	s, _ := a.game.Snapshot()
	switch ui.AlertDialog(c, &a.drawOffered, "Draw offer", colorName(s.DrawOffer)+" offers a draw. Does "+colorName(other(s.DrawOffer))+" accept?", "Decline", "Accept draw") {
	case acceptAlert:
		a.answerDraw(true)
	case cancelAlert, ui.AlertDismissed:
		a.answerDraw(false)
	}

	ui.Modal(c, &a.choosing, func() {
		ui.Text(c, "Promote to").FontSize(18).Bold()
		ui.Row(c).Gap(10).Children(func() {
			for _, option := range promotionOptions {
				if ui.Button(c, pieceGlyphs[option.piece]).Label(option.label).FontSize(40).Size(76, 76).Clicked() {
					a.promote(option.piece)
				}
			}
		})
	})
	if !a.choosing {
		a.promoting = nil
	}

	if ui.Modal(c, &a.showResult, func() {
		title, why := a.result(s)
		ui.Column(c).Width(300).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "♚").FontSize(56).TextColor(t.Accent)
			ui.Text(c, title).FontSize(28).Bold()
			ui.Text(c, why).TextColor(t.TextMuted).Margin(0, 0, 12, 0)
			if menuButton("Play again", true) {
				a.StartGame(a.vsAI)
			}
			if menuButton("View board", false) {
				a.showResult, a.resultDismissed = false, true
			}
			if menuButton("Main menu", false) {
				a.toMenu()
			}
		})
	}).Dismissed() {
		a.resultDismissed = true
	}
}
