package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	miniMapSize   = 100
	miniMapRadius = 30
)

// Render a clipped wall on the minimap, ensuring it stays within the minimap radius
func renderClippedWallOnMinimap(wall Vec64, playerX, playerY float32, screen *ebiten.Image) {
	// Translate wall coordinates relative to player position
	dx1 := float32(wall.X1) - float32(player.pos.X)
	dy1 := float32(wall.Y1) - float32(player.pos.Y)
	dx2 := float32(wall.X2) - float32(player.pos.X)
	dy2 := float32(wall.Y2) - float32(player.pos.Y)

	// Scale wall coordinates to minimap size based on miniMapRadius
	x1 := playerX + dx1*(miniMapSize/miniMapRadius)
	y1 := playerY + dy1*(miniMapSize/miniMapRadius)
	x2 := playerX + dx2*(miniMapSize/miniMapRadius)
	y2 := playerY + dy2*(miniMapSize/miniMapRadius)

	// Clip the wall to ensure it's within the minimap radius
	if !clipLineToCircle(&x1, &y1, &x2, &y2, playerX, playerY, miniMapRadius*(miniMapSize/miniMapRadius)) {
		return // Skip if the wall is entirely outside the minimap radius
	}

	// Draw the clipped wall line on the minimap
	vector.StrokeLine(screen, x1, y1, x2, y2, 1, colornames.Teal, false)
}

// Traverse BSP and render walls within minimap radius
func traverseBSPForMinimap(node *BSPNode, playerX, playerY float32, screen *ebiten.Image) {
	if node == nil {
		return
	}

	// Calculate the distance to both wall endpoints using the player's position
	distToWall1 := calculateDistance(float32(node.wall.X1), float32(node.wall.Y1), float32(player.pos.X), float32(player.pos.Y))
	distToWall2 := calculateDistance(float32(node.wall.X2), float32(node.wall.Y2), float32(player.pos.X), float32(player.pos.Y))

	// Check if either endpoint is within the minimap radius, or if the wall intersects the radius
	if distToWall1 <= miniMapRadius || distToWall2 <= miniMapRadius || wallIntersectsMinimap(node.wall) {
		// Clip the wall if needed and render the portion inside the minimap radius
		renderClippedWallOnMinimap(node.wall, playerX, playerY, screen)
	}

	// Recursively traverse the front and back subtrees
	traverseBSPForMinimap(node.front, playerX, playerY, screen)
	traverseBSPForMinimap(node.back, playerX, playerY, screen)
}

// Render the minimap, ensuring it fits fully on the screen
func renderMinimap(screen *ebiten.Image) {
	// Calculate the center of the minimap, based on screenWidth and miniMapSize
	miniMapCenterX := float32(miniMapSize + 10)
	miniMapCenterY := float32(miniMapSize + 10)

	// Calculate the top-left corner of the minimap
	miniMapTopLeftX := miniMapCenterX - miniMapSize/2
	miniMapTopLeftY := miniMapCenterY - miniMapSize/2

	// Precalculate player's position within the minimap (center of the minimap)
	playerX := miniMapTopLeftX + miniMapSize/2
	playerY := miniMapTopLeftY + miniMapSize/2

	// Traverse the BSP tree and render walls within the minimap's clipping radius
	traverseBSPForMinimap(bspData, playerX, playerY, screen)

	// Draw the player as a circle in the center of the minimap
	vector.DrawFilledCircle(screen, playerX, playerY, 5, colornames.Yellow, false)

	// Optionally, draw the player's facing direction on the minimap
	facingX := playerX - float32(math.Cos(player.angle))*10
	facingY := playerY - float32(math.Sin(player.angle))*10
	vector.StrokeLine(screen, playerX, playerY, facingX, facingY, 2, colornames.Red, false)
}

// Calculate the 2D Euclidean distance between two points
func calculateDistance(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))))
}

// Check if a wall intersects the minimap radius
func wallIntersectsMinimap(wall Vec64) bool {
	// Calculate distances from both endpoints to the player position (minimap center)
	dist1 := calculateDistance(float32(wall.X1), float32(wall.Y1), float32(player.pos.X), float32(player.pos.Y))
	dist2 := calculateDistance(float32(wall.X2), float32(wall.Y2), float32(player.pos.X), float32(player.pos.Y))

	// Check if one endpoint is inside the minimap radius and the other is outside
	return (dist1 <= miniMapRadius && dist2 > miniMapRadius) || (dist2 <= miniMapRadius && dist1 > miniMapRadius)
}

// Clip a line to a circle and modify the endpoints to keep them within the circle
func clipLineToCircle(x1, y1, x2, y2 *float32, cx, cy, r float32) bool {
	// Vector from the center of the circle to the first point
	dx1 := *x1 - cx
	dy1 := *y1 - cy
	dist1 := float32(math.Sqrt(float64(dx1*dx1 + dy1*dy1)))

	// Vector from the center of the circle to the second point
	dx2 := *x2 - cx
	dy2 := *y2 - cy
	dist2 := float32(math.Sqrt(float64(dx2*dx2 + dy2*dy2)))

	// If both points are inside the circle, no need to clip
	if dist1 <= r && dist2 <= r {
		return true
	}

	// If both points are outside the circle, return false (skip this line)
	if dist1 > r && dist2 > r {
		return false
	}

	// Normalize the direction vector from point 1 to point 2
	dx := *x2 - *x1
	dy := *y2 - *y1
	len := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	dx /= len
	dy /= len

	// Clip point 1 if it's outside the circle
	if dist1 > r {
		*x1, *y1 = intersectWithCircle(*x1, *y1, dx, dy, cx, cy, r)
	}

	// Clip point 2 if it's outside the circle
	if dist2 > r {
		*x2, *y2 = intersectWithCircle(*x2, *y2, -dx, -dy, cx, cy, r)
	}

	return true
}

// Calculate the intersection of a line with a circle
func intersectWithCircle(x, y, dx, dy, cx, cy, r float32) (float32, float32) {
	// Translate the line's starting point so the circle is centered at the origin
	x -= cx
	y -= cy

	// Quadratic equation coefficients for line-circle intersection
	a := dx*dx + dy*dy
	b := 2 * (x*dx + y*dy)
	c := x*x + y*y - r*r

	// Solve the quadratic equation using the discriminant
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		// No intersection; return original point
		return x + cx, y + cy
	}

	// Calculate the intersection point
	t := (-b - float32(math.Sqrt(float64(discriminant)))) / (2 * a)
	return x + t*dx + cx, y + t*dy + cy
}