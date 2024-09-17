package main

// Define a struct for a 2D vector with start and end points
type Vector2D struct {
	X1, Y1, X2, Y2 float64
}

type XY struct {
	X, Y float64
}

// Game struct to hold game state
type Game struct {
	camera,
	start,
	lastMouse XY

	createMode,
	firstClick,
	secondClick bool

	screenWidth,
	screenHeight int
}
