package engine

import (
	"github.com/hashicorp/golang-lru"
	"github.com/notnil/chess"
)

// Strategy contains postion tables and piece values
type Strategy struct {
	PieceValues map[chess.PieceType]int
	PieceTables map[chess.PieceType][64]int
	Lines       int
	Depth       int
	Cache       *lru.Cache
}

// DefaultStrategy returns default piece values and position tables
func DefaultStrategy(cache *lru.Cache) *Strategy {
	PieceValues := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Knight: 280,
		chess.Bishop: 320,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0,
	}
	TableValues := map[chess.PieceType][64]int{
		chess.Pawn:   defaultPawnTable,
		chess.Knight: defaultKnightTable,
		chess.Bishop: defaultBishopTable,
		chess.Rook:   defaultRookTable,
		chess.Queen:  defaultQueenTable,
		chess.King:   defaultKingTable,
	}
	strat := &Strategy{PieceValues: PieceValues, PieceTables: TableValues, Lines: 5, Depth: 5, Cache: cache}
	return strat
}

// DefaultBlackStrategy reverses boards for black
func DefaultBlackStrategy(cache *lru.Cache) *Strategy {
	PieceValues := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Knight: 280,
		chess.Bishop: 320,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0,
	}
	TableValues := map[chess.PieceType][64]int{
		chess.Pawn:   reverseArray(defaultPawnTable),
		chess.Knight: reverseArray(defaultKnightTable),
		chess.Bishop: reverseArray(defaultBishopTable),
		chess.Rook:   reverseArray(defaultRookTable),
		chess.Queen:  reverseArray(defaultQueenTable),
		chess.King:   reverseArray(defaultKingTable),
	}
	strat := &Strategy{PieceValues: PieceValues, PieceTables: TableValues, Lines: 5, Depth: 5, Cache: cache}
	return strat
}

func reverseArray(boardArray [64]int) [64]int {
	var revBoard [64]int
	count := 0
	for i := 63; i >= 0; i-- {
		revBoard[count] = i
		count++
	}
	return revBoard
}

// SelectDefaultStrategy returns strategy based on color
func SelectDefaultStrategy(cache *lru.Cache, color chess.Color) *Strategy {
	if color == chess.White {
		return DefaultStrategy(cache)
	}
	return DefaultBlackStrategy(cache)
}
