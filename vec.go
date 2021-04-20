package main

type Vec struct {
	x float64
	y float64
	z float64
}

func (v *Vec) ToVec2() *Vec2 {
	return &Vec2{
		x: int(v.x),
		y: int(v.y),
		z: int(v.z),
	}
}

type Vec2 struct {
	x int
	y int
	z int
}
