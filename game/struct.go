package main

type line32 struct {
	X1, Y1, X2, Y2 float32
}

type pos32 struct {
	X, Y float32
}

type playerData struct {
	pos      pos32
	velocity pos32
	size     float32
	angle    float32
	speed    float32
}

type Game struct {
}
