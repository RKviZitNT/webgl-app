//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/resourcemanager"
)

type GLContext struct {
	Canvas     js.Value
	CanvasSize primitives.Vec2
	GL         js.Value
	Program    js.Value
}

func InitWebGL(canvasID string) (*GLContext, error) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", canvasID)

	width := js.Global().Get("innerWidth").Float()
	height := js.Global().Get("innerHeight").Float()
	canvas.Set("width", width)
	canvas.Set("height", height)

	gl := canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		gl = canvas.Call("getContext", "experimental-webgl")
		if gl.IsNull() {
			return nil, fmt.Errorf("WebGL not supported")
		}
	}

	shaderLoaded := make(chan bool, 1)
	var loadErr error
	var vertexSrc, fragmentSrc string

	js.Global().Call("setLoadingProgress", 0, "Loading vertex shader...")
	resourcemanager.LoadFile("shaders/vertex.glsl",
		func(src string) {
			vertexSrc = src
			shaderLoaded <- true
		},
		func(err error) {
			loadErr = err
			shaderLoaded <- false
		})

	if !<-shaderLoaded {
		return nil, fmt.Errorf("failed to load vertex shader: %v", loadErr)
	} else {
		println("Vertex shader loaded. Size:", len(vertexSrc))
	}

	js.Global().Call("setLoadingProgress", 20, "Loading fragment shader...")
	resourcemanager.LoadFile("shaders/fragment.glsl",
		func(src string) {
			fragmentSrc = src
			shaderLoaded <- true
		},
		func(err error) {
			loadErr = err
			shaderLoaded <- false
		})

	if !<-shaderLoaded {
		return nil, fmt.Errorf("failed to load fragment shader: %v", loadErr)
	} else {
		println("Fragment shader loaded. Size:", len(fragmentSrc))
	}

	js.Global().Call("setLoadingProgress", 40, "Compiling vertex shader...")
	vShader, err := CompileShader(gl, vertexSrc, gl.Get("VERTEX_SHADER"))
	if err != nil {
		return nil, fmt.Errorf("vertex shader compilation failed: %v", err)
	} else {
		println("Vertext shader compiled")
	}

	js.Global().Call("setLoadingProgress", 60, "Compiling vertex shader...")
	fShader, err := CompileShader(gl, fragmentSrc, gl.Get("FRAGMENT_SHADER"))
	if err != nil {
		return nil, fmt.Errorf("fragment shader compilation failed: %v", err)
	} else {
		println("Fragment shader compiled")
	}

	js.Global().Call("setLoadingProgress", 80, "Creating program...")
	program, err := CreateProgram(gl, vShader, fShader)
	if err != nil {
		return nil, fmt.Errorf("program creation failed: %v", err)
	} else {
		println("Program created")
	}

	return &GLContext{
		Canvas:     canvas,
		CanvasSize: primitives.NewVec2(width, height),
		GL:         gl,
		Program:    program,
	}, nil
}
