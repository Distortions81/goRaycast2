package main

// Define a struct for a 2D vector with start and end points
type line32 struct {
	X1, Y1, X2, Y2 float32
}

type pos32 struct {
	X, Y float32
}

// Game struct to hold game state
type Game struct {
	camera,
	start,
	lastMouse pos32

	createMode, pStartMode,
	firstClick,
	secondClick bool

	screenWidth,
	screenHeight int
}
