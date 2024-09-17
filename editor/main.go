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
	vectors []Vector2D

	cameraX        float64
	cameraY        float64
	startX, startY float64

	lastMouseX, lastMouseY int

	createMode,
	firstClick,
	secondClick,
	dragging bool
}

// Update is called every frame
func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()

	if ebiten.IsKeyPressed(ebiten.KeyC) && !g.createMode {
		g.createMode = true
		g.firstClick = false
		g.secondClick = false
		g.startX = 0
		g.startY = 0
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.createMode {
			if !g.secondClick && !g.firstClick {
				// Start creating a new vector
				g.startX = float64(mouseX) - g.cameraX
				g.startY = float64(mouseY) - g.cameraY
				g.firstClick = true
			} else if g.firstClick && !g.secondClick {
				// Finish creating the vector
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
				g.createMode = false
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.cameraX += float64(mouseX - int(g.lastMouseX))
		g.cameraY += float64(mouseY - int(g.lastMouseY))
	}

	g.lastMouseX, g.lastMouseY = mouseX, mouseY
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
		x1 := vec.X1 + g.cameraX
		y1 := vec.Y1 + g.cameraY
		x2 := vec.X2 + g.cameraX
		y2 := vec.Y2 + g.cameraY
		vector.StrokeLine(vectorScreen, float32(x1), float32(y1), float32(x2), float32(y2), 1, color.White, true)
	}

	if g.createMode && g.firstClick && !g.secondClick {
		mouseX, mouseY := ebiten.CursorPosition()
		vector.StrokeLine(screen, float32(g.startX), float32(g.startY), float32(mouseX), float32(mouseY), 1, colornames.Red, true)
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
