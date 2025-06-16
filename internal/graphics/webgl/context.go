//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/config"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/jsfunc"
)

type GLContext struct {
	GL           js.Value
	Screen       *Screen
	textureQueue *textureQueue
	debugQueue   *debugQueue
}

func NewWebGLContext(canvasId string) (*GLContext, error) {
	screen, err := newScreen(canvasId, primitives.NewRect(0, 0, config.ProgramConfig.Window.Width, config.ProgramConfig.Window.Height))
	if err != nil {
		return nil, err
	}

	gl := screen.canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		gl = screen.canvas.Call("getContext", "experimental-webgl")
		if gl.IsNull() {
			return nil, fmt.Errorf("WebGL not supported")
		}
	}

	gl.Call("clearColor", 0.0, 0.0, 0.0, 1.0)

	gl.Call("enable", gl.Get("BLEND"))
	gl.Call("blendFunc", gl.Get("SRC_ALPHA"), gl.Get("ONE_MINUS_SRC_ALPHA"))

	return &GLContext{
		GL:     gl,
		Screen: screen,
	}, nil
}

func (ctx *GLContext) InitWebGL() error {
	var shadersSources ShadersSources
	err := config.LoadSources("shaders-manifest.json", &shadersSources)
	if err != nil {
		return err
	}

	jsfunc.SetLoadingProgress(0, "Loading texture shaders...")
	texVertSrc, texFragSrc, err := ctx.loadShaders(shadersSources.TextureShaders)
	if err != nil {
		return fmt.Errorf("failed to load texture shaders: %v", err)
	}
	jsfunc.LogInfo("Texture shaders loaded")

	jsfunc.SetLoadingProgress(2, "Compiling texture shaders...")
	texVertShader, texFragShader, err := ctx.compileShaders(texVertSrc, texFragSrc)
	if err != nil {
		return fmt.Errorf("texture shaders compilation failed: %v", err)
	}
	jsfunc.LogInfo("Texture shaders compiled")

	jsfunc.SetLoadingProgress(8, "Creating program...")
	textureProgram, err := ctx.createProgram(texVertShader, texFragShader)
	if err != nil {
		return fmt.Errorf("program creation failed: %v", err)
	}
	ctx.textureQueue = newTextureQueue(textureProgram)
	jsfunc.LogInfo("Texture queue created")

	if config.ProgramConfig.Debug {
		jsfunc.SetLoadingProgress(4, "Loading debug shaders...")
		debugVertSrc, debugFragSrc, err := ctx.loadShaders(shadersSources.DebugShaders)
		if err != nil {
			return fmt.Errorf("failed to load debug shaders: %v", err)
		}
		jsfunc.LogInfo("Debug shaders loaded")

		jsfunc.SetLoadingProgress(6, "Compiling debug shaders...")
		debugVertShader, debugFragShader, err := ctx.compileShaders(debugVertSrc, debugFragSrc)
		if err != nil {
			return fmt.Errorf("debug shaders compilation failed: %v", err)
		}
		jsfunc.LogInfo("Debug shaders compiled")

		jsfunc.SetLoadingProgress(8, "Creating program...")
		debugProgram, err := ctx.createProgram(debugVertShader, debugFragShader)
		if err != nil {
			return fmt.Errorf("debug program creation failed: %v", err)
		}
		ctx.debugQueue = newDebugQueue(debugProgram)
		jsfunc.LogInfo("Debug queue created")
	}

	js.Global().Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.handleResizeScreen()
		return nil
	}))

	return nil
}

func (ctx *GLContext) handleResizeScreen() {
	ctx.Screen.Update()
	ctx.GL.Call("viewport", ctx.Screen.ScreenRect.Left(), ctx.Screen.ScreenRect.Top(), ctx.Screen.ScreenRect.Width(), ctx.Screen.ScreenRect.Height())
}
