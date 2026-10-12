package rules_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

type enPassantCase struct {
	name           string
	color, enemy   chess.Color
	move           chess.Move
	capturedSquare chess.Square
}

var (
	whiteEnPassant = enPassantCase{"white to d6", chess.White, chess.Black,
		chess.Move{From: chess.E5, To: chess.D6}, chess.D5}
	blackEnPassant = enPassantCase{"black to d3", chess.Black, chess.White,
		chess.Move{From: chess.E4, To: chess.D3}, chess.D4}
	enPassantCases = []enPassantCase{
		whiteEnPassant,
		{"white to f6", chess.White, chess.Black,
			chess.Move{From: chess.E5, To: chess.F6}, chess.F5},
		blackEnPassant,
		{"black to f3", chess.Black, chess.White,
			chess.Move{From: chess.E4, To: chess.F3}, chess.F4},
	}
)

func enPassantPosition(t *testing.T, tc enPassantCase) chess.Position {
	t.Helper()
	position := kingSidePosition(t, tc.color)
	mustSetPiece(t, &position.Board, tc.move.From, chess.Piece{Type: chess.Pawn, Color: tc.color})
	mustSetPiece(t, &position.Board, tc.capturedSquare, chess.Piece{Type: chess.Pawn, Color: tc.enemy})
	position.EnPassant = tc.move.To
	position.HalfmoveClock = initialHalfmoveClock
	position.FullmoveNumber = initialFullmoveNumber
	return position
}

func TestEnPassant(t *testing.T) {
	for _, tc := range enPassantCases {
		t.Run(tc.name, func(t *testing.T) {
			position := enPassantPosition(t, tc)
			before := position
			move, ok := rules.EnPassantMove(position, tc.move.From)
			if !ok || move != tc.move {
				t.Fatalf("Expected en passant %v, got %v with ok=%v", tc.move, move, ok)
			}
			moves, err := rules.LegalMovesFrom(position, tc.move.From)
			if err != nil {
				t.Fatalf("Unexpected generation error: %v", err)
			}
			if !slices.Contains(moves, tc.move) {
				t.Errorf("Expected en passant %v in %v", tc.move, moves)
			}
			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected application error: %v", err)
			}
			assertPositionPiece(t, next, tc.move.From, chess.Piece{})
			assertPositionPiece(t, next, tc.capturedSquare, chess.Piece{})
			assertPositionPiece(t, next, tc.move.To, chess.Piece{Type: chess.Pawn, Color: tc.color})
			wantFullmove := initialFullmoveNumber
			if tc.color == chess.Black {
				wantFullmove++
			}
			assertPositionState(t, next, positionState{tc.enemy, before.CastlingRights, chess.NoSquare, 0, wantFullmove})
			if position != before {
				t.Error("Expected generation and application to preserve original position")
			}
		})
	}
}

func TestEnPassantUnavailable(t *testing.T) {
	cases := []struct {
		name          string
		shiftRow      bool
		invalidSource bool
		arrange       func(*testing.T, *chess.Position, enPassantCase)
	}{
		{name: "without opportunity", arrange: func(_ *testing.T, p *chess.Position, _ enPassantCase) {
			p.EnPassant = chess.NoSquare
		}},
		{name: "invalid recorded square", arrange: func(_ *testing.T, p *chess.Position, _ enPassantCase) {
			p.EnPassant = chess.Square(255)
		}},
		{name: "invalid turn", arrange: func(_ *testing.T, p *chess.Position, _ enPassantCase) {
			p.SideToMove = chess.NoColor
		}},
		{name: "opposing turn", arrange: func(_ *testing.T, p *chess.Position, tc enPassantCase) {
			p.SideToMove = tc.enemy
		}},
		{name: "invalid source", invalidSource: true},
		{name: "empty source", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.move.From, chess.Piece{})
		}},
		{name: "source is not pawn", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.move.From, chess.Piece{Type: chess.Knight, Color: tc.color})
		}},
		{name: "wrong capture row", shiftRow: true},
		{name: "destination is not diagonal", arrange: func(_ *testing.T, p *chess.Position, tc enPassantCase) {
			p.EnPassant, _ = tc.move.From.Offset(0, tc.move.To.Row()-tc.move.From.Row())
		}},
		{name: "occupied destination", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.move.To, chess.Piece{Type: chess.Knight, Color: tc.enemy})
		}},
		{name: "missing captured pawn", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.capturedSquare, chess.Piece{})
		}},
		{name: "captured piece is not pawn", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.capturedSquare, chess.Piece{Type: chess.Knight, Color: tc.enemy})
		}},
		{name: "captured pawn belongs to own side", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.capturedSquare, chess.Piece{Type: chess.Pawn, Color: tc.color})
		}},
		{name: "captured pawn has invalid color", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, tc.capturedSquare, chess.Piece{Type: chess.Pawn, Color: chess.NoColor})
		}},
		{name: "missing king", arrange: func(t *testing.T, p *chess.Position, tc enPassantCase) {
			mustSetPiece(t, &p.Board, rules.FindKing(p.Board, tc.color), chess.Piece{})
		}},
	}
	for _, base := range []enPassantCase{whiteEnPassant, blackEnPassant} {
		t.Run(base.name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					fixture := base
					if tc.shiftRow {
						fixture.move.From, _ = fixture.move.From.Offset(0, -1)
						fixture.move.To, _ = fixture.move.To.Offset(0, -1)
						fixture.capturedSquare, _ = fixture.capturedSquare.Offset(0, -1)
					}
					position := enPassantPosition(t, fixture)
					if tc.arrange != nil {
						tc.arrange(t, &position, fixture)
					}
					from := fixture.move.From
					if tc.invalidSource {
						from = chess.NoSquare
					}
					before := position
					if move, ok := rules.EnPassantMove(position, from); ok {
						t.Errorf("Expected en passant to be unavailable, got %v", move)
					}
					if position != before {
						t.Error("Expected unavailable capture to preserve original position")
					}
				})
			}
		})
	}
}

func TestEnPassantKingSafety(t *testing.T) {
	sides := []struct {
		fixture                                                    enPassantCase
		kingHome, enemyHome, enemySafe                             chess.Square
		discoveredKing, discoveredRook, checkedKing, unrelatedRook chess.Square
	}{
		{whiteEnPassant, chess.E1, chess.E8, chess.A8, chess.A5, chess.H5, chess.E4, chess.B1},
		{blackEnPassant, chess.E8, chess.E1, chess.A1, chess.A4, chess.H4, chess.E5, chess.B8},
	}
	for _, side := range sides {
		t.Run(side.fixture.name, func(t *testing.T) {
			cases := []struct {
				name                                     string
				king, rook                               chess.Square
				moveEnemyKing, initiallyInCheck, allowed bool
			}{
				{"opens horizontal rook attack", side.discoveredKing, side.discoveredRook, false, false, false},
				{"pinned pawn leaves king exposed", side.kingHome, side.enemyHome, true, false, false},
				{"captures checking pawn", side.checkedKing, chess.NoSquare, false, true, true},
				{"does not resolve existing check", side.kingHome, side.unrelatedRook, false, true, false},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					fixture := side.fixture
					position := enPassantPosition(t, fixture)
					mustSetPiece(t, &position.Board, side.kingHome, chess.Piece{})
					mustSetPiece(t, &position.Board, tc.king, chess.Piece{Type: chess.King, Color: fixture.color})
					if tc.moveEnemyKing {
						mustSetPiece(t, &position.Board, side.enemyHome, chess.Piece{})
						mustSetPiece(t, &position.Board, side.enemySafe, chess.Piece{Type: chess.King, Color: fixture.enemy})
					}
					if tc.rook.Valid() {
						mustSetPiece(t, &position.Board, tc.rook, chess.Piece{Type: chess.Rook, Color: fixture.enemy})
					}
					inCheck, err := rules.IsInCheck(position.Board, fixture.color)
					if err != nil || inCheck != tc.initiallyInCheck {
						t.Fatalf("Invalid setup: expected check=%v, got %v, %v", tc.initiallyInCheck, inCheck, err)
					}
					before := position
					move, ok := rules.EnPassantMove(position, fixture.move.From)
					if ok != tc.allowed || ok && move != fixture.move {
						t.Errorf("Expected availability %v and move %v, got %v, %v", tc.allowed, fixture.move, ok, move)
					}
					moves, err := rules.LegalMovesFrom(position, fixture.move.From)
					if err != nil {
						t.Fatalf("Unexpected generation error: %v", err)
					}
					if slices.Contains(moves, fixture.move) != tc.allowed {
						t.Errorf("Expected capture membership %v in %v", tc.allowed, moves)
					}
					next, err := rules.PositionAfterMove(position, fixture.move)
					if tc.allowed {
						if err != nil {
							t.Fatalf("Unexpected application error: %v", err)
						}
						if check, err := rules.IsInCheck(next.Board, fixture.color); err != nil || check {
							t.Errorf("Expected king to be safe, got check=%v, error=%v", check, err)
						}
					} else if !errors.Is(err, chess.ErrIllegalMove) {
						t.Errorf("Expected ErrIllegalMove, got %v", err)
					}
					if position != before {
						t.Error("Expected original position to remain unchanged")
					}
				})
			}
		})
	}
}

func TestEnPassantLifecycle(t *testing.T) {
	cases := []struct {
		fixture                       enPassantCase
		doubleAdvance, decline, reply chess.Move
	}{
		{whiteEnPassant, chess.Move{From: chess.D7, To: chess.D5}, chess.Move{From: chess.H1, To: chess.H2}, chess.Move{From: chess.E8, To: chess.F8}},
		{blackEnPassant, chess.Move{From: chess.D2, To: chess.D4}, chess.Move{From: chess.H8, To: chess.H7}, chess.Move{From: chess.E1, To: chess.F1}},
	}
	for _, tc := range cases {
		t.Run(tc.fixture.name, func(t *testing.T) {
			position := enPassantPosition(t, tc.fixture)
			victim := chess.Piece{Type: chess.Pawn, Color: tc.fixture.enemy}
			mustSetPiece(t, &position.Board, tc.fixture.capturedSquare, chess.Piece{})
			mustSetPiece(t, &position.Board, tc.doubleAdvance.From, victim)
			position.SideToMove = tc.fixture.enemy
			position.EnPassant = chess.NoSquare
			next, err := rules.PositionAfterMove(position, tc.doubleAdvance)
			if err != nil {
				t.Fatalf("Unexpected double-advance error: %v", err)
			}
			if next.EnPassant != tc.fixture.move.To {
				t.Errorf("Expected crossed square %v, got %v", tc.fixture.move.To, next.EnPassant)
			}
			if move, ok := rules.EnPassantMove(next, tc.fixture.move.From); !ok || move != tc.fixture.move {
				t.Fatalf("Expected capture after double advance, got %v, %v", move, ok)
			}
			for _, move := range []chess.Move{tc.decline, tc.reply} {
				next, err = rules.PositionAfterMove(next, move)
				if err != nil {
					t.Fatalf("Unexpected continuation error for %v: %v", move, err)
				}
				if next.EnPassant != chess.NoSquare {
					t.Errorf("Expected opportunity to expire after %v, got %v", move, next.EnPassant)
				}
			}
			assertPositionPiece(t, next, tc.fixture.move.From, chess.Piece{Type: chess.Pawn, Color: tc.fixture.color})
			assertPositionPiece(t, next, tc.fixture.capturedSquare, victim)
			if _, ok := rules.EnPassantMove(next, tc.fixture.move.From); ok {
				t.Error("Expected expired capture to be unavailable")
			}
			before := next
			if _, err := rules.PositionAfterMove(next, tc.fixture.move); !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Expected expired capture to return ErrIllegalMove, got %v", err)
			}
			if next != before {
				t.Error("Expected rejected capture to preserve original position")
			}
		})
	}
}

func TestNormalPawnMoveWithEnPassantAvailable(t *testing.T) {
	cases := []struct {
		fixture enPassantCase
		move    chess.Move
	}{
		{whiteEnPassant, chess.Move{From: chess.E5, To: chess.E6}},
		{blackEnPassant, chess.Move{From: chess.E4, To: chess.E3}},
	}
	for _, tc := range cases {
		t.Run(tc.fixture.name, func(t *testing.T) {
			position := enPassantPosition(t, tc.fixture)
			next, err := rules.PositionAfterMove(position, tc.move)
			if err != nil {
				t.Fatalf("Unexpected normal-move error: %v", err)
			}
			assertPositionPiece(t, next, tc.move.From, chess.Piece{})
			assertPositionPiece(t, next, tc.move.To, chess.Piece{Type: chess.Pawn, Color: tc.fixture.color})
			assertPositionPiece(t, next, tc.fixture.capturedSquare, chess.Piece{Type: chess.Pawn, Color: tc.fixture.enemy})
			if next.EnPassant != chess.NoSquare {
				t.Error("Expected normal move to expire en passant opportunity")
			}
		})
	}
}
