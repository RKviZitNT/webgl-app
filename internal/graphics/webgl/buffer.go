//go:build js

package webgl

import (
	"syscall/js"
)

func (ctx *GLContext) createBuffer(data []float32, usage js.Value) js.Value {
	gl := ctx.GL

	buffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

	uint8Array := js.Global().Get("Float32Array").New(len(data))
	for i, v := range data {
		uint8Array.SetIndex(i, v)
	}

	gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), uint8Array, usage)
	return buffer
}
