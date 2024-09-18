package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

// Draw is called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a white background
	screen.Fill(color.Black)

	drawGrid(g, screen)

	// Create a new vector image to draw on

	// Draw each vector with respect to the camera position
	for _, vec := range walls {
		x1 := vec.X1 + g.camera.X
		y1 := vec.Y1 + g.camera.Y
		x2 := vec.X2 + g.camera.X
		y2 := vec.Y2 + g.camera.Y
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), lineWidth, color.White, true)
	}

	if g.createMode {
		mouseX, mouseY := ebiten.CursorPosition()
		mpos := XY{X: float64(mouseX), Y: float64(mouseY)}
		snappedPos := snapPos(mpos, walls, lineSnapDist)

		if g.createMode {
			vector.DrawFilledCircle(screen, float32(snappedPos.X), float32(snappedPos.Y), lineWidth*2, colornames.Yellow, true)
		}
		if g.firstClick && !g.secondClick {
			vector.StrokeLine(screen, float32(g.start.X), float32(g.start.Y), float32(snappedPos.X), float32(snappedPos.Y), lineWidth, colornames.Red, true)
		}
	}

	// Draw text for clarity
	if g.createMode {
		ebitenutil.DebugPrint(screen, "Vector created, click again to specify vector end.")
	} else {
		ebitenutil.DebugPrint(screen, "Press 'c' to create a vector. Hold right click to move camera.")
	}
}

func drawGrid(g *Game, screen *ebiten.Image) {
	for x := float32(0); x < float32(g.screenWidth); x += gridSize {
		nx := x + float32(int(g.camera.X)%int(gridSize))
		vector.StrokeLine(screen, nx, 0, nx, float32(g.screenHeight), 1, gridColor, false)
	}
	for y := float32(0); y < float32(g.screenWidth); y += gridSize {
		ny := y + float32(int(g.camera.Y)%int(gridSize))
		vector.StrokeLine(screen, 0, ny, float32(g.screenWidth), ny, 1, gridColor, false)
	}
}

// Layout sets the size of the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}
