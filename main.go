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

var (
	wallColor color.NRGBA = HSVtoRGB(1.0, 0, 1.0)
)

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) Draw(screen *ebiten.Image) {

	for x := 0; x < screenWidth; x++ {
		cameraX := 2*float64(x)/float64(screenWidth) - 1
		rayDir := XY64{X: g.player.dir.X + g.player.plane.X*cameraX, Y: g.player.dir.Y + g.player.plane.Y*cameraX}

		//Need to optimize, this is a slow way to do this
		nearestDist := math.MaxFloat64
		for _, wall := range walls {
			if dist, hit := rayIntersectsSegment(g.player.pos, rayDir, wall); hit {
				if dist < nearestDist {
					nearestDist = dist
				}
			}
		}

		if nearestDist < math.MaxFloat64 {
			wallColor := wallColor

			// Calculate color with light falloff and gamma correction
			valFloat := applyFalloffWithGammaCorrection(nearestDist, 1, 1)
			value := valFloat * 255

			if value > 255 {
				value = 255
			} else if value < 0 {
				value = 0
			}
			wallColor.A = uint8(value)

			lineHeight := int(float64(screenHeight) / nearestDist)
			drawStart := -lineHeight/2 + screenHeight/2
			if drawStart < 0 {
				drawStart = 0
			}
			drawEnd := lineHeight/2 + screenHeight/2
			if drawEnd >= screenHeight {
				drawEnd = screenHeight - 1
			}

			vector.StrokeLine(screen, float32(x), float32(drawStart), float32(x), float32(drawEnd), 1, wallColor, false)
		}

	}

	const mapMag = 10
	const mapOff = 10

	for _, v := range walls {
		// Convert vector coordinates to screen coordinates
		x1, y1 := float32(v.X1), float32(v.Y1)
		x2, y2 := float32(v.X2), float32(v.Y2)

		x1, y1, x2, y2 = mapOff+x1*mapMag, mapOff+y1*mapMag, mapOff+x2*mapMag, mapOff+y2*mapMag

		vector.StrokeLine(screen, x1, y1, x2, y2, 1, colornames.Teal, false)
	}
	vector.DrawFilledCircle(screen, mapOff+float32(g.player.pos.X)*mapMag, mapOff+float32(g.player.pos.Y)*mapMag, 5, colornames.Yellow, false)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v", int(ebiten.ActualFPS())))
}

func main() {
	//walls = BoxToVectors(1, 1, 10, 10)
	//walls = append(walls, BoxToVectors(5, 5, 1, 1)...)

	game := &Game{
		player: Player{
			pos: XY64{X: 3, Y: 3}, dir: XY64{X: -1, Y: -0.25}, plane: XY64{X: -0.15, Y: 0.65},
		},
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
