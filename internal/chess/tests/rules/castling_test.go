package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

type castlingCase struct {
	name                    string
	color, enemy            chess.Color
	kingSide                bool
	right                   chess.CastlingRights
	wrongRights             []chess.CastlingRights
	kingMove, rookMove      chess.Move
	emptySquares, attackers []chess.Square
	canCastle               func(chess.Position) bool
	generate                func(chess.Position) (chess.Move, bool)
	apply                   func(chess.Position) (chess.Position, error)
}

var castlingCases = []castlingCase{
	{
		name: "white king-side", color: chess.White, enemy: chess.Black, kingSide: true,
		right: chess.WhiteKingSide, wrongRights: []chess.CastlingRights{chess.WhiteQueenSide},
		kingMove: chess.Move{From: chess.E1, To: chess.G1}, rookMove: chess.Move{From: chess.H1, To: chess.F1},
		emptySquares: []chess.Square{chess.F1, chess.G1}, attackers: []chess.Square{chess.E7, chess.F7, chess.G7},
		canCastle: rules.CanCastleKingSide, generate: rules.KingSideCastlingMove, apply: rules.ApplyKingSideCastling,
	},
	{
		name: "black king-side", color: chess.Black, enemy: chess.White, kingSide: true,
		right: chess.BlackKingSide, wrongRights: []chess.CastlingRights{chess.BlackQueenSide, chess.WhiteKingSide},
		kingMove: chess.Move{From: chess.E8, To: chess.G8}, rookMove: chess.Move{From: chess.H8, To: chess.F8},
		emptySquares: []chess.Square{chess.F8, chess.G8}, attackers: []chess.Square{chess.E2, chess.F2, chess.G2},
		canCastle: rules.CanCastleKingSide, generate: rules.KingSideCastlingMove, apply: rules.ApplyKingSideCastling,
	},
	{
		name: "white queen-side", color: chess.White, enemy: chess.Black,
		right:    chess.WhiteQueenSide,
		kingMove: chess.Move{From: chess.E1, To: chess.C1}, rookMove: chess.Move{From: chess.A1, To: chess.D1},
		emptySquares: []chess.Square{chess.B1, chess.C1, chess.D1}, attackers: []chess.Square{chess.E7, chess.D7, chess.C7},
		canCastle: rules.CanCastleQueenSide, generate: rules.QueenSideCastlingMove, apply: rules.ApplyQueenSideCastling,
	},
	{
		name: "black queen-side", color: chess.Black, enemy: chess.White,
		right:    chess.BlackQueenSide,
		kingMove: chess.Move{From: chess.E8, To: chess.C8}, rookMove: chess.Move{From: chess.A8, To: chess.D8},
		emptySquares: []chess.Square{chess.B8, chess.C8, chess.D8}, attackers: []chess.Square{chess.E2, chess.D2, chess.C2},
		canCastle: rules.CanCastleQueenSide, generate: rules.QueenSideCastlingMove, apply: rules.ApplyQueenSideCastling,
	},
}

func TestCanCastle(t *testing.T) {
	for _, castle := range castlingCases {
		t.Run(castle.name, func(t *testing.T) {
			type condition struct {
				name   string
				rights chess.CastlingRights
				square chess.Square
				piece  chess.Piece
				want   bool
			}
			cases := []condition{
				{"allowed", castle.right, chess.NoSquare, chess.Piece{}, true},
				{"without right", 0, chess.NoSquare, chess.Piece{}, false},
				{"missing king", castle.right, castle.kingMove.From, chess.Piece{}, false},
				{"missing rook", castle.right, castle.rookMove.From, chess.Piece{}, false},
			}
			for _, right := range castle.wrongRights {
				cases = append(cases, condition{"wrong right", right, chess.NoSquare, chess.Piece{}, false})
			}
			for _, square := range castle.emptySquares {
				cases = append(cases, condition{"blocked " + square.String(), castle.right, square, chess.Piece{Type: chess.Knight, Color: castle.color}, false})
			}
			for _, square := range castle.attackers {
				cases = append(cases, condition{"attacker " + square.String(), castle.right, square, chess.Piece{Type: chess.Queen, Color: castle.enemy}, false})
			}
			if !castle.kingSide {
				square := chess.A2
				if castle.color == chess.Black {
					square = chess.A7
				}
				cases = append(cases, condition{"b-square attacked", castle.right, square, chess.Piece{Type: chess.Bishop, Color: castle.enemy}, true})
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					position := castlingPosition(t, castle.color, castle.kingSide)
					position.CastlingRights = tc.rights
					if tc.square.Valid() {
						mustSetPiece(t, &position.Board, tc.square, tc.piece)
					}
					if got := castle.canCastle(position); got != tc.want {
						t.Errorf("Expected castling availability %v, got %v", tc.want, got)
					}
				})
			}
		})
	}
}

func TestCanCastleInvalidColor(t *testing.T) {
	for _, color := range []chess.Color{chess.NoColor, chess.Color(255)} {
		position := kingSidePosition(t, chess.White)
		position.SideToMove = color
		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected invalid color %v to prevent castling", color)
		}
	}
}

func TestCastlingMove(t *testing.T) {
	for _, castle := range castlingCases {
		t.Run(castle.name, func(t *testing.T) {
			position := castlingPosition(t, castle.color, castle.kingSide)
			move, ok := castle.generate(position)
			if !ok || move != castle.kingMove {
				t.Errorf("Expected move %v with ok=true, got %v with ok=%v", castle.kingMove, move, ok)
			}
			position.CastlingRights = 0
			if _, ok := castle.generate(position); ok {
				t.Error("Expected castling to be unavailable without rights")
			}
		})
	}
}

func TestApplyCastling(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide
	for _, castle := range castlingCases {
		t.Run(castle.name, func(t *testing.T) {
			position := castlingPosition(t, castle.color, castle.kingSide)
			position.CastlingRights = whiteRights | blackRights
			position.EnPassant = chess.C6
			position.HalfmoveClock = 7
			position.FullmoveNumber = 20
			before := position
			next, err := castle.apply(position)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			assertPositionPiece(t, next, castle.kingMove.From, chess.Piece{})
			assertPositionPiece(t, next, castle.rookMove.From, chess.Piece{})
			assertPositionPiece(t, next, castle.kingMove.To, chess.Piece{Type: chess.King, Color: castle.color})
			assertPositionPiece(t, next, castle.rookMove.To, chess.Piece{Type: chess.Rook, Color: castle.color})
			wantRights, wantFullmove := blackRights, uint16(20)
			if castle.color == chess.Black {
				wantRights, wantFullmove = whiteRights, 21
			}
			assertPositionState(t, next, positionState{castle.enemy, wantRights, chess.NoSquare, 8, wantFullmove})
			if position != before {
				t.Error("Expected original position to remain unchanged")
			}
			position.CastlingRights = 0
			before = position
			if _, err := castle.apply(position); !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Expected ErrIllegalMove, got %v", err)
			}
			if position != before {
				t.Error("Expected rejected castling to preserve original position")
			}
		})
	}
}
