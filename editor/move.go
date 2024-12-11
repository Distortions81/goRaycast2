package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Update is called every frame
func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	mpos := pos32{X: float32(mouseX), Y: float32(mouseY)}
	wpos := pos32{X: mpos.X - g.camera.X, Y: mpos.Y - g.camera.Y}

	//Follow cursor while placing player start
	if g.pStartMode {
		pStartPos = wpos
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) && !g.createMode {
		g.createMode = true
		g.firstClick = false
		g.secondClick = false
		g.start = pos32{X: 0, Y: 0}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		handlePMode(g)
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.createMode {
			snappedPos := snapPos(wpos, walls, lineSnapDist)
			if !g.secondClick && !g.firstClick {

				// Start creating a new vector
				g.start = snappedPos
				g.firstClick = true
			} else if g.firstClick && !g.secondClick {
				// Finish creating the vector
				endX := snappedPos.X
				endY := snappedPos.Y
				walls = append(walls, line32{
					X1: g.start.X,
					Y1: g.start.Y,
					X2: endX,
					Y2: endY,
				})
				fmt.Printf("created: %v,%v - %v,%v\n", g.start.X, g.start.Y, endX, endY)
				g.writeLevel()
				g.secondClick = true
				g.createMode = false
			}
		} else if g.pStartMode {
			handlePMode(g)
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.camera.X += float32(mouseX - int(g.lastMouse.X))
		g.camera.Y += float32(mouseY - int(g.lastMouse.Y))
	}

	g.lastMouse = mpos
	return nil
}

func handlePMode(g *Game) {
	if g.pStartMode {
		g.writeLevel()
	}
	g.pStartMode = !g.pStartMode
}
