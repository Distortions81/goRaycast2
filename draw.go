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
	miniMapSize    = 5
	miniMapOffset  = 5
	lightIntensity = 500
)

var (
	wallColor color.NRGBA = HSVtoRGB(180, 0.9, 0.25)
)

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
			valFloat := applyFalloff(nearestDist, lightIntensity, (float64(wallColor.R+wallColor.G+wallColor.B) / 765 / 3.0))
			wallColor.A = uint8(valFloat * 255)

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

	for _, v := range walls {
		// Convert vector coordinates to screen coordinates
		x1, y1 := float32(v.X1), float32(v.Y1)
		x2, y2 := float32(v.X2), float32(v.Y2)

		x1, y1, x2, y2 = miniMapOffset+x1*miniMapSize, miniMapOffset+y1*miniMapSize, miniMapOffset+x2*miniMapSize, miniMapOffset+y2*miniMapSize

		vector.StrokeLine(screen, x1, y1, x2, y2, 1, colornames.Teal, false)
	}
	vector.DrawFilledCircle(screen, miniMapOffset+float32(g.player.pos.X)*miniMapSize, miniMapOffset+float32(g.player.pos.Y)*miniMapSize, 5, colornames.Yellow, false)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v", int(ebiten.ActualFPS())))
}
