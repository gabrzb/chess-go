package domain

type CastlingRights uint8

const (
	WhiteKingSide CastlingRights = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

func (r CastlingRights) Has(right CastlingRights) bool {
	return right != 0 && r&right == right
}

func (r CastlingRights) Without(right CastlingRights) CastlingRights {
	return r &^ right
}
