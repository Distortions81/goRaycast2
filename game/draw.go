package main

import (
	"fmt"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	lightIntensity = 1000
)

var (
	wallColor   color.NRGBA = HSVtoRGB(180, 0.0, 0.8)
	frameNumber int
)

var (
	renderLock sync.Mutex
	worstFrame int
	bestFrame  int = 10000
)

func (g *Game) Draw(screen *ebiten.Image) {
	renderLock.Lock()
	defer renderLock.Unlock()

	frameNumber++
	start := time.Now()

	//renderFloorAndCeiling(screen)
	renderScene(screen)
	renderMinimap(screen)

	took := time.Since(start).Microseconds()
	if frameNumber%6000 == 0 {
		worstFrame = 0
		bestFrame = 10000
	}
	worstFrame = max(worstFrame, int(took))
	bestFrame = min(bestFrame, int(took))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %3v, Took: %4vus / Max: %4vus / Min: %4vus", int(ebiten.ActualFPS()), took, worstFrame, bestFrame))
}

var FOVDeg float32 = 60
var fovRadians = FOVDeg * (math.Pi / 180.0)
var planeLength = math32.Tan(fovRadians / 2.0)

func renderFloorAndCeiling(screenImage *ebiten.Image) {
	textureWidth, textureHeight := wallImg.Bounds().Dx(), wallImg.Bounds().Dy()
	screenWidth, screenHeight := screenImage.Size()
	posX, posY := float32(player.pos.X), float32(player.pos.Y)
	dir := angleToXY(player.angle, 1)

	// Use the same planeX and planeY as in wall rendering
	planeX := -dir.Y * planeLength
	planeY := dir.X * planeLength

	posZ := 0.5 * float32(screenHeight)

	for y := screenHeight / 2; y < screenHeight; y++ {
		p := y - screenHeight/2
		if p == 0 {
			continue // Avoid division by zero
		}

		rowDistance := -posZ / float32(p)

		// Leftmost and rightmost ray directions
		rayDirX0 := dir.X + planeX*(-1)
		rayDirY0 := dir.Y + planeY*(-1)
		rayDirX1 := dir.X + planeX*1
		rayDirY1 := dir.Y + planeY*1

		// Step size per screen pixel
		floorStepX := (rayDirX1 - rayDirX0) / float32(screenWidth)
		floorStepY := (rayDirY1 - rayDirY0) / float32(screenWidth)

		// Starting position
		floorX := posX + rowDistance*rayDirX0
		floorY := posY + rowDistance*rayDirY0

		for x := 0; x < screenWidth; x++ {
			// The cell coordinates
			cellX := int(floorX)
			cellY := int(floorY)

			// Texture coordinates
			tx := int((floorX-float32(cellX))*float32(textureWidth)) % textureWidth
			ty := int((floorY-float32(cellY))*float32(textureHeight)) % textureHeight

			if tx < 0 {
				tx += textureWidth
			}
			if ty < 0 {
				ty += textureHeight
			}

			// Sample the floor and ceiling textures
			floorColor := wallImg.At(tx, ty)
			ceilingColor := wallImg.At(tx, ty) // Use ceiling texture if available

			// Set the floor pixel
			screenImage.Set(x, y, floorColor)

			// Set the ceiling pixel
			screenImage.Set(x, screenHeight-y-1, ceilingColor)

			// Move to the next position
			floorX += floorStepX * rowDistance
			floorY += floorStepY * rowDistance
		}
	}
}
