//go:build js

package level

import (
	"webgl-app/internal/graphics/webgl"
)

type LevelName string

const (
	DefaultLevel LevelName = "default"
)

type Level struct {
	Name       LevelName
	Background *webgl.Sprite
}

func NewLevel(name LevelName, background *webgl.Sprite) *Level {
	return &Level{
		Name:       name,
		Background: background,
	}
}
