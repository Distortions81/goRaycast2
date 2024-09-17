package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	snapDistance = 10
	lineWidth    = 2
	scaleDiv     = 1
	gridSpacing  = 10
	levelPath    = "../level1.txt"
	gridBright   = 16
)

var (
	walls     = []Vector2D{}
	gridColor = color.NRGBA{R: gridBright, G: gridBright, B: gridBright, A: 255}
)

func main() {
	// Create a new game instance
	game := &Game{}

	readVecs()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Start the Ebiten game loop
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
