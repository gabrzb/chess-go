package chess

type Status uint8

const (
	StatusOngoing Status = iota
	StatusCheckmate
	StatusStalemate
)
