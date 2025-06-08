//go:build js

package webgl

import (
	"syscall/js"
)

func CreateBuffer(gl js.Value, data []float32, usage js.Value) js.Value {
	buffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

	uint8Array := js.Global().Get("Float32Array").New(len(data))
	for i, v := range data {
		uint8Array.SetIndex(i, v)
	}

	// не работает (undefined: js.TypedArrayOf [js,wasm])
	// uint8Array := js.TypedArrayOf(data)
	// defer uint8Array.Release()

	gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), uint8Array, usage)
	return buffer
}
