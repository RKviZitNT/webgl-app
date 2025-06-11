//go:build js

package webgl

import "webgl-app/internal/graphics/texture"

type DrawCommand struct {
	BufferData []float32
	Texture    *texture.Texture
}

func (ctx *GLContext) addDrawCommand(bufferData []float32, texture *texture.Texture) {
	ctx.drawQueue = append(ctx.drawQueue, DrawCommand{
		BufferData: bufferData,
		Texture:    texture,
	})
}

func (ctx *GLContext) clearDrawQueue() {
	ctx.drawQueue = make([]DrawCommand, 0)
}

func (ctx *GLContext) FlushDrawQueue() {
	gl := ctx.GL
	program := ctx.Program

	if len(ctx.drawQueue) == 0 {
		return
	}

	gl.Call("viewport", 0, 0, ctx.CanvasRect.Width(), ctx.CanvasRect.Height())
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))
	gl.Call("useProgram", program)

	gl.Call("enable", gl.Get("BLEND"))
	gl.Call("blendFunc", gl.Get("SRC_ALPHA"), gl.Get("ONE_MINUS_SRC_ALPHA"))

	for _, cmd := range ctx.drawQueue {
		buffer := ctx.CreateBuffer(cmd.BufferData, gl.Get("STATIC_DRAW"))
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
		gl.Call("bindTexture", gl.Get("TEXTURE_2D"), *cmd.Texture.GetTexture())

		texUniform := gl.Call("getUniformLocation", program, "u_texture")
		gl.Call("uniform1i", texUniform, 0)

		gl.Call("drawArrays", gl.Get("TRIANGLE_STRIP"), 0, 4)
	}

	ctx.clearDrawQueue()
}
