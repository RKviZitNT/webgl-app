//go:build js

package graphics

import (
	"fmt"
	"syscall/js"
)

func CompileShader(gl js.Value, source string, shaderType js.Value) (js.Value, error) {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)

	if !gl.Call("getShaderParameter", shader, gl.Get("COMPILE_STATUS")).Bool() {
		log := gl.Call("getShaderInfoLog", shader).String()
		return js.Null(), fmt.Errorf("shader compilation failed: %s", log)
	}
	return shader, nil
}
