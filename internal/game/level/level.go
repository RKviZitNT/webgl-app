//go:build js

package level

import (
	"fmt"
	"webgl-app/internal/graphics/webgl"
)

type Level struct {
	Name       string
	background *webgl.Sprite
}

func NewLevel(name string, tex *webgl.Texture) (*Level, error) {
	if tex == nil {
		return nil, fmt.Errorf("Texture is nil")
	}

	background := webgl.NewSprite(tex, nil, 1, nil, nil)

	return &Level{
		Name:       name,
		background: background,
	}, nil
}

func (l *Level) Draw(glCtx *webgl.GLContext) {
	glCtx.RenderSprite(l.background, glCtx.Screen.BaseScreenRect, false)
}
