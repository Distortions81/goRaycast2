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

type playerData struct {
	pos      XY64
	velocity XY64
	size     float64
	angle    float64
	speed    float64
}

type Game struct {
}
