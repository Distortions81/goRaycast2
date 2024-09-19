package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

// Optimized Minimap Rendering with BSP Traversal and Proper Endpoint Clipping and Intersection Check
func renderMinimap(screen *ebiten.Image) {
	// Calculate the center of the minimap on the screen
	miniMapCenterX := float32(screenWidth) - miniMapSize/2 - 10 // Example: place the minimap on the right side of the screen
	miniMapCenterY := float32(miniMapSize/2 + 10)               // Cast to float32

	// Precalculate player position on the minimap (the player is at the center)
	playerX := miniMapCenterX
	playerY := miniMapCenterY

	// Traverse the BSP tree and render walls within the minimap's clipping radius
	traverseBSPForMinimap(bspData, playerX, playerY, miniMapSize, miniMapRadius, screen)

	// Draw the player as a circle in the center of the minimap
	vector.DrawFilledCircle(screen, playerX, playerY, 5, colornames.Yellow, false)

	// Optionally, draw the player's facing direction on the minimap
	facingX := playerX + float32(math.Cos(player.angle))*10
	facingY := playerY + float32(math.Sin(player.angle))*10
	vector.StrokeLine(screen, playerX, playerY, facingX, facingY, 2, colornames.Red, false)
}

// Traverse BSP and render walls within minimap radius by checking both endpoints of the wall
func traverseBSPForMinimap(node *BSPNode, playerX, playerY, miniMapSize, miniMapRadius float32, screen *ebiten.Image) {
	if node == nil {
		return
	}

	// Calculate the distance to both wall endpoints using the player's position
	distToWall1 := calculateDistance(float32(node.wall.X1), float32(node.wall.Y1), float32(player.pos.X), float32(player.pos.Y))
	distToWall2 := calculateDistance(float32(node.wall.X2), float32(node.wall.Y2), float32(player.pos.X), float32(player.pos.Y))

	// Check if either endpoint is within the minimap radius, or if the wall intersects the radius
	if distToWall1 <= miniMapRadius || distToWall2 <= miniMapRadius || wallIntersectsMinimap(node.wall, miniMapRadius) {
		// Render the wall if it’s within the minimap radius or intersects it
		renderWallOnMinimap(node.wall, playerX, playerY, miniMapSize, miniMapRadius, screen)
	}

	// Recursively traverse the front and back subtrees
	traverseBSPForMinimap(node.front, playerX, playerY, miniMapSize, miniMapRadius, screen)
	traverseBSPForMinimap(node.back, playerX, playerY, miniMapSize, miniMapRadius, screen)
}

// Calculate the 2D Euclidean distance between two points
func calculateDistance(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))))
}

// Check if a wall intersects the minimap radius
func wallIntersectsMinimap(wall Vec64, miniMapRadius float32) bool {
	// Basic intersection check: does the wall cross the circular minimap boundary?
	// We’ll use a simple bounding box check or a more precise geometric check to see if
	// the line from (X1, Y1) to (X2, Y2) crosses the circular radius.

	// For simplicity, we use a bounding box approximation first.
	// Calculate the distances from both endpoints to the origin (player's position at the minimap center)
	dist1 := calculateDistance(float32(wall.X1), float32(wall.Y1), float32(player.pos.X), float32(player.pos.Y))
	dist2 := calculateDistance(float32(wall.X2), float32(wall.Y2), float32(player.pos.X), float32(player.pos.Y))

	// If one endpoint is inside and the other is outside, we assume an intersection
	return (dist1 <= miniMapRadius && dist2 > miniMapRadius) || (dist2 <= miniMapRadius && dist1 > miniMapRadius)
}

// Render a single wall on the minimap
func renderWallOnMinimap(wall Vec64, playerX, playerY, miniMapSize, miniMapRadius float32, screen *ebiten.Image) {
	// Translate wall coordinates relative to player position
	dx1 := float32(wall.X1) - float32(player.pos.X)
	dy1 := float32(wall.Y1) - float32(player.pos.Y)
	dx2 := float32(wall.X2) - float32(player.pos.X)
	dy2 := float32(wall.Y2) - float32(player.pos.Y)

	// Scale wall coordinates to minimap size
	x1 := playerX + dx1*(miniMapSize/miniMapRadius)
	y1 := playerY + dy1*(miniMapSize/miniMapRadius)
	x2 := playerX + dx2*(miniMapSize/miniMapRadius)
	y2 := playerY + dy2*(miniMapSize/miniMapRadius)

	// Draw the wall line on the minimap
	vector.StrokeLine(screen, x1, y1, x2, y2, 1, colornames.Teal, false)
}
