package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

// Define a struct for a 2D vector with start and end points
type Vector2D struct {
	X1, Y1, X2, Y2 float64
}

// Game struct to hold game state
type Game struct {
	vectors        []Vector2D
	cameraX        float64
	cameraY        float64
	startX, startY float64
	createMode,
	firstClick,
	secondClick bool
}

// Update is called every frame
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyC) && !g.createMode {
		g.createMode = true
		g.firstClick = false
		g.secondClick = false
		g.startX = 0
		g.startY = 0
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.createMode {
			if !g.secondClick && !g.firstClick {
				// Start creating a new vector
				mouseX, mouseY := ebiten.CursorPosition()
				g.startX = float64(mouseX) - g.cameraX
				g.startY = float64(mouseY) - g.cameraY
				g.firstClick = true
			} else if g.firstClick && !g.secondClick {
				// Finish creating the vector
				mouseX, mouseY := ebiten.CursorPosition()
				endX := float64(mouseX) - g.cameraX
				endY := float64(mouseY) - g.cameraY
				g.vectors = append(g.vectors, Vector2D{
					X1: g.startX,
					Y1: g.startY,
					X2: endX,
					Y2: endY,
				})
				fmt.Printf("created: %v,%v - %v,%v\n", g.startX, g.startY, endX, endY)
				g.secondClick = true
			}
		}
	}
	return nil
}

// Draw is called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a white background
	screen.Fill(color.White)

	// Create a new vector image to draw on
	vectorImg := image.NewRGBA(image.Rect(0, 0, 640, 480))
	vectorScreen := ebiten.NewImageFromImage(vectorImg)

	// Draw each vector with respect to the camera position
	for _, vec := range g.vectors {
		x1 := vec.X1 + g.cameraX
		y1 := vec.Y1 + g.cameraY
		x2 := vec.X2 + g.cameraX
		y2 := vec.Y2 + g.cameraY
		vector.StrokeLine(vectorScreen, float32(x1), float32(y1), float32(x2), float32(y2), 1, color.Black, true)
	}

	if g.createMode && g.firstClick && !g.secondClick {
		mouseX, mouseY := ebiten.CursorPosition()
		vector.StrokeLine(screen, float32(g.startX), float32(g.startY), float32(mouseX), float32(mouseY), 1, colornames.Darkred, true)
	}

	// Draw the vector image onto the screen
	screen.DrawImage(vectorScreen, nil)

	// Draw text for clarity
	ebitenutil.DebugPrint(screen, "Press 'C' to Create Vector\nLeft Click to Drag Vector\nRight Click and Drag to Move Camera")
}

// Check if a point is on a line segment
func isPointOnLine(px, py, x1, y1, x2, y2 float64) bool {
	// Check if the point is within the bounding box of the line segment
	if px < min(x1, x2) || px > max(x1, x2) || py < min(y1, y2) || py > max(y1, y2) {
		return false
	}

	// Calculate the distance from the point to the line segment
	dx := x2 - x1
	dy := y2 - y1
	// Avoid division by zero
	if dx == 0 && dy == 0 {
		return px == x1 && py == y1
	}

	t := ((px-x1)*dx + (py-y1)*dy) / (dx*dx + dy*dy)
	if t < 0 || t > 1 {
		return false
	}

	// Projected point
	projX := x1 + t*dx
	projY := y1 + t*dy
	return px == projX && py == projY
}

// Helper functions for bounding box calculations
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Layout sets the size of the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// Create a new game instance
	game := &Game{}

	// Start the Ebiten game loop
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
