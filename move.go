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
		g.player.posX += g.player.dirX * moveSpeed
		g.player.posY += g.player.dirY * moveSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.posX -= g.player.dirX * moveSpeed
		g.player.posY -= g.player.dirY * moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(turnSpeed) - g.player.dirY*math.Sin(turnSpeed)
		g.player.dirY = oldDirX*math.Sin(turnSpeed) + g.player.dirY*math.Cos(turnSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(turnSpeed) - g.player.planeY*math.Sin(turnSpeed)
		g.player.planeY = oldPlaneX*math.Sin(turnSpeed) + g.player.planeY*math.Cos(turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		oldDirX := g.player.dirX
		g.player.dirX = g.player.dirX*math.Cos(-turnSpeed) - g.player.dirY*math.Sin(-turnSpeed)
		g.player.dirY = oldDirX*math.Sin(-turnSpeed) + g.player.dirY*math.Cos(-turnSpeed)
		oldPlaneX := g.player.planeX
		g.player.planeX = g.player.planeX*math.Cos(-turnSpeed) - g.player.planeY*math.Sin(-turnSpeed)
		g.player.planeY = oldPlaneX*math.Sin(-turnSpeed) + g.player.planeY*math.Cos(-turnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.Key0) {
		fmt.Printf("posX: %v, posY: %v, dirX: %v, dirY: %v, planeX: %v, planeY: %v\n", g.player.posX, g.player.posY, g.player.dirX, g.player.dirY, g.player.planeX, g.player.planeY)
	}

	for _, wall := range walls {
		if dist, hit := rayIntersectsSegment(g.player.posX, g.player.posY, g.player.dirX, g.player.dirY, wall); hit {
			if dist < 0.5 {
				g.player = oldPlayer
			}
		}
	}
	return nil
}
