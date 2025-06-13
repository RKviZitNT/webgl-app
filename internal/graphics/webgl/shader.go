//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/config"
	"webgl-app/internal/resourceloader"
)

func (ctx *GLContext) compileShaders(vertexSrc, fragmentSrc string) (js.Value, js.Value, error) {
	gl := ctx.GL

	compileShader := func(source string, shaderType js.Value) (js.Value, error) {
		shader := gl.Call("createShader", shaderType)
		gl.Call("shaderSource", shader, source)
		gl.Call("compileShader", shader)

		if !gl.Call("getShaderParameter", shader, gl.Get("COMPILE_STATUS")).Bool() {
			log := gl.Call("getShaderInfoLog", shader).String()
			return js.Null(), fmt.Errorf("shader compilation failed: %s", log)
		}
		return shader, nil
	}

	vertexShader, err := compileShader(vertexSrc, gl.Get("VERTEX_SHADER"))
	if err != nil {
		return js.Null(), js.Null(), fmt.Errorf("vertex shader error: %v", err)
	}

	fragmentShader, err := compileShader(fragmentSrc, gl.Get("FRAGMENT_SHADER"))
	if err != nil {
		return js.Null(), js.Null(), fmt.Errorf("fragment shader error: %v", err)
	}

	return vertexShader, fragmentShader, nil
}

func (ctx *GLContext) loadShaders(shaders config.Shaders) (string, string, error) {
	var (
		loadErr     error
		vertexSrc   string
		fragmentSrc string
	)

	done := make(chan struct{}, 0)

	resourceloader.LoadFile(shaders.Vertex,
		func(src js.Value) {
			vertexSrc = src.String()

			resourceloader.LoadFile(shaders.Fragment,
				func(src js.Value) {
					fragmentSrc = src.String()
					close(done)
				},
				func(err error) {
					loadErr = err
					close(done)
				})
		},
		func(err error) {
			loadErr = err
			close(done)
		})

	<-done

	if loadErr != nil {
		return "", "", loadErr
	}

	return vertexSrc, fragmentSrc, nil
}
