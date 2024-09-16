package main

import (
	"fmt"
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

type Player struct {
	posX, posY, dirX, dirY, planeX, planeY float64
}

type Game struct {
	player Player
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

	game := &Game{
		player: Player{
			posX: 3, posY: 3, dirX: -1, dirY: -0.25, planeX: -0.15, planeY: 0.65,
		},
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
