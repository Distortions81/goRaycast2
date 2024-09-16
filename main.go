package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 640
	screenHeight = 480
	mapWidth     = 8
	mapHeight    = 8
)

var worldMap = [mapWidth][mapHeight]int{
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 120, 0, 0, 1},
	{1, 0, 0, 180, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1},
}

type Player struct {
	posX, posY, dirX, dirY, planeX, planeY float64
}

type Game struct {
	player Player
}

const turnSpeed = 0.05

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.posX += g.player.dirX * turnSpeed
		g.player.posY += g.player.dirY * turnSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.posX -= g.player.dirX * turnSpeed
		g.player.posY -= g.player.dirY * turnSpeed
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
	return nil
}

func (g *Game) Layout(h int, w int) (int, int) {
	return h, w
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)

	for x := 0; x < screenWidth; x++ {
		cameraX := 2*float64(x)/float64(screenWidth) - 1
		rayDirX := g.player.dirX + g.player.planeX*cameraX
		rayDirY := g.player.dirY + g.player.planeY*cameraX

		mapX := int(g.player.posX)
		mapY := int(g.player.posY)

		deltaDistX := math.Abs(1 / rayDirX)
		deltaDistY := math.Abs(1 / rayDirY)

		var stepX, stepY int
		var sideDistX, sideDistY float64

		if rayDirX < 0 {
			stepX = -1
			sideDistX = (g.player.posX - float64(mapX)) * deltaDistX
		} else {
			stepX = 1
			sideDistX = (float64(mapX+1) - g.player.posX) * deltaDistX
		}
		if rayDirY < 0 {
			stepY = -1
			sideDistY = (g.player.posY - float64(mapY)) * deltaDistY
		} else {
			stepY = 1
			sideDistY = (float64(mapY+1) - g.player.posY) * deltaDistY
		}

		var hit, side int
		for hit == 0 {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = 0
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = 1
			}

			if mapX >= mapWidth {
				side = -1
				break
			} else if mapX < 0 {
				side = -1
				break
			}
			if mapY >= mapHeight {
				side = -1
				break
			} else if mapY < 0 {
				side = -1
				break
			}

			if worldMap[mapX][mapY] > 0 {
				hit = 1
			}
		}

		var perpWallDist float64
		if side == 0 {
			perpWallDist = (float64(mapX) - g.player.posX + (1-float64(stepX))/2) / rayDirX
		} else {
			perpWallDist = (float64(mapY) - g.player.posY + (1-float64(stepY))/2) / rayDirY
		}

		lineHeight := int(float64(screenHeight) / perpWallDist)
		drawStart := -lineHeight/2 + screenHeight/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + screenHeight/2
		if drawEnd >= screenHeight {
			drawEnd = screenHeight - 1
		}

		var wallColor color.NRGBA
		if side == 0 {
			wallColor = HSVtoRGB(float64(worldMap[mapX][mapY]), 1.0, 1.0)
		} else if side == 1 {
			wallColor = HSVtoRGB(float64(worldMap[mapX][mapY]), 1.0, 0.7)
		} else {
			return
		}

		ldist := (perpWallDist * 3)
		dist := 255 - (ldist*ldist)/2
		if dist > 255 {
			dist = 255
		} else if dist < 0 {
			dist = 0
		}
		wallColor.A = uint8(dist)

		vector.DrawFilledRect(screen, float32(x), float32(drawStart), 1, float32(drawEnd-drawStart), wallColor, false)

	}
}

func main() {
	game := &Game{
		player: Player{
			posX:   4.0,
			posY:   4.0,
			dirX:   -1.0,
			dirY:   0.0,
			planeX: 0.0,
			planeY: 0.66,
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Minimal Raycaster")
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
