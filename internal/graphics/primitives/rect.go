package primitives

type Rect struct {
	Pos  Vec2
	Size Vec2
}

func NewRect(x, y, width, height float64) Rect {
	return Rect{
		Pos:  NewVec2(x, y),
		Size: NewVec2(width, height),
	}
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
	return NewVec2(r.Pos.X+r.Size.X/2, r.Pos.Y+r.Size.Y/2)
}

func (r *Rect) SetLeft(x float64) {
	r.Pos.X = x
}

func (r *Rect) SetRight(x float64) {
	r.Pos.X = x - r.Size.X
}

func (r *Rect) SetTop(y float64) {
	r.Pos.Y = y
}

func (r *Rect) SetBottom(y float64) {
	r.Pos.Y = y - r.Size.Y
}

func (r *Rect) SetCenter(pos Vec2) {
	r.Pos = NewVec2(pos.X-(r.Size.X/2), pos.Y-(r.Size.Y/2))
}

func (r *Rect) Width() float64 {
	return r.Size.X
}

func (r *Rect) Height() float64 {
	return r.Size.Y
}

func (r *Rect) Move(offset Vec2) Rect {
	return Rect{
		Pos:  r.Pos.AddVec2(offset),
		Size: r.Size,
	}
}

func (r *Rect) Intersection(other Rect) bool {
	return r.Left() < other.Right() && r.Right() > other.Left() && r.Top() < other.Bottom() && r.Bottom() > other.Top()
}

func (r *Rect) ContsinsVec2(point Vec2) bool {
	return point.X >= r.Left() && point.X <= r.Right() && point.Y >= r.Top() && point.Y <= r.Bottom()
}
