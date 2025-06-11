//go:build js

package webgl

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/jsfunc"
)

func (ctx *GLContext) RenderTexture(texture *texture.Texture, rect primitives.Rect) {
	if texture == nil {
		jsfunc.LogError("RenderTexture: texture is nil")
		return
	}

	tex := texture.GetTexture()
	if tex == nil || tex.IsUndefined() || tex.IsNull() {
		jsfunc.LogError("RenderTexture: texture not loaded")
		return
	}

	x1 := float32(rect.Left())
	y1 := float32(rect.Top())
	x2 := float32(rect.Right())
	y2 := float32(rect.Bottom())

	u1, v1 := float32(0.0), float32(0.0)
	u2, v2 := float32(1.0), float32(1.0)

	bufferData := []float32{
		x1, y1, u1, v1,
		x2, y1, u2, v1,
		x1, y2, u1, v2,
		x2, y2, u2, v2,
	}

	ctx.addDrawCommand(bufferData, texture)
}

func (ctx *GLContext) RenderSprite(sprite *sprite.Sprite, pos primitives.Vec2, scale float64) {
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

	x1 := float32(pos.X)
	y1 := float32(pos.Y)
	x2 := float32(pos.X + sprite.Rect.Width()*scale)
	y2 := float32(pos.Y + sprite.Rect.Height()*scale)

	texSize := sprite.Texture.Size()
	u1 := float32(sprite.Rect.Left()) / float32(texSize.X)
	v1 := float32(sprite.Rect.Top()) / float32(texSize.Y)
	u2 := float32(sprite.Rect.Right()) / float32(texSize.X)
	v2 := float32(sprite.Rect.Bottom()) / float32(texSize.Y)

	bufferData := []float32{
		x1, y1, u1, v1,
		x2, y1, u2, v1,
		x1, y2, u1, v2,
		x2, y2, u2, v2,
	}

	ctx.addDrawCommand(bufferData, sprite.Texture)
}
