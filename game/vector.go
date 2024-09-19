package main

import "math"

// Dot product of two 2D vectors
func dotXY(v1, v2 XY64) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Subtract two vectors
func subXY(v1, v2 XY64) XY64 {
	return XY64{v1.X - v2.X, v1.Y - v2.Y}
}

// Subtract two vectors
func addXY(v1, v2 XY64) XY64 {
	return XY64{v1.X + v2.X, v1.Y + v2.Y}
}

// Scale a vector by a scalar
func scaleXY(v XY64, scalar float64) XY64 {
	return XY64{v.X * scalar, v.Y * scalar}
}

// Normalize a vector
func normalizeXY(v XY64) XY64 {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if magnitude == 0 {
		return XY64{0, 0}
	}
	return XY64{v.X / magnitude, v.Y / magnitude}
}

func movementDirection(wall Vec64) XY64 {
	return XY64{
		X: wall.X2 - wall.X1,
		Y: wall.Y2 - wall.Y1,
	}
}

func rayIntersectsSegment(p, rayDir XY64, wall Vec64) (float64, XY64, bool) {
	// Using line intersection formula
	x1, y1, x2, y2 := wall.X1, wall.Y1, wall.X2, wall.Y2

	denom := (x1-x2)*(p.Y+rayDir.Y-p.Y) - (y1-y2)*(p.X+rayDir.X-p.X)
	if denom == 0 {
		return 0, XY64{}, false // Parallel lines
	}

	// t and u parameters for intersection formula
	t := ((x1-p.X)*(p.Y+rayDir.Y-p.Y) - (y1-p.Y)*(p.X+rayDir.X-p.X)) / denom
	u := -((x1-x2)*(y1-p.Y) - (y1-y2)*(x1-p.X)) / denom

	// If t and u are valid, we have an intersection
	if t >= 0 && t <= 1 && u > 0 {
		// Calculate the intersection point using t
		intersection := XY64{
			X: x1 + t*(x2-x1),
			Y: y1 + t*(y2-y1),
		}
		return u, intersection, true
	}

	return 0, XY64{}, false
}

func BoxToVectors(x, y, width, height float64) []Vec64 {
	// Define the four corners of the box
	topLeft := Vec64{X1: x, Y1: y, X2: x + width, Y2: y}                       // Top edge
	topRight := Vec64{X1: x + width, Y1: y, X2: x + width, Y2: y + height}     // Right edge
	bottomRight := Vec64{X1: x + width, Y1: y + height, X2: x, Y2: y + height} // Bottom edge
	bottomLeft := Vec64{X1: x, Y1: y + height, X2: x, Y2: y}                   // Left edge

	// Return the four edges of the box
	return []Vec64{topLeft, topRight, bottomRight, bottomLeft}
}
