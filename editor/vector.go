package main

import "github.com/chewxy/math32"

func distance(p1, p2 pos32) float32 {
	return math32.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}

// snapPos snaps a new position to the nearest existing position within a threshold
func snapPos(newPos pos32, existingPositions []line32, threshold float32) pos32 {
	minDistance := threshold // Initialize with threshold to ensure snapping only within the threshold

	for _, pos := range existingPositions {

		apos := pos32{X: pos.X1, Y: pos.Y1}
		bpos := pos32{X: pos.X2, Y: pos.Y2}

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

	return SnapToGrid(newPos, gridSize, gridSnapDist)
}

func SnapToGrid(pos pos32, gridSize, threshold float32) pos32 {
	snapX := math32.Round(pos.X/gridSize) * gridSize
	snapY := math32.Round(pos.Y/gridSize) * gridSize

	// Check if the coordinates are within the snapping threshold
	if math32.Abs(pos.X-snapX) <= threshold {
		pos.X = snapX
	}
	if math32.Abs(pos.Y-snapY) <= threshold {
		pos.Y = snapY
	}

	return pos32{X: pos.X, Y: pos.Y}
}
