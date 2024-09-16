package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const turnSpeed = 0.03
const moveSpeed = 0.1

func (g *Game) Update() error {
	oldPlayer := g.player
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.pos.X += g.player.dir.X * moveSpeed
		g.player.pos.Y += g.player.dir.Y * moveSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.pos.X -= g.player.dir.X * moveSpeed
		g.player.pos.Y -= g.player.dir.Y * moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldDirX := g.player.dir.X
		g.player.dir.X = g.player.dir.X*math.Cos(turnSpeed) - g.player.dir.Y*math.Sin(turnSpeed)
		g.player.dir.Y = oldDirX*math.Sin(turnSpeed) + g.player.dir.Y*math.Cos(turnSpeed)
		oldPlaneX := g.player.plane.X
		g.player.plane.X = g.player.plane.X*math.Cos(turnSpeed) - g.player.plane.Y*math.Sin(turnSpeed)
		g.player.plane.Y = oldPlaneX*math.Sin(turnSpeed) + g.player.plane.Y*math.Cos(turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldDirX := g.player.dir.X
		g.player.dir.X = g.player.dir.X*math.Cos(-turnSpeed) - g.player.dir.Y*math.Sin(-turnSpeed)
		g.player.dir.Y = oldDirX*math.Sin(-turnSpeed) + g.player.dir.Y*math.Cos(-turnSpeed)
		oldPlaneX := g.player.plane.X
		g.player.plane.X = g.player.plane.X*math.Cos(-turnSpeed) - g.player.plane.Y*math.Sin(-turnSpeed)
		g.player.plane.Y = oldPlaneX*math.Sin(-turnSpeed) + g.player.plane.Y*math.Cos(-turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.Key0) {
		fmt.Printf("posX: %v, posY: %v, dirX: %v, dirY: %v, planeX: %v, planeY: %v\n", g.player.pos.X, g.player.pos.Y, g.player.dir.X, g.player.dir.Y, g.player.plane.X, g.player.plane.Y)
	}

	for _, wall := range walls {
		if dist, hit := rayIntersectsSegment(g.player.pos, g.player.dir, wall); hit {
			if dist < 0.5 {
				move := subtract(oldPlayer.pos, g.player.pos)
				newMove := clipMovement(move, XY64{X: wall.X1, Y: wall.Y1})
				newPos := add(oldPlayer.pos, newMove)
				g.player.pos = newPos
			}
		}
	}
	return nil
}

func clipMovement(movement, collisionNormal XY64) XY64 {
	// Normalize the collision normal
	normal := normalize(collisionNormal)

	// Project movement onto the normal (component to block)
	projection := scale(normal, dot(movement, normal))

	// Subtract projection from the movement to get the clipped movement
	clippedMovement := subtract(movement, projection)

	return clippedMovement
}
