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
	miniMapOffset  = 20
	lightIntensity = 1000
	lightDither    = true
)

var (
	wallColor   color.NRGBA = HSVtoRGB(180, 0.9, 0.25)
	frameNumber int
)

func (g *Game) Draw(screen *ebiten.Image) {
	frameNumber++

	for x := 0; x < screenWidth; x++ {
		cameraX := 2*float64(x)/float64(screenWidth) - 1
		rayDir := angleToXY(player.angle+cameraX, 1)

		//Need to optimize, this is a slow way to do this
		nearestDist := math.MaxFloat64
		for _, wall := range walls {
			if dist, hit := rayIntersectsSegment(player.pos, rayDir, wall); hit {
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

			if lightDither {
				if frameNumber%2 == 0 {
					if wallColor.A > 0 && wallColor.A != 255 {
						wallColor.A = wallColor.A - 2
					}
				}
			}

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
	vector.DrawFilledCircle(screen, miniMapOffset+float32(player.pos.X)*miniMapSize, miniMapOffset+float32(player.pos.Y)*miniMapSize, 5, colornames.Yellow, false)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v, vel: %3.2f,%3.2f, angle: %3.2f, speed: %3.2f", int(ebiten.ActualFPS()), player.velocity.X, player.velocity.Y, player.angle, math.Sqrt(float64(player.velocity.X*player.velocity.X+player.velocity.Y*player.velocity.Y))))
}
