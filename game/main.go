package main

import (
	"image"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	mapScaleDiv  = 50
	levelPath    = "../level1.txt"
	spriteFile   = "test.png"
)

var (
	walls   = []Line32{}
	wallImg *ebiten.Image
	bspData *BSPNode
)

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

var (
	textureBounds image.Rectangle
	textureWidth,
	textureHeight int
	textureRepeatDistance float32 = 1.0
	workSize              int     = 32
	rayList               [screenWidth]renderData
)

func main() {
	player = playerData{
		pos: pos32{X: 3, Y: 3}, angle: 4,
	}

	ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")

	readVecs()

	bspData = buildBSPTree(walls)

	//Update level if written
	go func() {
		var oldModTime time.Time
		for {
			time.Sleep(time.Millisecond * 500)
			stat, _ := os.Stat(levelPath)
			if stat.ModTime() != oldModTime {
				oldModTime = stat.ModTime()
				readVecs()
			}
		}
	}()

	//Load sprite
	var err error
	wallImg, _, err = ebitenutil.NewImageFromFile(spriteFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Precompute texture width and height (cached outside the loop in the main render function)
	textureBounds = wallImg.Bounds()
	textureWidth = textureBounds.Dx()
	textureHeight = textureBounds.Dy()
	workSize = int(math.Round(float64(screenWidth)/float64(runtime.NumCPU()))) / 2

	//Start game
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

func readVecs() {

	data, err := os.ReadFile(levelPath)
	if err != nil {
		log.Fatalln("Unable to read " + levelPath)
	}

	tmp := []Line32{}
	text := string(data)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		args := strings.Split(line, ",")
		if len(args) != 4 {
			continue
		}
		x1, _ := strconv.ParseFloat(args[0], 32)
		y1, _ := strconv.ParseFloat(args[1], 32)
		x2, _ := strconv.ParseFloat(args[2], 32)
		y2, _ := strconv.ParseFloat(args[3], 32)

		tmp = append(tmp, Line32{X1: float32(x1) / mapScaleDiv, Y1: float32(y1) / mapScaleDiv, X2: float32(x2) / mapScaleDiv, Y2: float32(y2) / mapScaleDiv})
	}

	renderLock.Lock()
	walls = tmp
	defer renderLock.Unlock()
}
