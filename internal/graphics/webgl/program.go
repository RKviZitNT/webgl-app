//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
)

func (ctx *GLContext) CreateProgram(vertexShader, fragmentShader js.Value) (js.Value, error) {
	gl := ctx.GL

	program := gl.Call("createProgram")
	gl.Call("attachShader", program, vertexShader)
	gl.Call("attachShader", program, fragmentShader)
	gl.Call("linkProgram", program)

	if !gl.Call("getProgramParameter", program, gl.Get("LINK_STATUS")).Bool() {
		log := gl.Call("getProgramInfoLog", program).String()
		return js.Null(), fmt.Errorf("program linking failed: %s", log)
	}

	gl.Call("useProgram", program)
	return program, nil
}
