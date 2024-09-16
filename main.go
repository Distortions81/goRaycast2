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

	game := &Game{
		player: Player{
			pos: XY64{X: 3, Y: 3}, dir: XY64{X: -1, Y: -0.25}, plane: XY64{X: -0.15, Y: 0.65},
		},
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
