//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/config"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/resourceloader"
)

type GLContext struct {
	GL        js.Value
	Program   js.Value
	Screen    *Screen
	drawQueue []DrawCommand
}

func NewWebGLContext(canvasId string) (*GLContext, error) {
	screen, err := NewScreen(canvasId, primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(config.GlobalConfig.Window.Width, config.GlobalConfig.Window.Height)))
	if err != nil {
		return nil, err
	}

	gl := screen.Canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		gl = screen.Canvas.Call("getContext", "experimental-webgl")
		if gl.IsNull() {
			return nil, fmt.Errorf("WebGL not supported")
		}
	}

	return &GLContext{
		GL:        gl,
		Screen:    screen,
		drawQueue: make([]DrawCommand, 0),
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
		jsfunc.LogInfo("Vertex shader loaded")
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
		jsfunc.LogInfo("Fragment shader loaded")
	}

	js.Global().Call("setLoadingProgress", 4, "Compiling vertex shader...")
	vShader, err := ctx.CompileShader(vertexSrc, gl.Get("VERTEX_SHADER"))
	if err != nil {
		return fmt.Errorf("vertex shader compilation failed: %v", err)
	} else {
		jsfunc.LogInfo("Vertex shader compiled")
	}

	js.Global().Call("setLoadingProgress", 6, "Compiling vertex shader...")
	fShader, err := ctx.CompileShader(fragmentSrc, gl.Get("FRAGMENT_SHADER"))
	if err != nil {
		return fmt.Errorf("fragment shader compilation failed: %v", err)
	} else {
		jsfunc.LogInfo("Fragment shader compiled")
	}

	js.Global().Call("setLoadingProgress", 8, "Creating program...")
	ctx.Program, err = ctx.CreateProgram(vShader, fShader)
	if err != nil {
		return fmt.Errorf("program creation failed: %v", err)
	} else {
		jsfunc.LogInfo("Program created")
	}

	return nil
}
