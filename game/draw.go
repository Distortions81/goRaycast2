package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	miniMapSize     = 5
	miniMapOffset   = 20
	lightIntensity  = 1000
	lightDither     = true
	occlusionRadius = 2.0 // Adjust as needed
	numSamples      = 32  // Number of samples to take around each pixel
)

var (
	wallColor   color.NRGBA = HSVtoRGB(180, 0.0, 0.8)
	frameNumber int
)

var renderLock sync.Mutex

func (g *Game) Draw(screen *ebiten.Image) {
	renderLock.Lock()
	defer renderLock.Unlock()

	frameNumber++

	// Precompute texture width and height
	textureBounds := wallImg.Bounds()
	textureWidth := textureBounds.Dx()
	textureHeight := textureBounds.Dy()

	// Precompute constants for screen dimensions
	screenCenter := screenHeight / 2

	for x := 0; x < screenWidth; x++ {
		// 1. Map screen X to camera plane (-1 to 1)
		cameraX := 2*float64(x)/float64(screenWidth) - 1

		// 2. Adjust ray direction based on player's current angle + offset from center
		rayDir := angleToXY(player.angle+math.Atan(cameraX), 1)

		// Variables to track nearest wall and hit position
		nearestDist := math.MaxFloat64
		var nearwall Vec64
		var hitPos XY64 // Intersection point on the wall

		// Loop through all the walls in the scene to find the nearest wall hit by the ray
		for _, wall := range walls {
			if dist, hPos, hit := rayIntersectsSegment(player.pos, rayDir, wall); hit {
				// Correct the distance to avoid fish-eye effect
				correctedDist := dist * math.Cos(math.Atan(cameraX))
				if correctedDist < nearestDist {
					nearestDist = correctedDist
					hitPos = hPos
					nearwall = wall
				}
			}
		}

		if nearestDist < math.MaxFloat64 {
			// 3. Shading Calculation (done per column, not per pixel)
			valFloat := applyFalloff(nearestDist, lightIntensity, float64(wallColor.R+wallColor.G+wallColor.B)/765.0/3.0)
			wallColor.A = uint8(valFloat * 255)

			// Calculate line height and start/end of the slice
			lineHeight := int(float64(screenHeight) / nearestDist)
			drawStart := -lineHeight/2 + screenCenter
			if drawStart < 0 {
				drawStart = 0
			}
			drawEnd := lineHeight/2 + screenCenter
			if drawEnd >= screenHeight {
				drawEnd = screenHeight - 1
			}

			// 4. Calculate the correct X position on the texture based on hitPos

			// Wall's direction vector (the vector from wall start to wall end)
			wallDirX := nearwall.X2 - nearwall.X1
			wallDirY := nearwall.Y2 - nearwall.Y1
			wallLength := math.Sqrt(wallDirX*wallDirX + wallDirY*wallDirY)

			// Normalize the direction vector
			wallDirX /= wallLength
			wallDirY /= wallLength
			dx := hitPos.X - nearwall.X1
			dy := hitPos.Y - nearwall.Y1
			wallHitPosition := (dx*wallDirX + dy*wallDirY)

			// Convert wallHitPosition into texture space
			textureX := int((wallHitPosition * float64(textureWidth))) % textureWidth
			if textureX < 0 {
				textureX += textureWidth
			}

			// 5. Texture clipping and scaling
			textureStep := float64(textureHeight) / float64(lineHeight) // Step size for texture Y
			textureY := 0.0

			// If the wall height is larger than the screen, adjust textureY and clip the texture
			if lineHeight > screenHeight {
				textureY = float64(lineHeight-screenHeight) / 2 * textureStep
				drawStart = 0 // Clamp drawStart to 0 (top of screen)
			}

			// 6. Create a sub-image of the texture slice to draw (from textureX to textureX + 1)
			srcRect := image.Rect(textureX, int(textureY), textureX+1, textureHeight)
			textureSlice := wallImg.SubImage(srcRect).(*ebiten.Image)

			// 7. Apply shading and draw the texture slice
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
			op.GeoM.Scale(1, float64(lineHeight)/float64(textureHeight)) // Scale texture to line height
			op.GeoM.Translate(float64(x), float64(drawStart))            // Position the texture slice
			op.ColorM.Scale(float64(wallColor.R)/255, float64(wallColor.G)/255, float64(wallColor.B)/255, float64(wallColor.A)/255)

			// Draw the texture slice
			screen.DrawImage(textureSlice, op)
		}
	}

	//Minimap
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

func RelativeCrop(source *ebiten.Image, r image.Rectangle) *ebiten.Image {
	rx, ry := source.Bounds().Min.X+r.Min.X, source.Bounds().Min.Y+r.Min.Y
	return source.SubImage(image.Rect(rx, ry, rx+r.Max.X, ry+r.Max.Y)).(*ebiten.Image)
}
