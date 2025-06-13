//go:build js

package webgl

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/jsfunc"
)

func (ctx *GLContext) RenderSprite(sprite *Sprite, drawRect *primitives.Rect) {
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

	if sprite.Offset != nil {
		drawRect = drawRect.Move(sprite.Offset)
	}

	var (
		x1, y1, x2, y2 float32
		u1, v1, u2, v2 float32
	)

	x1 = float32(drawRect.Left())
	y1 = float32(drawRect.Top())
	if sprite.Rect == nil {
		x2 = float32(drawRect.Left() + drawRect.Width()*sprite.Scale)
		y2 = float32(drawRect.Top() + drawRect.Height()*sprite.Scale)
	} else {
		x2 = float32(drawRect.Left() + sprite.Rect.Width()*sprite.Scale)
		y2 = float32(drawRect.Top() + sprite.Rect.Height()*sprite.Scale)
	}

	if sprite.Rect == nil {
		u1, v1 = float32(0), float32(0)
		u2, v2 = float32(1), float32(1)
	} else {
		texSize := sprite.Texture.Size()
		u1 = float32(sprite.Rect.Left()) / float32(texSize.X)
		v1 = float32(sprite.Rect.Top()) / float32(texSize.Y)
		u2 = float32(sprite.Rect.Right()) / float32(texSize.X)
		v2 = float32(sprite.Rect.Bottom()) / float32(texSize.Y)
	}

	bufferData := []float32{
		x1, y1, u1, v1,
		x2, y1, u2, v1,
		x1, y2, u1, v2,
		x2, y2, u2, v2,
	}

	ctx.textureQueue.addCommand(bufferData, sprite.Texture)
}

func (ctx *GLContext) RenderRect(rect *primitives.Rect, thickness float32, color *Color) {
	var (
		x1, y1, x2, y2 float32
	)

	x1 = float32(rect.Left())
	y1 = float32(rect.Top())
	x2 = float32(rect.Right())
	y2 = float32(rect.Bottom())

	bufferData := []float32{
		x1, y1,
		x2, y1,
		x2, y2,
		x1, y2,
	}

	ctx.debugQueue.addCommand(bufferData, thickness, color)
}
