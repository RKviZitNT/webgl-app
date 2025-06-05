//go:build js

package graphics

import (
	"fmt"
	"syscall/js"
)

func CreateProgram(gl js.Value, vertexShader, fragmentShader js.Value) (js.Value, error) {
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
