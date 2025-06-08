//go:build js

package sprite

import (
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/graphics/webgl"
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

func (s *Sprite) Draw(gl js.Value, program js.Value, pos primitives.Vec2, canvasSize primitives.Vec2) {
	if s.Texture == nil {
		println("Draw error: texture is nil")
		return
	}

	texture := s.Texture.GetTexture()
	if texture == nil || texture.IsUndefined() || texture.IsNull() {
		println("Draw error: texture not loaded")
		return
	}

	gl.Call("useProgram", program)

	scale := 5.0

	x1 := float32(pos.X)
	y1 := float32(pos.Y)
	x2 := float32(pos.X + s.Rect.Width()*scale)
	y2 := float32(pos.Y + s.Rect.Height()*scale)

	texSize := s.Texture.Size()
	u1 := float32(s.Rect.Left()) / float32(texSize.X)
	v1 := float32(s.Rect.Top()) / float32(texSize.Y)
	u2 := float32(s.Rect.Right()) / float32(texSize.X)
	v2 := float32(s.Rect.Bottom()) / float32(texSize.Y)

	bufferData := []float32{
		x1, y1, u1, v1,
		x2, y1, u2, v1,
		x1, y2, u1, v2,
		x2, y2, u2, v2,
	}

	buffer := webgl.CreateBuffer(gl, bufferData, gl.Get("STATIC_DRAW"))
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

	posAttrib := gl.Call("getAttribLocation", program, "a_position")
	texAttrib := gl.Call("getAttribLocation", program, "a_texCoord")

	gl.Call("vertexAttribPointer", posAttrib, 2, gl.Get("FLOAT"), false, 16, 0)
	gl.Call("enableVertexAttribArray", posAttrib)

	gl.Call("vertexAttribPointer", texAttrib, 2, gl.Get("FLOAT"), false, 16, 8)
	gl.Call("enableVertexAttribArray", texAttrib)

	resUniform := gl.Call("getUniformLocation", program, "u_resolution")
	gl.Call("uniform2f", resUniform, canvasSize.X, canvasSize.Y)

	gl.Call("activeTexture", gl.Get("TEXTURE0"))
	gl.Call("bindTexture", gl.Get("TEXTURE_2D"), *texture)
	texUniform := gl.Call("getUniformLocation", program, "u_texture")
	gl.Call("uniform1i", texUniform, 0)

	gl.Call("enable", gl.Get("BLEND"))
	gl.Call("blendFunc", gl.Get("SRC_ALPHA"), gl.Get("ONE_MINUS_SRC_ALPHA"))

	gl.Call("drawArrays", gl.Get("TRIANGLE_STRIP"), 0, 4)
}
