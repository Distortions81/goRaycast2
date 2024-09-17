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
		// Map screen X to camera plane (-1 to 1)
		cameraX := 2*float64(x)/float64(screenWidth) - 1

		// Adjust the ray direction based on player's current angle + offset from center view
		rayDir := angleToXY(player.angle+math.Atan(cameraX), 1)

		// Calculate distance to the nearest wall
		nearestDist := math.MaxFloat64
		for _, wall := range walls {
			if dist, hit := rayIntersectsSegment(player.pos, rayDir, wall); hit {
				// Correct the fisheye effect by adjusting with cos of angle offset
				correctedDist := dist * math.Cos(math.Atan(cameraX))
				if correctedDist < nearestDist {
					nearestDist = correctedDist
				}
			}
		}

		if nearestDist < math.MaxFloat64 {
			// Light falloff and gamma correction (unchanged)
			wallColor := wallColor
			valFloat := applyFalloff(nearestDist, lightIntensity, float64(wallColor.R+wallColor.G+wallColor.B)/765.0/3.0)
			wallColor.A = uint8(valFloat * 255)

			if lightDither && frameNumber%2 == 0 && wallColor.A > 0 && wallColor.A != 255 {
				wallColor.A = wallColor.A - 2
			}

			// Draw the wall slice with corrected distance
			lineHeight := int(float64(screenHeight) / nearestDist)
			drawStart := -lineHeight/2 + screenHeight/2
			if drawStart < 0 {
				drawStart = 0
			}
			drawEnd := lineHeight/2 + screenHeight/2
			if drawEnd >= screenHeight {
				drawEnd = screenHeight - 1
			}

			// Draw the line representing the wall column
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
