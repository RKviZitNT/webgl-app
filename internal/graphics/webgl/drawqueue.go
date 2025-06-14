//go:build js

package webgl

import "syscall/js"

// ----- Texture queue -----

type textureCommand struct {
	bufferData []float32
	texture    *Texture
}

type textureQueue struct {
	program js.Value
	queue   []textureCommand
}

func newTextureQueue(program js.Value) *textureQueue {
	return &textureQueue{
		program: program,
		queue:   make([]textureCommand, 0),
	}
}

func (tq *textureQueue) addCommand(bufferData []float32, texture *Texture) {
	tq.queue = append(tq.queue, textureCommand{
		bufferData: bufferData,
		texture:    texture,
	})
}

// ----- Debug queue -----

type debugCommand struct {
	bufferData []float32
	thickness  float32
	color      Color
}

type debugQueue struct {
	program js.Value
	queue   []debugCommand
}

func newDebugQueue(program js.Value) *debugQueue {
	return &debugQueue{
		program: program,
		queue:   make([]debugCommand, 0),
	}
}

func (dq *debugQueue) addCommand(bufferData []float32, thickness float32, color Color) {
	dq.queue = append(dq.queue, debugCommand{
		bufferData: bufferData,
		thickness:  thickness,
		color:      color,
	})
}

// -----------------------

func (ctx *GLContext) drawTextureQueue() {
	gl := ctx.GL
	program := ctx.textureQueue.program

	if len(ctx.textureQueue.queue) == 0 {
		return
	}

	gl.Call("useProgram", program)

	uResolution := gl.Call("getUniformLocation", program, "uResolution")
	gl.Call("uniform2f", uResolution, ctx.Screen.ScreenRect.Width(), ctx.Screen.ScreenRect.Height())

	for _, cmd := range ctx.textureQueue.queue {
		buffer := ctx.createBuffer(cmd.bufferData, gl.Get("STATIC_DRAW"))
		gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

		aPosition := gl.Call("getAttribLocation", program, "aPosition")
		aTexCoords := gl.Call("getAttribLocation", program, "aTexCoords")

		gl.Call("vertexAttribPointer", aPosition, 2, gl.Get("FLOAT"), false, 16, 0)
		gl.Call("enableVertexAttribArray", aPosition)

		gl.Call("vertexAttribPointer", aTexCoords, 2, gl.Get("FLOAT"), false, 16, 8)
		gl.Call("enableVertexAttribArray", aTexCoords)

		gl.Call("activeTexture", gl.Get("TEXTURE0"))
		gl.Call("bindTexture", gl.Get("TEXTURE_2D"), *cmd.texture.GetTexture())

		uTexture := gl.Call("getUniformLocation", program, "uTexture")
		gl.Call("uniform1i", uTexture, 0)

		gl.Call("drawArrays", gl.Get("TRIANGLE_STRIP"), 0, 4)
	}

	ctx.textureQueue.queue = make([]textureCommand, 0)
}

func (ctx *GLContext) drawDebugQueue() {
	gl := ctx.GL
	program := ctx.debugQueue.program

	if len(ctx.debugQueue.queue) == 0 {
		return
	}

	gl.Call("useProgram", program)

	uResolution := gl.Call("getUniformLocation", program, "uResolution")
	gl.Call("uniform2f", uResolution, ctx.Screen.ScreenRect.Width(), ctx.Screen.ScreenRect.Height())

	for _, cmd := range ctx.debugQueue.queue {
		buffer := ctx.createBuffer(cmd.bufferData, gl.Get("STATIC_DRAW"))
		gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

		aPosition := gl.Call("getAttribLocation", program, "aPosition")
		gl.Call("vertexAttribPointer", aPosition, 2, gl.Get("FLOAT"), false, 8, 0)
		gl.Call("enableVertexAttribArray", aPosition)

		uColor := gl.Call("getUniformLocation", program, "uColor")
		gl.Call("uniform4f", uColor, cmd.color.R, cmd.color.G, cmd.color.B, cmd.color.A)

		gl.Call("drawArrays", gl.Get("LINE_LOOP"), 0, 4)
	}

	ctx.debugQueue.queue = make([]debugCommand, 0)
}

func (ctx *GLContext) DrawQueue() {
	gl := ctx.GL
	gl.Call("viewport", ctx.Screen.ScreenRect.Left(), ctx.Screen.ScreenRect.Top(), ctx.Screen.ScreenRect.Width(), ctx.Screen.ScreenRect.Height())
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	ctx.drawTextureQueue()
	ctx.drawDebugQueue()
}
