package chess

import "testing"

func checkMoveString(t *testing.T, move Move, expected string) {
	got := move.String()

	if got != expected {
		t.Errorf("Move.String() = %s; expected %s", got, expected)
	}
}

func checkMovePromotionString(t *testing.T, move Move, expected string) {
	got := move.String()

	if got != expected {
		t.Errorf("Move.String() with promotion = %s; expected %s", got, expected)
	}
}

func TestMoveString(t *testing.T) {
	checkMoveString(
		t,
		Move{From: Square(12), To: Square(28)},
		"e2e4",
	)
	checkMovePromotionString(
		t,
		Move{
			From: Square(52), To: Square(60), Promotion: Knight,
		},
		"e7e8n",
	)
	checkMovePromotionString(
		t,
		Move{
			From: Square(52), To: Square(60), Promotion: Bishop,
		},
		"e7e8b",
	)
	checkMovePromotionString(
		t,
		Move{
			From: Square(52), To: Square(60), Promotion: Rook,
		},
		"e7e8r",
	)
	checkMovePromotionString(
		t,
		Move{
			From: Square(52), To: Square(60), Promotion: Queen,
		},
		"e7e8q",
	)
}
