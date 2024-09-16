package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Vector struct {
	X1, Y1, X2, Y2 float64
}

var walls = []Vector{}

func BoxToVectors(x, y, width, height float64) []Vector {
	// Define the four corners of the box
	topLeft := Vector{X1: x, Y1: y, X2: x + width, Y2: y}                       // Top edge
	topRight := Vector{X1: x + width, Y1: y, X2: x + width, Y2: y + height}     // Right edge
	bottomRight := Vector{X1: x + width, Y1: y + height, X2: x, Y2: y + height} // Bottom edge
	bottomLeft := Vector{X1: x, Y1: y + height, X2: x, Y2: y}                   // Left edge

	// Return the four edges of the box
	return []Vector{topLeft, topRight, bottomRight, bottomLeft}
}

type Player struct {
	posX, posY, dirX, dirY, planeX, planeY float64
}

type Game struct {
	player Player
}

const pSpeed = 0.05

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.posX += g.player.dirX * pSpeed
		g.player.posY += g.player.dirY * pSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.posX -= g.player.dirX * pSpeed
		g.player.posY -= g.player.dirY * pSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(pSpeed) - g.player.dirY*math.Sin(pSpeed)
		g.player.dirY = oldDirX*math.Sin(pSpeed) + g.player.dirY*math.Cos(pSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(pSpeed) - g.player.planeY*math.Sin(pSpeed)
		g.player.planeY = oldPlaneX*math.Sin(pSpeed) + g.player.planeY*math.Cos(pSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(-pSpeed) - g.player.dirY*math.Sin(-pSpeed)
		g.player.dirY = oldDirX*math.Sin(-pSpeed) + g.player.dirY*math.Cos(-pSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(-pSpeed) - g.player.planeY*math.Sin(-pSpeed)
		g.player.planeY = oldPlaneX*math.Sin(-pSpeed) + g.player.planeY*math.Cos(-pSpeed)
	}
	return nil
}

func rayIntersectsSegment(px, py, rayDirX, rayDirY float64, wall Vector) (float64, bool) {
	// Using line intersection formula
	x1, y1, x2, y2 := wall.X1, wall.Y1, wall.X2, wall.Y2

	denom := (x1-x2)*(py+rayDirY-py) - (y1-y2)*(px+rayDirX-px)
	if denom == 0 {
		return 0, false // Parallel lines
	}

	// t and u parameters for intersection formula
	t := ((x1-px)*(py+rayDirY-py) - (y1-py)*(px+rayDirX-px)) / denom
	u := -((x1-x2)*(y1-py) - (y1-y2)*(x1-px)) / denom

	// If t and u are valid, we have an intersection
	if t >= 0 && t <= 1 && u > 0 {
		return u, true
	}

	return 0, false
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)

	for x := 0; x < screenWidth; x++ {
		cameraX := 2*float64(x)/float64(screenWidth) - 1
		rayDirX := g.player.dirX + g.player.planeX*cameraX
		rayDirY := g.player.dirY + g.player.planeY*cameraX

		nearestDist := math.MaxFloat64
		for _, wall := range walls {
			if dist, hit := rayIntersectsSegment(g.player.posX, g.player.posY, rayDirX, rayDirY, wall); hit {
				if dist < nearestDist {
					nearestDist = dist
				}
			}
		}

		if nearestDist < math.MaxFloat64 {
			wallColor := HSVtoRGB(1.0, 1.0, 1.0)

			ldist := (nearestDist * 2)
			dist := 255 - (ldist*ldist)/2
			if dist > 255 {
				dist = 255
			} else if dist < 0 {
				dist = 0
			}
			wallColor.A = uint8(dist)

			lineHeight := int(float64(screenHeight) / nearestDist)
			drawStart := -lineHeight/2 + screenHeight/2
			if drawStart < 0 {
				drawStart = 0
			}
			drawEnd := lineHeight/2 + screenHeight/2
			if drawEnd >= screenHeight {
				drawEnd = screenHeight - 1
			}

			vector.DrawFilledRect(screen, float32(x), float32(drawStart), 1, float32(drawEnd-drawStart), wallColor, false)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v", int(ebiten.ActualFPS())))
}

func main() {
	walls = BoxToVectors(1, 1, 10, 10)
	walls = append(walls, BoxToVectors(5, 5, 1, 1)...)
	fmt.Println(walls)

	game := &Game{
		player: Player{
			posX:   2.0,
			posY:   2.0,
			dirX:   -1.0,
			dirY:   0.0,
			planeX: 0.0,
			planeY: 0.66,
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

// HSVtoRGB converts HSV values to RGB
func HSVtoRGB(h, s, v float64) color.NRGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r1, g1, b1 float64

	if h >= 0 && h < 60 {
		r1, g1, b1 = c, x, 0
	} else if h >= 60 && h < 120 {
		r1, g1, b1 = x, c, 0
	} else if h >= 120 && h < 180 {
		r1, g1, b1 = 0, c, x
	} else if h >= 180 && h < 240 {
		r1, g1, b1 = 0, x, c
	} else if h >= 240 && h < 300 {
		r1, g1, b1 = x, 0, c
	} else if h >= 300 && h < 360 {
		r1, g1, b1 = c, 0, x
	}

	// Convert to RGB by adding m and scaling to the range of 0-255
	r := (r1 + m) * 255
	g := (g1 + m) * 255
	b := (b1 + m) * 255

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 1}
}
