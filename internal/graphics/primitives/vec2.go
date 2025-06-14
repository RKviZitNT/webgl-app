package primitives

import "math"

type Vec2 struct {
	X, Y float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

func (v *Vec2) Set(x, y float64) {
	v.X = x
	v.Y = y
}

func (v *Vec2) AddValue(value float64) Vec2 {
	return Vec2{
		X: v.X + value,
		Y: v.Y + value,
	}
}

func (v *Vec2) SubValue(value float64) Vec2 {
	return Vec2{
		X: v.X - value,
		Y: v.Y - value,
	}
}

func (v *Vec2) MulValue(value float64) Vec2 {
	return Vec2{
		X: v.X * value,
		Y: v.Y * value,
	}
}

func (v *Vec2) DivValue(value float64) Vec2 {
	return Vec2{
		X: v.X / value,
		Y: v.Y / value,
	}
}

func (v *Vec2) AddVec2(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v *Vec2) SubVec2(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v *Vec2) MulVec2(other Vec2) Vec2 {
	return Vec2{
		X: v.X * other.X,
		Y: v.Y * other.Y,
	}
}

func (v *Vec2) DivVec2(other Vec2) Vec2 {
	return Vec2{
		X: v.X / other.X,
		Y: v.Y / other.Y,
	}
}

func (v *Vec2) Negative() Vec2 {
	return Vec2{
		X: -v.X,
		Y: -v.Y,
	}
}

func (v *Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vec2) Normalize() Vec2 {
	len := v.Length()
	if len > 0 {
		return Vec2{
			X: v.X / len,
			Y: v.Y / len,
		}
	}
	return Vec2{}
}
