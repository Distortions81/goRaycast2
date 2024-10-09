package main

import "github.com/chewxy/math32"

// Dot product of two 2D vectors
func dotXY(v1, v2 pos32) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Subtract two vectors
func subXY(v1, v2 pos32) pos32 {
	return pos32{v1.X - v2.X, v1.Y - v2.Y}
}

// Subtract two vectors
func addXY(v1, v2 pos32) pos32 {
	return pos32{v1.X + v2.X, v1.Y + v2.Y}
}

// Scale a vector by a scalar
func scaleXY(v pos32, scalar float32) pos32 {
	return pos32{v.X * scalar, v.Y * scalar}
}

// Normalize a vector
func normalizeXY(v pos32) pos32 {
	magnitude := math32.Sqrt(v.X*v.X + v.Y*v.Y)
	if magnitude == 0 {
		return pos32{0, 0}
	}
	return pos32{v.X / magnitude, v.Y / magnitude}
}

func movementDirection(wall Line32) pos32 {
	return pos32{
		X: wall.X2 - wall.X1,
		Y: wall.Y2 - wall.Y1,
	}
}

func rayIntersectsSegment(rayDir pos32, wall Line32) (float32, pos32, bool) {
	// Using line intersection formula
	x1, y1, x2, y2 := wall.X1, wall.Y1, wall.X2, wall.Y2

	denom := (x1-x2)*(player.pos.Y+rayDir.Y-player.pos.Y) - (y1-y2)*(player.pos.X+rayDir.X-player.pos.X)
	if denom == 0 {
		return 0, pos32{}, false // Parallel lines
	}

	// t and u parameters for intersection formula
	t := ((x1-player.pos.X)*(player.pos.Y+rayDir.Y-player.pos.Y) - (y1-player.pos.Y)*(player.pos.X+rayDir.X-player.pos.X)) / denom
	u := -((x1-x2)*(y1-player.pos.Y) - (y1-y2)*(x1-player.pos.X)) / denom

	// If t and u are valid, we have an intersection
	if t >= 0 && t <= 1 && u > 0 {
		// Calculate the intersection point using t
		intersection := pos32{
			X: x1 + t*(x2-x1),
			Y: y1 + t*(y2-y1),
		}
		return u, intersection, true
	}

	return 0, pos32{}, false
}

func BoxToVectors(x, y, width, height float32) []Line32 {
	// Define the four corners of the box
	topLeft := Line32{X1: x, Y1: y, X2: x + width, Y2: y}                       // Top edge
	topRight := Line32{X1: x + width, Y1: y, X2: x + width, Y2: y + height}     // Right edge
	bottomRight := Line32{X1: x + width, Y1: y + height, X2: x, Y2: y + height} // Bottom edge
	bottomLeft := Line32{X1: x, Y1: y + height, X2: x, Y2: y}                   // Left edge

	// Return the four edges of the box
	return []Line32{topLeft, topRight, bottomRight, bottomLeft}
}
