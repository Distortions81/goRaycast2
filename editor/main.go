package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	snapDistance = 10
	lineWidth    = 2
)

// Define a struct for a 2D vector with start and end points
type Vector2D struct {
	X1, Y1, X2, Y2 float64
}

type XY struct {
	X, Y float64
}

// Game struct to hold game state
type Game struct {
	vectors []Vector2D

	camera,
	start,
	lastMouse XY

	createMode,
	firstClick,
	secondClick bool
}

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
			snappedPos := SnapPosition(wpos, g.vectors, snapDistance)
			if !g.secondClick && !g.firstClick {

				// Start creating a new vector
				g.start = snappedPos
				g.firstClick = true
			} else if g.firstClick && !g.secondClick {
				// Finish creating the vector
				endX := snappedPos.X
				endY := snappedPos.Y
				g.vectors = append(g.vectors, Vector2D{
					X1: g.start.X,
					Y1: g.start.Y,
					X2: endX,
					Y2: endY,
				})
				fmt.Printf("created: %v,%v - %v,%v\n", g.start.X, g.start.Y, endX, endY)
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

// Draw is called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a white background
	screen.Fill(color.Black)

	// Create a new vector image to draw on
	vectorImg := image.NewRGBA(image.Rect(0, 0, 1280, 1440))
	vectorScreen := ebiten.NewImageFromImage(vectorImg)

	// Draw each vector with respect to the camera position
	for _, vec := range g.vectors {
		x1 := vec.X1 + g.camera.X
		y1 := vec.Y1 + g.camera.Y
		x2 := vec.X2 + g.camera.X
		y2 := vec.Y2 + g.camera.Y
		vector.StrokeLine(vectorScreen, float32(x1), float32(y1), float32(x2), float32(y2), lineWidth, color.White, true)
	}

	if g.createMode {
		mouseX, mouseY := ebiten.CursorPosition()
		mpos := XY{X: float64(mouseX), Y: float64(mouseY)}
		snappedPos := SnapPosition(mpos, g.vectors, snapDistance)

		if g.createMode {
			vector.DrawFilledCircle(screen, float32(snappedPos.X), float32(snappedPos.Y), lineWidth*2, colornames.Yellow, true)
		}
		if g.firstClick && !g.secondClick {
			vector.StrokeLine(screen, float32(g.start.X), float32(g.start.Y), float32(snappedPos.X), float32(snappedPos.Y), lineWidth, colornames.Red, true)
		}
	}

	// Draw the vector image onto the screen
	screen.DrawImage(vectorScreen, nil)

	// Draw text for clarity
	if g.createMode {
		ebitenutil.DebugPrint(screen, "Vector created, click again to specify vector end.")
	} else {
		ebitenutil.DebugPrint(screen, "Press 'c' to create a vector. Hold right click to move camera.")
	}
}

// Layout sets the size of the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Create a new game instance
	game := &Game{}

	// Start the Ebiten game loop
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func distance(p1, p2 XY) float64 {
	return math.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}

// SnapPosition snaps a new position to the nearest existing position within a threshold
func SnapPosition(newPos XY, existingPositions []Vector2D, threshold float64) XY {
	var snappedPosition XY = newPos
	minDistance := threshold // Initialize with threshold to ensure snapping only within the threshold

	for _, pos := range existingPositions {

		apos := XY{X: pos.X1, Y: pos.Y1}
		bpos := XY{X: pos.X2, Y: pos.Y2}

		dist := distance(newPos, apos)
		if dist < minDistance {
			minDistance = dist
			snappedPosition = apos
		}

		dist = distance(newPos, bpos)
		if dist < minDistance {
			minDistance = dist
			snappedPosition = bpos
		}
	}

	return snappedPosition
}
