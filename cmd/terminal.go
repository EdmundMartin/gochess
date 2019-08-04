package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EdmundMartin/gochess/pkg/engine"
	"github.com/hashicorp/golang-lru"
	"github.com/notnil/chess"
)

func selectColor(reader *bufio.Reader) chess.Color {
	fmt.Print("CPU Color: ")
	text, _ := reader.ReadString('\n')
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "white" {
		return chess.White
	}
	return chess.Black
}

func playersTurn(reader *bufio.Reader, game *chess.Game) {
	fmt.Println(game.Position().Board().Draw())
	fmt.Print("Enter move: ")
	text, _ := reader.ReadString('\n')
	t := strings.TrimSpace(text)
	err := game.MoveStr(t)
	if err != nil {
		fmt.Print("Invalid move")

	}
}

func checkWinner(CPU chess.Color, outcome chess.Outcome) string {
	if CPU == chess.White && outcome == "1-0" {
		return "CPU Wins"
	}
	if CPU == chess.Black && outcome == "0-1" {
		return "CPU Wins"
	}
	if outcome == "1/2-1/2" {
		return "Draw"
	}
	return "Player Wins"
}

func main() {
	l, _ := lru.New(200000)
	game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))
	reader := bufio.NewReader(os.Stdin)
	CPU := selectColor(reader)
	strat := engine.SelectDefaultStrategy(l, CPU)
	for game.Outcome() == chess.NoOutcome {
		if game.Position().Turn() == CPU {
			start := time.Now()
			move := engine.SelectMove(game, chess.White, strat)
			game.Move(move)
			end := time.Now()
			fmt.Printf("CPU took %f seconds", end.Sub(start).Seconds())
		} else {
			playersTurn(reader, game)
		}
	}
	result := checkWinner(CPU, game.Outcome())
	fmt.Println(result)
}
