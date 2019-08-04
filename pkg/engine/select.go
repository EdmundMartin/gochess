package engine

import (
	"fmt"
	"sort"

	"github.com/notnil/chess"
)

type chessLine struct {
	move  *chess.Move
	score int
}

type chessLines []chessLine

func (p chessLines) Len() int           { return len(p) }
func (p chessLines) Less(i, j int) bool { return p[i].score < p[j].score }
func (p chessLines) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func trimMoves(gameState *chess.Game, strat *Strategy) chessLines {
	moves := gameState.ValidMoves()
	lines := chessLines{}
	for _, m := range moves {
		clone := gameState.Clone()
		clone.Move(m)
		score := scorePosition(clone, strat)
		lines = append(lines, chessLine{move: m, score: score})
	}
	sort.Sort(sort.Reverse(lines))
	maxLines := min(len(lines), strat.Lines)
	return lines[:maxLines]
}

// SelectMove takes a Chess game and selects a move
func SelectMove(gameState *chess.Game, cpuColor chess.Color, strat *Strategy) *chess.Move {
	toAnalyse := trimMoves(gameState, strat)
	bestMove := searchTree(gameState, toAnalyse, strat)
	fmt.Println(bestMove)
	return bestMove
}

func quiesce(gameState *chess.Game, alpha int, beta int, strat *Strategy) int {
	standPat := scorePosition(gameState, strat)
	if standPat >= beta {
		return beta
	}
	if alpha < standPat {
		alpha = standPat
	}
	validMoves := trimMoves(gameState, strat)
	for _, m := range validMoves {
		if m.move.HasTag(chess.Capture) {
			clone := gameState.Clone()
			clone.Move(m.move)
			score := quiesce(clone, -beta, -alpha, strat)
			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}
	return alpha
}

func alphaBeta(gameState *chess.Game, alpha int, beta int, depthLeft int, strat *Strategy) int {
	bestScore := -9999
	if depthLeft == 0 {
		return quiesce(gameState, alpha, beta, strat)
	}
	moves := trimMoves(gameState, strat)
	for _, m := range moves {
		newState := gameState.Clone()
		newState.Move(m.move)
		score := -alphaBeta(newState, -beta, -alpha, depthLeft-1, strat)
		if score >= beta {
			return score
		}
		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
		}
	}
	return bestScore
}

func searchTree(gameState *chess.Game, lines chessLines, strat *Strategy) *chess.Move {
	var bestMove *chess.Move
	bestValue := -99999
	alpha := -100000
	beta := 100000
	for _, i := range lines {
		clone := gameState.Clone()
		clone.Move(i.move)
		boardValue := alphaBeta(clone, -beta, -alpha, strat.Depth, strat)
		if boardValue > bestValue {
			bestValue = boardValue
			bestMove = i.move
		}
		if boardValue > alpha {
			alpha = boardValue
		}
	}
	return bestMove
}

func calculateMaterial(posMap map[chess.Square]chess.Piece, strat *Strategy) int {
	score := 0
	for _, v := range posMap {
		pType := v.Type()
		if v.Color() == chess.White {
			score += strat.PieceValues[pType]
		} else {
			score -= strat.PieceValues[pType]
		}
	}
	return score
}

func calculatePosScore(posMap map[chess.Square]chess.Piece, strat *Strategy) int {
	score := 0
	for k, v := range posMap {
		pType := v.Type()
		pTable := strat.PieceTables[pType]
		val := pTable[int(k)]
		if v.Color() == chess.White {
			score += val
		} else {
			score -= val
		}
	}
	return score
}

func scorePosition(gameState *chess.Game, strat *Strategy) int {
	pos := gameState.Position()
	posHash := pos.Hash()
	posMap := pos.Board().SquareMap()
	cached, found := strat.Cache.Get(posHash)
	if found {
		return cached.(int)
	}
	score := calculateMaterial(posMap, strat)
	score += calculatePosScore(posMap, strat)
	strat.Cache.Add(posHash, score)
	return score
}
