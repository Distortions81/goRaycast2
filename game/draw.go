package main

import (
	"fmt"
	"image/color"
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v", int(ebiten.ActualFPS())))
}
