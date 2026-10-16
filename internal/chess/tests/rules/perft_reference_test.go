package rules_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func perftPositionFromFEN(t *testing.T, fen string) chess.Position {
	t.Helper()

	fields := strings.Fields(fen)
	if len(fields) != 6 {
		t.Fatalf("Expected six FEN fields: %q", fen)
	}

	placement, turn, castling, enPassant, halfmove, fullmove := fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]
	ranks := strings.Split(placement, "/")

	if len(ranks) != chess.BoardSize {
		t.Fatalf("Expected %d FEN ranks: %q", chess.BoardSize, placement)
	}

	types := map[string]chess.PieceType{
		"p": chess.Pawn, "n": chess.Knight, "b": chess.Bishop,
		"r": chess.Rook, "q": chess.Queen, "k": chess.King,
	}
	var board chess.Board

	for rankIndex, rank := range ranks {
		col := 0

		for _, symbol := range rank {
			if symbol >= '1' && symbol <= '8' {
				col += int(symbol - '0')
				continue
			}

			pieceType, ok := types[strings.ToLower(string(symbol))]
			if !ok {
				t.Fatalf("Invalid FEN piece: %q", symbol)
			}

			color := chess.Black

			if symbol >= 'A' && symbol <= 'Z' {
				color = chess.White
			}

			square, ok := chess.A8.Offset(col, -rankIndex)
			if !ok {
				t.Fatalf("FEN rank exceeds board width: %q", rank)
			}

			mustSetPiece(t, &board, square, chess.Piece{Type: pieceType, Color: color})
			col++
		}
		if col != chess.BoardSize {
			t.Fatalf("Expected %d squares in FEN rank: %q", chess.BoardSize, rank)
		}
	}
	var side chess.Color

	switch turn {
	case "w":
		side = chess.White
	case "b":
		side = chess.Black
	default:
		t.Fatalf("Invalid FEN side to move: %q", turn)
	}

	position, err := chess.NewPosition(board, side)
	if err != nil {
		t.Fatal(err)
	}

	rights := map[rune]chess.CastlingRights{
		'K': chess.WhiteKingSide, 'Q': chess.WhiteQueenSide,
		'k': chess.BlackKingSide, 'q': chess.BlackQueenSide,
	}

	if castling != "-" {
		for _, symbol := range castling {
			right, ok := rights[symbol]
			if !ok {
				t.Fatalf("Invalid FEN castling right: %q", symbol)
			}

			position.CastlingRights |= right
		}
	}

	if enPassant != "-" {
		position.EnPassant, err = chess.ParseSquare(enPassant)
		if err != nil {
			t.Fatalf("Invalid FEN en passant square: %v", err)
		}
	}

	halfmoveClock, err := strconv.ParseUint(halfmove, 10, 16)
	if err != nil {
		t.Fatalf("Invalid FEN halfmove clock: %v", err)
	}

	fullmoveNumber, err := strconv.ParseUint(fullmove, 10, 16)
	if err != nil || fullmoveNumber == 0 {
		t.Fatalf("Invalid FEN fullmove number: %q", fullmove)
	}

	position.HalfmoveClock = uint16(halfmoveClock)
	position.FullmoveNumber = uint16(fullmoveNumber)

	return position
}

func TestPerftReferencePositions(t *testing.T) {
	// Reference positions and counts: https://www.chessprogramming.org/Perft_Results
	cases := []struct {
		name string
		fen  string
		want uint64
	}{
		{
			name: "kiwipete",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			want: 97862,
		},
		{
			name: "rook and pawn endgame",
			fen:  "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			want: 2812,
		},
		{
			name: "promotions and check evasions",
			fen:  "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			want: 9467,
		},
		{
			name: "promotions and check evasions mirrored",
			fen:  "r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
			want: 9467,
		},
		{
			name: "tactical promotions",
			fen:  "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			want: 62379,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := perftPositionFromFEN(t, tc.fen)

			assertPerft(t, position, 3, tc.want)
		})
	}
}
