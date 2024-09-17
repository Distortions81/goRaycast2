package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func (g *Game) writeVecs() {
	buf := ""

	for _, item := range walls {
		buf = buf + fmt.Sprintf("%v,%v,%v,%v\n", item.X1/scaleDiv, item.Y1/scaleDiv, item.X2/scaleDiv, item.Y2/scaleDiv)
	}

	os.WriteFile(levelPath, []byte(buf), 0755)
}

func readVecs() {
	data, err := os.ReadFile(levelPath)
	if err != nil {
		log.Fatalln("Unable to read " + levelPath)
	}

	walls = []Vector2D{}
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

		walls = append(walls, Vector2D{X1: x1 / scaleDiv, Y1: y1 / scaleDiv, X2: x2 / scaleDiv, Y2: y2 / scaleDiv})
	}
}
