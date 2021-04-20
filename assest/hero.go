package main

type Location struct {
	x        int32
	y        int32
	offset_x int32
	offset_y int32
}

type Hero struct {
	Hp        int
	Mp        int
	MovePoint int
	location  Location
}
