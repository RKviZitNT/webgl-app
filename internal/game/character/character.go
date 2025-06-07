package character

import (
	"webgl-app/internal/graphics/primitives"
)

type Character struct {
	HitBox primitives.Rect
}

func NewCharacter(pos primitives.Vec2, size primitives.Vec2) *Character {
	return &Character{
		HitBox: primitives.Rect{
			Pos:  pos,
			Size: size,
		},
	}
}
