package main

import "github.com/hajimehoshi/ebiten/v2"

type BSPNode struct {
	wall   Vec64    // The wall that splits the space
	front  *BSPNode // The front subspace
	back   *BSPNode // The back subspace
	isLeaf bool     // Whether this node is a leaf node
	walls  []Vec64  // Walls in the node (for leaf nodes)
}

// Function to calculate which side of the wall the player is on
func pointSide(p XY64, wall Vec64) float64 {
	return (wall.X2-wall.X1)*(p.Y-wall.Y1) - (wall.Y2-wall.Y1)*(p.X-wall.X1)
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

// Traverse the BSP tree and render the walls in correct order
func renderBSPTree(node *BSPNode, playerPos XY64, screen *ebiten.Image) {
	if node == nil {
		return
	}

	// Determine which side of the partition wall the player is on
	playerSide := pointSide(playerPos, node.wall)

	// Back-to-front traversal to ensure proper occlusion (Painter's Algorithm)
	if playerSide > 0 {
		// Player is in front of the wall; traverse the back subspace first
		renderBSPTree(node.back, playerPos, screen)
		renderWall(node.wall, screen)
		renderBSPTree(node.front, playerPos, screen)
	} else {
		// Player is behind the wall; traverse the front subspace first
		renderBSPTree(node.front, playerPos, screen)
		renderWall(node.wall, screen)
		renderBSPTree(node.back, playerPos, screen)
	}
}
