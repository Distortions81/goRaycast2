package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Update is called every frame
func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	mpos := XY{X: float64(mouseX), Y: float64(mouseY)}
	wpos := XY{X: mpos.X - g.camera.X, Y: mpos.Y - g.camera.Y}

	if ebiten.IsKeyPressed(ebiten.KeyC) && !g.createMode {
		g.createMode = true
		g.firstClick = false
		g.secondClick = false
		g.start = XY{X: 0, Y: 0}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.createMode {
			snappedPos := snapPos(wpos, walls, snapDistance)
			if !g.secondClick && !g.firstClick {

				// Start creating a new vector
				g.start = snappedPos
				g.firstClick = true
			} else if g.firstClick && !g.secondClick {
				// Finish creating the vector
				endX := snappedPos.X
				endY := snappedPos.Y
				walls = append(walls, Vector2D{
					X1: g.start.X,
					Y1: g.start.Y,
					X2: endX,
					Y2: endY,
				})
				fmt.Printf("created: %v,%v - %v,%v\n", g.start.X, g.start.Y, endX, endY)
				g.writeVecs()
				g.secondClick = true
				g.createMode = false
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.camera.X += float64(mouseX - int(g.lastMouse.X))
		g.camera.Y += float64(mouseY - int(g.lastMouse.Y))
	}

	g.lastMouse = mpos
	return nil
}
