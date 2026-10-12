package rules_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

type promotionCase struct {
	name          string
	color, enemy  chess.Color
	move          chess.Move
	capturedRight chess.CastlingRights
}

var (
	whitePromotion        = promotionCase{"white advance", chess.White, chess.Black, chess.Move{From: chess.A7, To: chess.A8}, 0}
	blackPromotion        = promotionCase{"black advance", chess.Black, chess.White, chess.Move{From: chess.A2, To: chess.A1}, 0}
	whiteCapturePromotion = promotionCase{"white captures a8", chess.White, chess.Black, chess.Move{From: chess.B7, To: chess.A8}, chess.BlackQueenSide}
	blackCapturePromotion = promotionCase{"black captures a1", chess.Black, chess.White, chess.Move{From: chess.B2, To: chess.A1}, chess.WhiteQueenSide}
	promotionCases        = []promotionCase{
		whitePromotion, blackPromotion, whiteCapturePromotion,
		{"white captures h8", chess.White, chess.Black, chess.Move{From: chess.G7, To: chess.H8}, chess.BlackKingSide},
		blackCapturePromotion,
		{"black captures h1", chess.Black, chess.White, chess.Move{From: chess.G2, To: chess.H1}, chess.WhiteKingSide},
	}
	promotionChoices = []struct {
		name       string
		pieceType  chess.PieceType
		givesCheck bool
	}{
		{"queen", chess.Queen, true},
		{"rook", chess.Rook, true},
		{"bishop", chess.Bishop, false},
		{"knight", chess.Knight, false},
	}
)

func promotionPosition(t *testing.T, tc promotionCase) chess.Position {
	t.Helper()
	position := kingSidePosition(t, tc.color)
	mustSetPiece(t, &position.Board, tc.move.From, chess.Piece{Type: chess.Pawn, Color: tc.color})
	if tc.capturedRight != 0 {
		mustSetPiece(t, &position.Board, tc.move.To, chess.Piece{Type: chess.Rook, Color: tc.enemy})
		position.CastlingRights |= tc.capturedRight
	}
	position.EnPassant = chess.C3
	position.HalfmoveClock = initialHalfmoveClock
	position.FullmoveNumber = initialFullmoveNumber
	return position
}

func TestPromotion(t *testing.T) {
	for _, tc := range promotionCases {
		t.Run(tc.name, func(t *testing.T) {
			position := promotionPosition(t, tc)
			before := position
			moves, err := rules.LegalMovesFrom(position, tc.move.From)
			if err != nil {
				t.Fatalf("Unexpected generation error: %v", err)
			}
			for _, choice := range promotionChoices {
				t.Run(choice.name, func(t *testing.T) {
					move := tc.move
					move.Promotion = choice.pieceType
					if !slices.Contains(moves, move) {
						t.Errorf("Expected promotion %v in %v", move, moves)
					}
					next, err := rules.PositionAfterMove(position, move)
					if err != nil {
						t.Fatalf("Unexpected application error: %v", err)
					}
					assertPositionPiece(t, next, move.From, chess.Piece{})
					assertPositionPiece(t, next, move.To, chess.Piece{Type: choice.pieceType, Color: tc.color})
					wantFullmove := initialFullmoveNumber
					if tc.color == chess.Black {
						wantFullmove++
					}
					assertPositionState(t, next, positionState{tc.enemy, before.CastlingRights.Without(tc.capturedRight), chess.NoSquare, 0, wantFullmove})
					inCheck, err := rules.IsInCheck(next.Board, tc.enemy)
					if err != nil || inCheck != choice.givesCheck {
						t.Errorf("Expected opponent check=%v, got %v, %v", choice.givesCheck, inCheck, err)
					}
					if position != before {
						t.Error("Expected generation and application to preserve original position")
					}
				})
			}
		})
	}
}

func TestPromotionRejected(t *testing.T) {
	sides := []struct {
		fixture  promotionCase
		ordinary chess.Move
	}{
		{whitePromotion, chess.Move{From: chess.A6, To: chess.A7}},
		{blackPromotion, chess.Move{From: chess.A3, To: chess.A2}},
	}
	for _, side := range sides {
		t.Run(side.fixture.name, func(t *testing.T) {
			cases := []struct {
				name          string
				pieceType     chess.PieceType
				beforeLastRow bool
			}{
				{"missing choice", chess.NoPiece, false},
				{"pawn choice", chess.Pawn, false},
				{"king choice", chess.King, false},
				{"unknown choice", chess.PieceType(255), false},
				{"promotion before last row", chess.Queen, true},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					fixture := side.fixture
					if tc.beforeLastRow {
						fixture.move = side.ordinary
					}
					position := promotionPosition(t, fixture)
					before := position
					move := fixture.move
					move.Promotion = tc.pieceType
					if _, err := rules.PositionAfterMove(position, move); !errors.Is(err, chess.ErrIllegalMove) {
						t.Errorf("Expected ErrIllegalMove, got %v", err)
					}
					if position != before {
						t.Error("Expected rejected promotion to preserve original position")
					}
				})
			}
		})
	}
}

func TestPromotionCannotExposeKing(t *testing.T) {
	cases := []struct {
		fixture                       promotionCase
		enemyKing, safeKing, attacker chess.Square
	}{
		{promotionCase{"white pinned pawn", chess.White, chess.Black, chess.Move{From: chess.E7, To: chess.D8}, 0}, chess.E8, chess.A8, chess.E8},
		{promotionCase{"black pinned pawn", chess.Black, chess.White, chess.Move{From: chess.E2, To: chess.D1}, 0}, chess.E1, chess.A1, chess.E1},
	}
	for _, tc := range cases {
		t.Run(tc.fixture.name, func(t *testing.T) {
			fixture := tc.fixture
			position := promotionPosition(t, fixture)
			mustSetPiece(t, &position.Board, tc.enemyKing, chess.Piece{})
			mustSetPiece(t, &position.Board, tc.safeKing, chess.Piece{Type: chess.King, Color: fixture.enemy})
			mustSetPiece(t, &position.Board, tc.attacker, chess.Piece{Type: chess.Rook, Color: fixture.enemy})
			mustSetPiece(t, &position.Board, fixture.move.To, chess.Piece{Type: chess.Knight, Color: fixture.enemy})
			if check, err := rules.IsInCheck(position.Board, fixture.color); err != nil || check {
				t.Fatalf("Expected pawn to block rook before move, got check=%v, error=%v", check, err)
			}
			if candidates := movegen.PawnMoves(position.Board, fixture.move.From); len(candidates) != len(promotionChoices) {
				t.Fatalf("Expected four pseudolegal capture promotions, got %v", candidates)
			}
			before := position
			moves, err := rules.LegalMovesFrom(position, fixture.move.From)
			if err != nil {
				t.Fatalf("Unexpected generation error: %v", err)
			}
			for _, choice := range promotionChoices {
				t.Run(choice.name, func(t *testing.T) {
					move := fixture.move
					move.Promotion = choice.pieceType
					if slices.Contains(moves, move) {
						t.Errorf("Expected promotion %v to be filtered because it exposes the king", move)
					}
					if _, err := rules.PositionAfterMove(position, move); !errors.Is(err, chess.ErrIllegalMove) {
						t.Errorf("Expected ErrIllegalMove, got %v", err)
					}
					if position != before {
						t.Error("Expected rejected promotion to preserve original position")
					}
				})
			}
		})
	}
}

func TestPromotionCannotCaptureKing(t *testing.T) {
	sides := []struct {
		fixture   promotionCase
		enemyKing chess.Square
	}{
		{whiteCapturePromotion, chess.E8},
		{blackCapturePromotion, chess.E1},
	}
	for _, side := range sides {
		t.Run(side.fixture.name, func(t *testing.T) {
			fixture := side.fixture
			position := promotionPosition(t, fixture)
			mustSetPiece(t, &position.Board, side.enemyKing, chess.Piece{})
			mustSetPiece(t, &position.Board, fixture.move.To, chess.Piece{Type: chess.King, Color: fixture.enemy})
			moves, err := rules.LegalMovesFrom(position, fixture.move.From)
			if err != nil {
				t.Fatalf("Unexpected generation error: %v", err)
			}
			for _, move := range moves {
				if move.To == fixture.move.To {
					t.Errorf("Expected king capture to be filtered, got %v", move)
				}
			}
			move := fixture.move
			move.Promotion = chess.Queen
			before := position
			if _, err := rules.PositionAfterMove(position, move); !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Expected ErrIllegalMove, got %v", err)
			}
			if position != before {
				t.Error("Expected rejected king capture to preserve original position")
			}
		})
	}
}
