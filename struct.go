package main

type Vec64 struct {
	X1, Y1, X2, Y2 float64
}

type Vec32 struct {
	X1, Y1, X2, Y2 float64
}

type XY64 struct {
	X, Y float64
}

type XY32 struct {
	X, Y float64
}

type Player struct {
	pos, dir, plane XY64
}

type Game struct {
	player Player
}
