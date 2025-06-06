package primitives

type Rect struct {
	Pos  Vec2
	Size Vec2
}

func (r *Rect) Left() float64 {
	return r.Pos.X
}

func (r *Rect) Right() float64 {
	return r.Pos.X + r.Size.X
}

func (r *Rect) Top() float64 {
	return r.Pos.Y
}

func (r *Rect) Bottom() float64 {
	return r.Pos.Y + r.Size.Y
}

func (r *Rect) Center() Vec2 {
	return Vec2{
		X: r.Pos.X + r.Size.X/2,
		Y: r.Pos.Y + r.Size.Y/2,
	}
}

func (r *Rect) Move(offset Vec2) {
	r.Pos.AddVec2(offset)
}

func (r *Rect) Intersection(other Rect) bool {
	return r.Left() < other.Right() && r.Right() > other.Left() && r.Top() < other.Bottom() && r.Bottom() > other.Top()
}

func (r *Rect) ContsinsVec2(point Vec2) bool {
	return point.X >= r.Left() && point.X <= r.Right() && point.Y >= r.Top() && point.Y <= r.Bottom()
}
