package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	levelPath = "../level1.txt"

	lineSnapDist = 10
	gridSnapDist = 5

	scaleDiv   = 1
	lineWidth  = 2
	gridBright = 25
	gridSize   = 25
)

var (
	walls     = []line32{}
	gridColor = color.NRGBA{R: gridBright, G: gridBright, B: gridBright, A: 255}
)

func main() {
	// Create a new game instance
	game := &Game{}

	readLevel()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Start the Ebiten game loop
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
