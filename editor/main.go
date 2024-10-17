package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

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
	bgImage   *ebiten.Image
)

func main() {
	// Create a new game instance
	game := &Game{}

	readLevel()
	loadImg()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Start the Ebiten game loop
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func loadImg() {
	file, err := os.OpenFile("trace.png", os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Unable to open trace.png")
		return
	}
	defer file.Close()

	src, err := png.Decode(file)
	if err != nil {
		fmt.Println("Unable to decode trace.png")
		return
	}
	bgImage = ebiten.NewImageFromImage(src)
}
