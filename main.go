package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func main() {
	//walls = BoxToVectors(1, 1, 10, 10)
	//walls = append(walls, BoxToVectors(5, 5, 1, 1)...)

	player = playerData{
		pos: XY64{X: 3, Y: 3}, angle: 4,
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
