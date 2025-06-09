//go:build js

package level

import (
	"webgl-app/internal/graphics/texture"
)

type Level struct {
	Background *texture.Texture
	Floor      *texture.Texture
}

func NewLevel(background, floor *texture.Texture) *Level {
	return &Level{
		Background: background,
		Floor:      floor,
	}
}
