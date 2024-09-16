package main

var walls = []Vector{
	// Outer boundary
	{0.0, 0.0, 25.0, 0.0},
	{25.0, 0.0, 25.0, 25.0},
	{25.0, 25.0, 0.0, 25.0},
	{0.0, 25.0, 0.0, 0.0},

	// Main corridor
	{2.5, 2.5, 22.5, 2.5},
	{22.5, 2.5, 22.5, 22.5},
	{22.5, 22.5, 2.5, 22.5},
	{2.5, 22.5, 2.5, 2.5},

	// Rooms
	{5.0, 5.0, 10.0, 5.0},
	{10.0, 5.0, 10.0, 10.0},
	{10.0, 10.0, 5.0, 10.0},
	{5.0, 10.0, 5.0, 5.0},

	{15.0, 5.0, 20.0, 5.0},
	{20.0, 5.0, 20.0, 10.0},
	{20.0, 10.0, 15.0, 10.0},
	{15.0, 10.0, 15.0, 5.0},

	// Additional corridors
	{5.0, 15.0, 10.0, 15.0},
	{10.0, 15.0, 10.0, 20.0},
	{10.0, 20.0, 5.0, 20.0},
	{5.0, 20.0, 5.0, 15.0},

	{15.0, 15.0, 20.0, 15.0},
	{20.0, 15.0, 20.0, 20.0},
	{20.0, 20.0, 15.0, 20.0},
	{15.0, 20.0, 15.0, 15.0},
}
