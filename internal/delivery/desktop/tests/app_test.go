package desktop_test

import (
	"context"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/engine"
)

func newTester(t *testing.T, options desktop.Options) *ui.Tester {
	t.Helper()
	app, err := desktop.New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	return ui.NewTester(app.View, 1100, 780)
}

func click(t *testing.T, tt *ui.Tester, labels ...string) {
	t.Helper()
	for _, label := range labels {
		if err := tt.Click(label); err != nil {
			t.Fatalf("click %s: %v (texts %q)", label, err, tt.Texts())
		}
	}
}

func expect(t *testing.T, tt *ui.Tester, texts ...string) {
	t.Helper()
	for _, want := range texts {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
}

func TestMenuLeadsToSetupAndBack(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	expect(t, tt, "CHESS GO", "Player vs Player", "Player vs AI", "Quit")
	click(t, tt, "Player vs AI")
	expect(t, tt, "New game", "Random", "Medium")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	expect(t, tt, "CHESS GO")
}

func TestLocalGameMovesAndHistory(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player")
	expect(t, tt, "Player 1", "Player 2", "White to move", "No moves yet")
	click(t, tt, "e2", "e4", "e7", "e5", "g1", "f3")
	expect(t, tt, "Nf3", "2.", "Black to move")
	click(t, tt, "Previous move")
	expect(t, tt, "Reviewing an earlier position — press → to return")
	click(t, tt, "Next move")
	if tt.HasText("Reviewing an earlier position — press → to return") {
		t.Errorf("still reviewing: %q", tt.Texts())
	}
}

func TestIllegalClickSelectsNothing(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player", "e2", "e5")
	expect(t, tt, "No moves yet", "White to move")
}

func TestFoolsMateShowsResult(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player", "f2", "f3", "e7", "e5", "g2", "g4", "d8", "h4")
	expect(t, tt, "Black wins", "Checkmate.", "Play again", "View board")
	click(t, tt, "View board")
	expect(t, tt, "Main menu", "Play again")
	click(t, tt, "Play again")
	expect(t, tt, "White to move", "No moves yet")
}

func TestResignFromGameMenu(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player", "☰  Menu", "Resign", "Resign")
	expect(t, tt, "Black wins", "White resigned.")
	click(t, tt, "Main menu")
	expect(t, tt, "CHESS GO")
}

func TestLocalDrawOffer(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player", "☰  Menu", "Offer draw", "Accept draw")
	expect(t, tt, "Draw", "Draw agreed.")
}

func TestAIRepliesThroughUpdate(t *testing.T) {
	updates := make(chan func(), 1)
	tt := newTester(t, desktop.Options{
		Difficulty: engine.Easy, HumanColor: chess.White, Seed: 1,
		Update: func(fn func()) { updates <- fn },
	})
	expect(t, tt, "You", "AI · Easy")
	click(t, tt, "e2", "e4")
	select {
	case fn := <-updates:
		fn()
	case <-time.After(10 * time.Second):
		t.Fatal("AI never replied")
	}
	tt.Frame()
	expect(t, tt, "White to move")
}

func TestInvalidDifficulty(t *testing.T) {
	if _, err := desktop.New(context.Background(), desktop.Options{Difficulty: "impossible"}); err == nil {
		t.Error("expected an error")
	}
}

func drag(t *testing.T, tt *ui.Tester, from, to string) {
	t.Helper()
	a, okA := tt.Find(from)
	b, okB := tt.Find(to)
	if !okA || !okB {
		t.Fatalf("squares %s/%s not found", from, to)
	}
	ax, ay := a.X+a.W/2, a.Y+a.H/2
	bx, by := b.X+b.W/2, b.Y+b.H/2
	tt.Press(ax, ay)
	for i := float32(1); i <= 10; i++ {
		tt.Move(ax+(bx-ax)*i/10, ay+(by-ay)*i/10)
	}
	tt.Release(bx, by)
	tt.Frame()
}

func TestDragPlaysLegalMove(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player")
	drag(t, tt, "g1", "f3")
	expect(t, tt, "Nf3", "Black to move")
}

func TestDragToIllegalSquarePutsPieceBack(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player")
	drag(t, tt, "e2", "e5")
	expect(t, tt, "No moves yet", "White to move")
	drag(t, tt, "e7", "e5")
	expect(t, tt, "No moves yet", "White to move")
}

func TestDragPlaysLegalMoveOnFlippedBoard(t *testing.T) {
	tt := newTester(t, desktop.Options{})
	click(t, tt, "Player vs Player", "Flip board")
	drag(t, tt, "e2", "e4")
	expect(t, tt, "e4", "Black to move")
}
