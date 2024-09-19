package main

import (
	"fmt"
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

	renderScene(bspData, player.pos, player.angle, screen)

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
