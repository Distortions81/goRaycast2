package main

import (
	"fmt"
	"image/color"
	"math"
	"sync"

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

var renderLock sync.Mutex

func (g *Game) Draw(screen *ebiten.Image) {
	renderLock.Lock()
	defer renderLock.Unlock()

	frameNumber++

	renderScene(bspData, player.pos, player.angle, screen)
	renderMinimap(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v, vel: %3.2f,%3.2f, angle: %3.2f, speed: %3.2f", int(ebiten.ActualFPS()), player.velocity.X, player.velocity.Y, player.angle, math.Sqrt(float64(player.velocity.X*player.velocity.X+player.velocity.Y*player.velocity.Y))))
}
