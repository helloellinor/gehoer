package main

import (
	"gehoer/game"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1200, 800, "Gehør")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	g := game.New()
	g.Run()
}
