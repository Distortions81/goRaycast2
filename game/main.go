package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	mapScaleDiv  = 50
	levelPath    = "../editor/vecs.txt"
)

var walls = []Vec64{}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func main() {
	player = playerData{
		pos: XY64{X: 3, Y: 3}, angle: 4,
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Raycaster with Vectors")

	readVecs()

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

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

func readVecs() {

	data, err := os.ReadFile(levelPath)
	if err != nil {
		log.Fatalln("Unable to read " + levelPath)
	}

	tmp := []Vec64{}
	text := string(data)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		args := strings.Split(line, ",")
		if len(args) != 4 {
			continue
		}
		x1, _ := strconv.ParseFloat(args[0], 64)
		y1, _ := strconv.ParseFloat(args[1], 64)
		x2, _ := strconv.ParseFloat(args[2], 64)
		y2, _ := strconv.ParseFloat(args[3], 64)

		tmp = append(tmp, Vec64{X1: x1 / mapScaleDiv, Y1: y1 / mapScaleDiv, X2: x2 / mapScaleDiv, Y2: y2 / mapScaleDiv})
	}

	renderLock.Lock()
	walls = tmp
	defer renderLock.Unlock()
}
