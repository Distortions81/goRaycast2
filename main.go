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

var walls = []Vector{
	// Outer boundary
	{0.0, 0.0, 25.0, 0.0},
	{25.0, 0.0, 25.0, 25.0},
	{25.0, 25.0, 0.0, 25.0},
	{0.0, 25.0, 0.0, 0.0},

	// Main corridor
	{2.5, 2.5, 22.5, 2.5},
	{22.5, 2.5, 22.5, 22.5},
	{22.5, 22.5, 2.5, 22.5},
	{2.5, 22.5, 2.5, 2.5},

	// Rooms
	{5.0, 5.0, 10.0, 5.0},
	{10.0, 5.0, 10.0, 10.0},
	{10.0, 10.0, 5.0, 10.0},
	{5.0, 10.0, 5.0, 5.0},

	{15.0, 5.0, 20.0, 5.0},
	{20.0, 5.0, 20.0, 10.0},
	{20.0, 10.0, 15.0, 10.0},
	{15.0, 10.0, 15.0, 5.0},

	// Additional corridors
	{5.0, 15.0, 10.0, 15.0},
	{10.0, 15.0, 10.0, 20.0},
	{10.0, 20.0, 5.0, 20.0},
	{5.0, 20.0, 5.0, 15.0},

	{15.0, 15.0, 20.0, 15.0},
	{20.0, 15.0, 20.0, 20.0},
	{20.0, 20.0, 15.0, 20.0},
	{15.0, 20.0, 15.0, 15.0},
}

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

const turnSpeed = 0.03
const moveSpeed = 0.1

func (g *Game) Update() error {
	oldPlayer := g.player
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.posX += g.player.dirX * moveSpeed
		g.player.posY += g.player.dirY * moveSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.posX -= g.player.dirX * moveSpeed
		g.player.posY -= g.player.dirY * moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(turnSpeed) - g.player.dirY*math.Sin(turnSpeed)
		g.player.dirY = oldDirX*math.Sin(turnSpeed) + g.player.dirY*math.Cos(turnSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(turnSpeed) - g.player.planeY*math.Sin(turnSpeed)
		g.player.planeY = oldPlaneX*math.Sin(turnSpeed) + g.player.planeY*math.Cos(turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(-turnSpeed) - g.player.dirY*math.Sin(-turnSpeed)
		g.player.dirY = oldDirX*math.Sin(-turnSpeed) + g.player.dirY*math.Cos(-turnSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(-turnSpeed) - g.player.planeY*math.Sin(-turnSpeed)
		g.player.planeY = oldPlaneX*math.Sin(-turnSpeed) + g.player.planeY*math.Cos(-turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.Key0) {
		fmt.Printf("posX: %v, posY: %v, dirX: %v, dirY: %v, planeX: %v, planeY: %v\n", g.player.posX, g.player.posY, g.player.dirX, g.player.dirY, g.player.planeX, g.player.planeY)
	}

	for _, wall := range walls {
		if dist, hit := rayIntersectsSegment(g.player.posX, g.player.posY, g.player.dirX, g.player.dirY, wall); hit {
			if dist < 0.5 {
				g.player = oldPlayer
			}
		}
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

			ldist := (nearestDist)
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

	const mapMag = 10
	const mapOff = 10

	for _, v := range walls {
		// Convert vector coordinates to screen coordinates
		x1, y1 := float32(v.X1), float32(v.Y1)
		x2, y2 := float32(v.X2), float32(v.Y2)

		x1, y1, x2, y2 = mapOff+x1*mapMag, mapOff+y1*mapMag, mapOff+x2*mapMag, mapOff+y2*mapMag

		// Draw lines as filled rectangles
		vector.StrokeLine(screen, x1, y1, x2, y2, 1, colornames.Teal, false)
	}
	vector.DrawFilledCircle(screen, mapOff+float32(g.player.posX)*mapMag, mapOff+float32(g.player.posY)*mapMag, 5, colornames.Yellow, false)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v", int(ebiten.ActualFPS())))
}

func main() {
	//walls = BoxToVectors(1, 1, 10, 10)
	//walls = append(walls, BoxToVectors(5, 5, 1, 1)...)
	//fmt.Println(walls)

	game := &Game{
		player: Player{
			posX: 2.986815779357488, posY: 3.0406601123043306, dirX: -0.9751973713086025, dirY: -0.22133704387836103, planeX: -0.14608244895971825, planeY: 0.6436302650636798,
		},
	}

	ebiten.SetVsyncEnabled(true)
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
