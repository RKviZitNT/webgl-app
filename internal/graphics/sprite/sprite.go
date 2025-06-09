//go:build js

package sprite

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/texture"
)

type Sprite struct {
	Texture *texture.Texture
	Rect    primitives.Rect
}

func NewSprite(texture *texture.Texture, rect primitives.Rect) *Sprite {
	return &Sprite{
		Texture: texture,
		Rect:    rect,
	}
}
