package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type BSPNode struct {
	wall   Vec64    // The wall that splits the space
	front  *BSPNode // The front subspace
	back   *BSPNode // The back subspace
	isLeaf bool     // Whether this node is a leaf node
	walls  []Vec64  // Walls in the node (for leaf nodes)
}

// Build a BSP tree from a list of walls
func buildBSPTree(walls []Vec64) *BSPNode {
	if len(walls) == 0 {
		return nil
	}

	// Pick the first wall as the partitioning wall (you can optimize this choice)
	partitionWall := walls[0]

	// Initialize lists for front and back walls
	var frontWalls, backWalls []Vec64

	// Classify the remaining walls as either front or back of the partition wall
	for i := 1; i < len(walls); i++ {
		wall := walls[i]
		frontCount := 0
		backCount := 0

		// Check the endpoints of the wall
		if pointSide(XY64{wall.X1, wall.Y1}, partitionWall) > 0 {
			frontCount++
		} else {
			backCount++
		}
		if pointSide(XY64{wall.X2, wall.Y2}, partitionWall) > 0 {
			frontCount++
		} else {
			backCount++
		}

		// Add the wall to the appropriate list
		if frontCount == 2 {
			frontWalls = append(frontWalls, wall)
		} else if backCount == 2 {
			backWalls = append(backWalls, wall)
		} else {
			// Split wall logic can go here for walls that straddle both sides (if needed)
		}
	}

	// Recursively build the BSP tree
	return &BSPNode{
		wall:   partitionWall,
		front:  buildBSPTree(frontWalls),
		back:   buildBSPTree(backWalls),
		isLeaf: false,
	}
}

// Traverse the BSP tree and render the closest wall in correct order
func renderBSPTree(node *BSPNode, playerPos XY64, screen *ebiten.Image, nearestDist *float64, closestWall *Vec64) {
	if node == nil {
		return
	}

	// Determine which side of the partition wall the player is on
	playerSide := pointSide(playerPos, node.wall)

	// Back-to-front traversal to ensure proper occlusion (Painter's Algorithm)
	if playerSide > 0 {
		// Player is in front of the wall; traverse the back subspace first
		renderBSPTree(node.back, playerPos, screen, nearestDist, closestWall)
		checkAndRenderWall(node.wall, playerPos, nearestDist, closestWall)
		renderBSPTree(node.front, playerPos, screen, nearestDist, closestWall)
	} else {
		// Player is behind the wall; traverse the front subspace first
		renderBSPTree(node.front, playerPos, screen, nearestDist, closestWall)
		checkAndRenderWall(node.wall, playerPos, nearestDist, closestWall)
		renderBSPTree(node.back, playerPos, screen, nearestDist, closestWall)
	}
}

// Traverse the BSP tree and find the closest wall
func findClosestWall(node *BSPNode, playerPos XY64, nearestDist *float64, closestWall *Vec64) {
	if node == nil {
		return
	}

	// Determine which side of the partition wall the player is on
	playerSide := pointSide(playerPos, node.wall)

	// Back-to-front traversal to ensure proper occlusion (Painter's Algorithm)
	if playerSide > 0 {
		// Player is in front of the wall; traverse the back subspace first
		findClosestWall(node.back, playerPos, nearestDist, closestWall)
		checkAndTrackWall(node.wall, playerPos, nearestDist, closestWall)
		findClosestWall(node.front, playerPos, nearestDist, closestWall)
	} else {
		// Player is behind the wall; traverse the front subspace first
		findClosestWall(node.front, playerPos, nearestDist, closestWall)
		checkAndTrackWall(node.wall, playerPos, nearestDist, closestWall)
		findClosestWall(node.back, playerPos, nearestDist, closestWall)
	}
}

// Function to calculate which side of the wall the player is on
func pointSide(p XY64, wall Vec64) float64 {
	return (wall.X2-wall.X1)*(p.Y-wall.Y1) - (wall.Y2-wall.Y1)*(p.X-wall.X1)
}

// Check the distance to the current wall and update the closest wall if it's nearer
func checkAndTrackWall(wall Vec64, playerPos XY64, nearestDist *float64, closestWall *Vec64) {
	// Calculate distance from player to wall
	dist := distanceToWall(wall, playerPos)

	// If this wall is closer than the previous nearest, update the nearest wall
	if dist < *nearestDist {
		*nearestDist = dist
		*closestWall = wall
	}
}

// Calculate the distance from the player to a wall
func distanceToWall(wall Vec64, playerPos XY64) float64 {
	// This function should calculate the perpendicular distance from the player to the wall
	// For now, assuming a simple Euclidean distance to one endpoint as a placeholder
	// You can improve this with proper perpendicular distance calculation based on the player's position.
	return math.Sqrt(math.Pow(playerPos.X-wall.X1, 2) + math.Pow(playerPos.Y-wall.Y1, 2))
}

// Check the distance to the current wall and update the closest wall if it's nearer
func checkAndRenderWall(wall Vec64, playerPos XY64, nearestDist *float64, closestWall *Vec64) {
	// Calculate distance from player to wall
	dist := distanceToWall(wall, playerPos)

	// If this wall is closer than the previous nearest, update the nearest wall
	if dist < *nearestDist {
		*nearestDist = dist
		*closestWall = wall
	}
}

// Traverse the BSP tree and find the closest wall for the current ray
func findClosestWallForRay(node *BSPNode, playerPos XY64, rayDir XY64, nearestDist *float64, closestWall *Vec64, hitPos *XY64) {
	if node == nil {
		return
	}

	// Determine which side of the partition wall the player is on
	raySide := pointSide(playerPos, node.wall)

	// Back-to-front traversal to ensure proper occlusion (Painter's Algorithm)
	if raySide > 0 {
		findClosestWallForRay(node.back, playerPos, rayDir, nearestDist, closestWall, hitPos)
		checkAndTrackWallForRay(node.wall, playerPos, rayDir, nearestDist, closestWall, hitPos)
		findClosestWallForRay(node.front, playerPos, rayDir, nearestDist, closestWall, hitPos)
	} else {
		findClosestWallForRay(node.front, playerPos, rayDir, nearestDist, closestWall, hitPos)
		checkAndTrackWallForRay(node.wall, playerPos, rayDir, nearestDist, closestWall, hitPos)
		findClosestWallForRay(node.back, playerPos, rayDir, nearestDist, closestWall, hitPos)
	}
}

// Check if the ray intersects with the current wall, and track the closest wall if it does
func checkAndTrackWallForRay(wall Vec64, playerPos XY64, rayDir XY64, nearestDist *float64, closestWall *Vec64, hitPos *XY64) {
	// Ray-wall intersection logic
	if dist, hPos, hit := rayIntersectsSegment(playerPos, rayDir, wall); hit {
		// If this wall is closer than the previous nearest, update the nearest wall
		if dist < *nearestDist {
			*nearestDist = dist
			*closestWall = wall
			*hitPos = hPos // Store the hit position
		}
	}
}

// Main rendering loop
func renderScene(bspRoot *BSPNode, playerPos XY64, playerAngle float64, screen *ebiten.Image) {
	for x := 0; x < screenWidth; x++ {
		// 1. Map screen X to camera plane (-1 to 1)
		cameraX := 2*float64(x)/float64(screenWidth) - 1

		// 2. Calculate the ray direction based on the player's current angle and camera plane offset
		rayDir := angleToXY(playerAngle+math.Atan(cameraX), 1)

		// Variables to track the nearest wall and intersection for this ray
		nearestDist := math.MaxFloat64
		var closestWall Vec64
		var hitPos XY64

		// Traverse the BSP tree and find the closest wall for this ray
		findClosestWallForRay(bspRoot, playerPos, rayDir, &nearestDist, &closestWall, &hitPos)

		// Correct the fisheye effect by adjusting the distance based on the angle of the ray
		correctedDist := nearestDist * math.Cos(math.Atan(cameraX))

		// 3. If a closest wall is found, render the wall slice
		if nearestDist < math.MaxFloat64 {
			renderWallSlice(closestWall, hitPos, correctedDist, x, screen)
		}
	}
}

// Render a single slice of the wall for the current ray
func renderWallSlice(wall Vec64, hitPos XY64, nearestDist float64, x int, screen *ebiten.Image) {
	// Calculate the line height based on the corrected distance and project the wall slice
	lineHeight := int(float64(screenHeight) / nearestDist)
	drawStart := -lineHeight/2 + screenHeight/2
	if drawStart < 0 {
		drawStart = 0
	}
	drawEnd := lineHeight/2 + screenHeight/2
	if drawEnd >= screenHeight {
		drawEnd = screenHeight - 1
	}

	// Precompute texture width and height (cached outside the loop in the main render function)
	textureBounds := wallImg.Bounds()
	textureWidth := textureBounds.Dx()
	textureHeight := textureBounds.Dy()

	// Calculate the direction vector for the wall
	wallDirX := wall.X2 - wall.X1
	wallDirY := wall.Y2 - wall.Y1
	wallLength := math.Sqrt(wallDirX*wallDirX + wallDirY*wallDirY)

	// Normalize the direction vector
	wallDirX /= wallLength
	wallDirY /= wallLength

	// Use a fixed repetition distance for the texture (e.g., repeat every 1 unit of game space)
	textureRepeatDistance := 1.0

	// Calculate the hit position along the wall
	dx := hitPos.X - wall.X1
	dy := hitPos.Y - wall.Y1
	wallHitPosition := (dx*wallDirX + dy*wallDirY)

	// Calculate texture X based on the fixed texture repeat distance
	wallHitPosition = math.Mod(wallHitPosition, textureRepeatDistance)
	textureX := int((wallHitPosition/textureRepeatDistance)*float64(textureWidth)) % textureWidth
	if textureX < 0 {
		textureX += textureWidth
	}

	// Texture Y scaling and clipping
	textureStep := float64(textureHeight) / float64(lineHeight)
	textureY := 0.0

	// If the wall height is larger than the screen, adjust textureY and clip the texture
	if lineHeight > screenHeight {
		textureY = float64(lineHeight-screenHeight) / 2 * textureStep
		drawStart = 0 // Clamp drawStart to 0 (top of screen)
	}

	// Create a sub-image of the texture slice to draw (from textureX to textureX + 1)
	srcRect := image.Rect(textureX, int(textureY), textureX+1, textureHeight)
	textureSlice := wallImg.SubImage(srcRect).(*ebiten.Image)

	// Calculate the lighting/shading factor
	valFloat := applyFalloff(nearestDist, lightIntensity, float64(wallColor.R+wallColor.G+wallColor.B)/765.0/3.0)

	// Apply shading and draw the texture slice
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(1, float64(lineHeight)/float64(textureHeight)) // Scale texture to line height
	op.GeoM.Translate(float64(x), float64(drawStart))            // Position the texture slice
	op.ColorScale.Scale(valFloat, valFloat, valFloat, 1)

	// Draw the texture slice to the screen
	screen.DrawImage(textureSlice, op)
}
