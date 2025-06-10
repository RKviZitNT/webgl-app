//go:build js

package level

import (
	"webgl-app/internal/graphics/texture"
)

type LevelName string

const (
	DefaultLevel LevelName = "default"
)

type Level struct {
	Name       LevelName
	Background *texture.Texture
	Floor      *texture.Texture
}

func NewLevel(name LevelName, background, floor *texture.Texture) *Level {
	return &Level{
		Name:       name,
		Background: background,
		Floor:      floor,
	}
}
