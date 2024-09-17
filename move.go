package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	moveSpeed  = 0.02
	turnSpeed  = 0.05
	playerSize = 0.5

	friction = 0.009
	maxSpeed = 0.1
)

var player playerData

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if player.speed < maxSpeed {
			player.speed += moveSpeed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		if player.speed > -maxSpeed {
			player.speed -= moveSpeed
		}
	} else {
		if player.speed > 0 {
			if player.speed < friction {
				player.speed = 0
			}
			player.speed -= friction
		} else if player.speed < 0 {
			if player.speed > -friction {
				player.speed = 0
			}
			player.speed += friction
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.angle -= turnSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.angle += turnSpeed
	}

	player.velocity = AngleToVelocity(player.angle, player.speed)

	player.pos = addXY(player.pos, player.velocity)
	return nil
}

func clipMovement(movement, collisionNormal XY64) XY64 {
	// Normalize the collision normal
	normal := normalizeXY(collisionNormal)

	// Project movement onto the normal (component to block)
	projection := scaleXY(normal, dotXY(movement, normal))

	// Subtract projection from the movement to get the clipped movement
	clippedMovement := subXY(movement, projection)

	return clippedMovement
}

// Function to convert an angle in radians to a velocity vector with momentum
func AngleToVelocity(angle float64, magnitude float64) XY64 {
	// Calculate X and Y components using trigonometry
	vx := magnitude * math.Cos(angle)
	vy := magnitude * math.Sin(angle)
	return XY64{X: vx, Y: vy}
}
