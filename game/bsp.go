package main

import (
	"image"
	"sync"

	"github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"
)

type BSPNode struct {
	wall   line32   // The wall that splits the space
	front  *BSPNode // The front subspace
	back   *BSPNode // The back subspace
	isLeaf bool     // Whether this node is a leaf node
	walls  []line32 // Walls in the node (for leaf nodes)
}

type renderData struct {
	textureX, lineHeight, drawStart, x int
	textureY, valFloat                 float32
}

// Build a BSP tree from a list of walls
func buildBSPTree(walls []line32) *BSPNode {
	if len(walls) == 0 {
		return nil
	}

	// Pick the first wall as the partitioning wall (you can optimize this choice)
	partitionWall := walls[0]

	// Initialize lists for front and back walls
	var frontWalls, backWalls []line32

	// Classify the remaining walls as either front or back of the partition wall
	for i := 1; i < len(walls); i++ {
		wall := walls[i]
		frontCount := 0
		backCount := 0

		// Check the endpoints of the wall
		if pointSide(pos32{wall.X1, wall.Y1}, partitionWall) > 0 {
			frontCount++
		} else {
			backCount++
		}
		if pointSide(pos32{wall.X2, wall.Y2}, partitionWall) > 0 {
			frontCount++
		} else {
			backCount++
		}

		// Add the wall to the appropriate list
		if frontCount == 2 {
			frontWalls = append(frontWalls, wall)
			//fmt.Printf("F: %v, ", wall)
		} else if backCount == 2 {
			backWalls = append(backWalls, wall)
			//fmt.Printf("B: %v, ", wall)
		} else {
			//Split wallls
			backWalls = append(backWalls, wall)
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
func renderBSPTree(node *BSPNode, nearestDist *float32, closestWall *line32) {
	if node == nil {
		return
	}

	// Determine which side of the partition wall the player is on
	playerSide := pointSide(player.pos, node.wall)

	// Back-to-front traversal to ensure proper occlusion (Painter's Algorithm)
	if playerSide > 0 {
		// Player is in front of the wall; traverse the back subspace first
		renderBSPTree(node.back, nearestDist, closestWall)
		checkAndRenderWall(node.wall, nearestDist, closestWall)
		renderBSPTree(node.front, nearestDist, closestWall)
	} else {
		// Player is behind the wall; traverse the front subspace first
		renderBSPTree(node.front, nearestDist, closestWall)
		checkAndRenderWall(node.wall, nearestDist, closestWall)
		renderBSPTree(node.back, nearestDist, closestWall)
	}
}

// Traverse the BSP tree and find the closest wall
func findClosestWall(node *BSPNode, playerPos pos32, nearestDist *float32, closestWall *line32) {
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
func pointSide(p pos32, wall line32) float32 {
	return (wall.X2-wall.X1)*(p.Y-wall.Y1) - (wall.Y2-wall.Y1)*(p.X-wall.X1)
}

// Check the distance to the current wall and update the closest wall if it's nearer
func checkAndTrackWall(wall line32, playerPos pos32, nearestDist *float32, closestWall *line32) {
	// Calculate distance from player to wall
	dist := distanceToWall(wall, playerPos)

	// If this wall is closer than the previous nearest, update the nearest wall
	if dist < *nearestDist {
		*nearestDist = dist
		*closestWall = wall
	}
}

// Calculate the distance from the player to a wall
func distanceToWall(wall line32, playerPos pos32) float32 {
	// This function should calculate the perpendicular distance from the player to the wall
	// For now, assuming a simple Euclidean distance to one endpoint as a placeholder
	// You can improve this with proper perpendicular distance calculation based on the player's position.
	return math32.Sqrt(math32.Pow(playerPos.X-wall.X1, 2) + math32.Pow(playerPos.Y-wall.Y1, 2))
}

// Check the distance to the current wall and update the closest wall if it's nearer
func checkAndRenderWall(wall line32, nearestDist *float32, closestWall *line32) {
	// Calculate distance from player to wall
	dist := distanceToWall(wall, player.pos)

	// If this wall is closer than the previous nearest, update the nearest wall
	if dist < *nearestDist {
		*nearestDist = dist
		*closestWall = wall
	}
}

// Traverse the BSP tree and find the closest wall for the current ray
func findClosestWallForRay(node *BSPNode, rayDir pos32, nearestDist *float32, closestWall *line32, hitPos *pos32) {
	if node == nil {
		return
	}

	raySide := pointSide(player.pos, node.wall)

	if raySide > 0 {
		findClosestWallForRay(node.back, rayDir, nearestDist, closestWall, hitPos)
		checkAndTrackWallForRay(node.wall, rayDir, nearestDist, closestWall, hitPos)
		findClosestWallForRay(node.front, rayDir, nearestDist, closestWall, hitPos)
	} else {
		findClosestWallForRay(node.front, rayDir, nearestDist, closestWall, hitPos)
		checkAndTrackWallForRay(node.wall, rayDir, nearestDist, closestWall, hitPos)
		findClosestWallForRay(node.back, rayDir, nearestDist, closestWall, hitPos)
	}
}

// Check if the ray intersects with the current wall, and track the closest wall if it does
func checkAndTrackWallForRay(wall line32, rayDir pos32, nearestDist *float32, closestWall *line32, hitPos *pos32) {
	// Ray-wall intersection logic
	if dist, hPos, hit := rayIntersectsSegment(rayDir, wall); hit {
		// If this wall is closer than the previous nearest, update the nearest wall
		if dist < *nearestDist {
			*nearestDist = dist
			*closestWall = wall
			*hitPos = hPos // Store the hit position
		}
	}
}

var wg sync.WaitGroup

func renderScene(screen *ebiten.Image) {

	for x := 0; x < screenWidth; x += workSize {
		wg.Add(1)
		go func(start int) {
			end := min(start+workSize, screenWidth-1)
			for col := start; col < end; col++ {
				cameraX := 2*float32(col)/float32(screenWidth) - 1
				rayDir := angleToXY(player.angle+math32.Atan(cameraX), 1)

				var nearestDist float32 = math32.MaxFloat32
				var wall line32
				var hitPos pos32

				findClosestWallForRay(bspData, rayDir, &nearestDist, &wall, &hitPos)
				correctedDist := nearestDist * math32.Cos(math32.Atan(cameraX))

				lineHeight := int(float32(screenHeight) / correctedDist)
				drawStart := -lineHeight/2 + screenHeight/2
				if drawStart < 0 {
					drawStart = 0
				}
				drawEnd := lineHeight/2 + screenHeight/2
				if drawEnd >= screenHeight {
					drawEnd = screenHeight - 1
				}

				// Calculate the direction vector for the wall
				wallDirX := wall.X2 - wall.X1
				wallDirY := wall.Y2 - wall.Y1
				wallLength := math32.Sqrt(wallDirX*wallDirX + wallDirY*wallDirY)

				// Normalize the direction vector
				wallDirX /= wallLength
				wallDirY /= wallLength

				// Calculate the hit position along the wall
				dx := hitPos.X - wall.X1
				dy := hitPos.Y - wall.Y1
				wallHitPosition := (dx*wallDirX + dy*wallDirY)

				// Calculate texture X based on the fixed texture repeat distance
				wallHitPosition = math32.Mod(wallHitPosition, textureRepeatDistance)
				textureX := int((wallHitPosition/textureRepeatDistance)*float32(textureWidth)) % textureWidth
				if textureX < 0 {
					textureX += textureWidth
				}

				// Texture Y scaling and clipping
				var textureStep float32 = float32(textureHeight) / float32(lineHeight)
				var textureY float32 = 0.0

				// If the wall height is larger than the screen, adjust textureY and clip the texture
				if lineHeight > screenHeight {
					textureY = float32(lineHeight-screenHeight) / 2.0 * textureStep
					drawStart = 0 // Clamp drawStart to 0 (top of screen)
				}

				// Calculate the lighting/shading factor
				valFloat := applyFalloff(nearestDist, lightIntensity, float32(wallColor.R+wallColor.G+wallColor.B)/765.0/3.0)

				rayList[col] = renderData{
					textureX: textureX, lineHeight: lineHeight, x: col,
					drawStart: drawStart, textureY: textureY, valFloat: valFloat}
			}
			wg.Done()
		}(x)
	}
	wg.Wait()

	renderWallSlice(screen)
}

func renderWallSlice(screen *ebiten.Image) {

	for _, data := range rayList {
		// Create a sub-image of the texture slice to draw (from textureX to textureX + 1)
		srcRect := image.Rect(data.textureX, int(data.textureY), data.textureX+1, textureHeight)

		// Apply shading and draw the texture slice
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
		op.GeoM.Scale(1, float64(data.lineHeight)/float64(textureHeight)) // Scale texture to line height
		op.GeoM.Translate(float64(data.x), float64(data.drawStart))       // Position the texture slice
		op.ColorScale.Scale(data.valFloat, data.valFloat, data.valFloat, 1)

		screen.DrawImage(wallImg.SubImage(srcRect).(*ebiten.Image), op)
	}
}
