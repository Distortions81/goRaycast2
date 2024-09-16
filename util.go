package main

import (
	"image/color"
	"math"
)

// HSVtoRGB converts HSV values to RGB
func HSVtoRGB(h, s, v float64) color.NRGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r1, g1, b1 float64

	if h >= 0 && h < 60 {
		r1, g1, b1 = c, x, 0
	} else if h >= 60 && h < 120 {
		r1, g1, b1 = x, c, 0
	} else if h >= 120 && h < 180 {
		r1, g1, b1 = 0, c, x
	} else if h >= 180 && h < 240 {
		r1, g1, b1 = 0, x, c
	} else if h >= 240 && h < 300 {
		r1, g1, b1 = x, 0, c
	} else if h >= 300 && h < 360 {
		r1, g1, b1 = c, 0, x
	}

	// Convert to RGB by adding m and scaling to the range of 0-255
	r := (r1 + m) * 255
	g := (g1 + m) * 255
	b := (b1 + m) * 255

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 1}
}

func rayIntersectsSegment(px, py, rayDirX, rayDirY float64, wall Vector) (float64, bool) {
	// Using line intersection formula
	x1, y1, x2, y2 := wall.X1, wall.Y1, wall.X2, wall.Y2

	denom := (x1-x2)*(py+rayDirY-py) - (y1-y2)*(px+rayDirX-px)
	if denom == 0 {
		return 0, false // Parallel lines
	}

	// t and u parameters for intersection formula
	t := ((x1-px)*(py+rayDirY-py) - (y1-py)*(px+rayDirX-px)) / denom
	u := -((x1-x2)*(y1-py) - (y1-y2)*(x1-px)) / denom

	// If t and u are valid, we have an intersection
	if t >= 0 && t <= 1 && u > 0 {
		return u, true
	}

	return 0, false
}

func BoxToVectors(x, y, width, height float64) []Vector {
	// Define the four corners of the box
	topLeft := Vector{X1: x, Y1: y, X2: x + width, Y2: y}                       // Top edge
	topRight := Vector{X1: x + width, Y1: y, X2: x + width, Y2: y + height}     // Right edge
	bottomRight := Vector{X1: x + width, Y1: y + height, X2: x, Y2: y + height} // Bottom edge
	bottomLeft := Vector{X1: x, Y1: y + height, X2: x, Y2: y}                   // Left edge

	// Return the four edges of the box
	return []Vector{topLeft, topRight, bottomRight, bottomLeft}
}
