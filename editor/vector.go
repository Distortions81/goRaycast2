package main

import "math"

func distance(p1, p2 XY) float64 {
	return math.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}

// snapPos snaps a new position to the nearest existing position within a threshold
func snapPos(newPos XY, existingPositions []Vector2D, threshold float64) XY {
	var snappedPosition XY = newPos
	minDistance := threshold // Initialize with threshold to ensure snapping only within the threshold

	for _, pos := range existingPositions {

		apos := XY{X: pos.X1, Y: pos.Y1}
		bpos := XY{X: pos.X2, Y: pos.Y2}

		dist := distance(newPos, apos)
		if dist < minDistance {
			minDistance = dist
			snappedPosition = apos
		}

		dist = distance(newPos, bpos)
		if dist < minDistance {
			minDistance = dist
			snappedPosition = bpos
		}
	}

	return snappedPosition
}
