//go:build js

package webgl

import (
	"webgl-app/internal/graphics/primitives"
)

type Sprite struct {
	Texture *Texture
	Rect    *primitives.Rect
	Scale   float64
	Offset  *primitives.Vec2
}

func NewSprite(texture *Texture, rect *primitives.Rect, scale float64, offset *primitives.Vec2) *Sprite {
	return &Sprite{
		Texture: texture,
		Rect:    rect,
		Scale:   scale,
		Offset:  offset,
	}
}
