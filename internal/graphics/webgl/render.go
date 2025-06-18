//go:build js

package webgl

import (
	"webgl-app/internal/config"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/jsfunc"
)

func (ctx *GLContext) RenderSprite(sprite *Sprite, drawRect primitives.Rect, specular bool) {
	if sprite == nil {
		jsfunc.LogError("RenderSprite: sprite is nil")
		return
	}

	if sprite.Texture == nil {
		jsfunc.LogError("RenderSprite: sprite.Texture is nil")
		return
	}

	texture := sprite.Texture.GetTexture()
	if texture == nil || texture.IsUndefined() || texture.IsNull() {
		jsfunc.LogError("RenderSprite: texture not loaded")
		return
	}

	if specular {
		drawRect.Pos = drawRect.Move(*sprite.SpecularOffset)
	} else {
		drawRect.Pos = drawRect.Move(*sprite.Offset)
	}

	var (
		x1, y1, x2, y2 float32
		u1, v1, u2, v2 float32
	)

	x1 = float32(drawRect.Left() * ctx.Screen.Scale.X)
	y1 = float32(drawRect.Top() * ctx.Screen.Scale.Y)
	if sprite.Rect == nil {
		x2 = float32((drawRect.Left() + drawRect.Width()*sprite.Scale) * ctx.Screen.Scale.X)
		y2 = float32((drawRect.Top() + drawRect.Height()*sprite.Scale) * ctx.Screen.Scale.Y)
	} else {
		x2 = float32((drawRect.Left() + sprite.Rect.Width()*sprite.Scale) * ctx.Screen.Scale.X)
		y2 = float32((drawRect.Top() + sprite.Rect.Height()*sprite.Scale) * ctx.Screen.Scale.Y)
	}

	if sprite.Rect == nil {
		u1, v1 = float32(0), float32(0)
		u2, v2 = float32(1), float32(1)
	} else {
		u1 = float32(sprite.Rect.Left()) / float32(sprite.Texture.Width)
		v1 = float32(sprite.Rect.Top()) / float32(sprite.Texture.Height)
		u2 = float32(sprite.Rect.Right()) / float32(sprite.Texture.Width)
		v2 = float32(sprite.Rect.Bottom()) / float32(sprite.Texture.Height)
	}

	if specular {
		u1, u2 = u2, u1
	}

	bufferData := []float32{
		x1, y1, u1, v1,
		x2, y1, u2, v1,
		x1, y2, u1, v2,
		x2, y2, u2, v2,
	}

	ctx.textureQueue.addCommand(bufferData, sprite.Texture)
}

func (ctx *GLContext) RenderRect(rect primitives.Rect, color Color) {
	if !config.ProgramConfig.Debug {
		return
	}

	var (
		x1, y1, x2, y2 float32
	)

	x1 = float32(rect.Left() * ctx.Screen.Scale.X)
	y1 = float32(rect.Top() * ctx.Screen.Scale.Y)
	x2 = float32(rect.Right() * ctx.Screen.Scale.X)
	y2 = float32(rect.Bottom() * ctx.Screen.Scale.Y)

	bufferData := []float32{
		x1, y1,
		x2, y1,
		x2, y2,
		x1, y2,
	}

	ctx.debugQueue.addCommand(bufferData, color)
}
