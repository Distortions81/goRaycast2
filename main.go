package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
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

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.posX += g.player.dirX * 0.1
		g.player.posY += g.player.dirY * 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.posX -= g.player.dirX * 0.1
		g.player.posY -= g.player.dirY * 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(0.1) - g.player.dirY*math.Sin(0.1)
		g.player.dirY = oldDirX*math.Sin(0.1) + g.player.dirY*math.Cos(0.1)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(0.1) - g.player.planeY*math.Sin(0.1)
		g.player.planeY = oldPlaneX*math.Sin(0.1) + g.player.planeY*math.Cos(0.1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(-0.1) - g.player.dirY*math.Sin(-0.1)
		g.player.dirY = oldDirX*math.Sin(-0.1) + g.player.dirY*math.Cos(-0.1)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(-0.1) - g.player.planeY*math.Sin(-0.1)
		g.player.planeY = oldPlaneX*math.Sin(-0.1) + g.player.planeY*math.Cos(-0.1)
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

		wallColor := colornames.Red
		if side == 1 {
			wallColor = colornames.Darkred
		}

		ebitenutil.DrawLine(screen, float64(x), float64(drawStart), float64(x), float64(drawEnd), wallColor)

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
