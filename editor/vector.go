package main

import "math"

func distance(p1, p2 XY) float64 {
	return math.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}

// snapPos snaps a new position to the nearest existing position within a threshold
func snapPos(newPos XY, existingPositions []Vector2D, threshold float64) XY {
	minDistance := threshold // Initialize with threshold to ensure snapping only within the threshold

	for _, pos := range existingPositions {

		apos := XY{X: pos.X1, Y: pos.Y1}
		bpos := XY{X: pos.X2, Y: pos.Y2}

		dist := distance(newPos, apos)
		if dist < minDistance {
			minDistance = dist
			newPos = apos
		}

		dist = distance(newPos, bpos)
		if dist < minDistance {
			minDistance = dist
			newPos = bpos
		}
	}

	return SnapToGrid(newPos, 10, 3)
}

func SnapToGrid(pos XY, gridSize, threshold float64) XY {
	snapX := math.Round(pos.X/gridSize) * gridSize
	snapY := math.Round(pos.Y/gridSize) * gridSize

	// Check if the coordinates are within the snapping threshold
	if math.Abs(pos.X-snapX) <= threshold {
		pos.X = snapX
	}
	if math.Abs(pos.Y-snapY) <= threshold {
		pos.Y = snapY
	}

	return XY{X: pos.X, Y: pos.Y}
}
