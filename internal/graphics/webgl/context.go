//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/resourceloader"
)

type GLContext struct {
	Canvas     js.Value
	CanvasRect primitives.Rect
	GL         js.Value
	Program    js.Value
}

func NewWebGLCtx(canvasID string) (*GLContext, error) {
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

	return &GLContext{
		Canvas:     canvas,
		CanvasRect: primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(width, height)),
		GL:         gl,
	}, nil
}

func (ctx *GLContext) InitWebGL() error {
	gl := ctx.GL

	shaderLoaded := make(chan bool, 1)
	var loadErr error
	var vertexSrc, fragmentSrc string

	js.Global().Call("setLoadingProgress", 0, "Loading vertex shader...")
	resourceloader.LoadFile("shaders/vertex.glsl",
		func(src js.Value) {
			vertexSrc = src.String()
			shaderLoaded <- true
		},
		func(err error) {
			loadErr = err
			shaderLoaded <- false
		})

	if !<-shaderLoaded {
		return fmt.Errorf("failed to load vertex shader: %v", loadErr)
	} else {
		js.Global().Get("console").Call("log", "Vertex shader loaded")
	}

	js.Global().Call("setLoadingProgress", 2, "Loading fragment shader...")
	resourceloader.LoadFile("shaders/fragment.glsl",
		func(src js.Value) {
			fragmentSrc = src.String()
			shaderLoaded <- true
		},
		func(err error) {
			loadErr = err
			shaderLoaded <- false
		})

	if !<-shaderLoaded {
		return fmt.Errorf("failed to load fragment shader: %v", loadErr)
	} else {
		js.Global().Get("console").Call("log", "Fragment shader loaded")
	}

	js.Global().Call("setLoadingProgress", 4, "Compiling vertex shader...")
	vShader, err := ctx.CompileShader(vertexSrc, gl.Get("VERTEX_SHADER"))
	if err != nil {
		return fmt.Errorf("vertex shader compilation failed: %v", err)
	} else {
		js.Global().Get("console").Call("log", "Vertext shader compiled")
	}

	js.Global().Call("setLoadingProgress", 6, "Compiling vertex shader...")
	fShader, err := ctx.CompileShader(fragmentSrc, gl.Get("FRAGMENT_SHADER"))
	if err != nil {
		return fmt.Errorf("fragment shader compilation failed: %v", err)
	} else {
		js.Global().Get("console").Call("log", "Fragment shader compiled")
	}

	js.Global().Call("setLoadingProgress", 8, "Creating program...")
	ctx.Program, err = ctx.CreateProgram(vShader, fShader)
	if err != nil {
		return fmt.Errorf("program creation failed: %v", err)
	} else {
		js.Global().Get("console").Call("log", "Program created")
	}

	return nil
}
