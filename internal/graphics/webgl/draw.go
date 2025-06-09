//go:build js

package webgl

import (
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/texture"
)

func (ctx *GLContext) draw(bufferData []float32, texture *js.Value) {
	gl := ctx.GL
	program := ctx.Program

	buffer := ctx.CreateBuffer(bufferData, gl.Get("STATIC_DRAW"))
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

	posAttrib := gl.Call("getAttribLocation", program, "a_position")
	texAttrib := gl.Call("getAttribLocation", program, "a_texCoord")

	gl.Call("vertexAttribPointer", posAttrib, 2, gl.Get("FLOAT"), false, 16, 0)
	gl.Call("enableVertexAttribArray", posAttrib)

	gl.Call("vertexAttribPointer", texAttrib, 2, gl.Get("FLOAT"), false, 16, 8)
	gl.Call("enableVertexAttribArray", texAttrib)

	resUniform := gl.Call("getUniformLocation", program, "u_resolution")
	gl.Call("uniform2f", resUniform, ctx.CanvasRect.Width(), ctx.CanvasRect.Height())

	gl.Call("activeTexture", gl.Get("TEXTURE0"))
	gl.Call("bindTexture", gl.Get("TEXTURE_2D"), *texture)
	texUniform := gl.Call("getUniformLocation", program, "u_texture")
	gl.Call("uniform1i", texUniform, 0)

	gl.Call("enable", gl.Get("BLEND"))
	gl.Call("blendFunc", gl.Get("SRC_ALPHA"), gl.Get("ONE_MINUS_SRC_ALPHA"))

	gl.Call("drawArrays", gl.Get("TRIANGLE_STRIP"), 0, 4)
}

func (ctx *GLContext) DrawTexture(texture *texture.Texture, rect primitives.Rect) {
	if texture == nil {
		js.Global().Get("console").Call("error", "DrawTexture error: texture is nil")
		return
	}

	gl := ctx.GL
	program := ctx.Program

	gl.Call("useProgram", program)

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

	ctx.draw(bufferData, texture.GetTexture())
}

func (ctx *GLContext) DrawSprite(sprite *sprite.Sprite, pos primitives.Vec2, scale float64) {
	if sprite.Texture == nil {
		js.Global().Get("console").Call("error", "Draw error: texture is nil")
		return
	}

	texture := sprite.Texture.GetTexture()
	if texture == nil || texture.IsUndefined() || texture.IsNull() {
		js.Global().Get("console").Call("error", "Draw error: texture not loaded")
		return
	}

	gl := ctx.GL
	program := ctx.Program

	gl.Call("useProgram", program)

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

	ctx.draw(bufferData, texture)
}
