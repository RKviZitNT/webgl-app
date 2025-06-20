//go:build js

package webgl

import (
	"webgl-app/internal/graphics/primitives"
)

type Sprite struct {
	Texture        *Texture
	Rect           *primitives.Rect
	Scale          float64
	Offset         *primitives.Vec2
	SpecularOffset *primitives.Vec2
}

func NewSprite(tex *Texture, rect *primitives.Rect, scale float64, offset *primitives.Vec2, specularOffset *primitives.Vec2) *Sprite {
	if offset == nil {
		offset = &primitives.Vec2{}
	}
	if specularOffset == nil {
		specularOffset = &primitives.Vec2{}
	}
	return &Sprite{
		Texture:        tex,
		Rect:           rect,
		Scale:          scale,
		Offset:         offset,
		SpecularOffset: specularOffset,
	}
}
