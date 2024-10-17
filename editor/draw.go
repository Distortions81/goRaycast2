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
	screen.Fill(color.Black)

	if bgImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.ColorScale.ScaleAlpha(0.3)
		screen.DrawImage(bgImage, op)
	}

	drawGrid(g, screen)

	// Draw each vector with respect to the camera position
	for _, vec := range walls {
		x1 := vec.X1 + g.camera.X
		y1 := vec.Y1 + g.camera.Y
		x2 := vec.X2 + g.camera.X
		y2 := vec.Y2 + g.camera.Y
		vector.StrokeLine(screen, (x1), (y1), (x2), (y2), lineWidth, color.White, true)
	}

	if g.createMode {
		mouseX, mouseY := ebiten.CursorPosition()
		mpos := pos32{X: float32(mouseX), Y: float32(mouseY)}
		snappedPos := snapPos(mpos, walls, lineSnapDist)

		if g.createMode {
			vector.DrawFilledCircle(screen, (snappedPos.X), (snappedPos.Y), lineWidth*2, colornames.Yellow, true)
		}
		if g.firstClick && !g.secondClick {
			vector.StrokeLine(screen, (g.start.X), (g.start.Y), (snappedPos.X), (snappedPos.Y), lineWidth, colornames.Red, true)
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
	for x := (0); x < (g.screenWidth); x += gridSize {
		nx := float32(x + (int(g.camera.X) % int(gridSize)))
		vector.StrokeLine(screen, nx, 0, nx, float32(g.screenHeight), 1, gridColor, false)
	}
	for y := (0); y < (g.screenWidth); y += gridSize {
		ny := float32(y + (int(g.camera.Y) % int(gridSize)))
		vector.StrokeLine(screen, 0, ny, float32(g.screenWidth), ny, 1, gridColor, false)
	}
}

// Layout sets the size of the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}
