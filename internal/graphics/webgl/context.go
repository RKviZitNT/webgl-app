//go:build js

package webgl

import (
	"syscall/js"
	"webgl-app/internal/resourcemanager"
)

type GLContext struct {
	GL js.Value
}

var (
	VertexShaderSrc   string
	FragmentShaderSrc string
)

func InitWebGL(canvasID string) (*GLContext, error) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", canvasID)

	gl := canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		return nil, js.Error{Value: js.ValueOf("WebGL not supported")}
	}

	done := make(chan struct{})
	go func() {
		VertexShaderSrc = resourcemanager.LoadShader("shaders/vertex.glsl")
		FragmentShaderSrc = resourcemanager.LoadShader("shaders/fragment.glsl")
		close(done)
	}()
	<-done

	return &GLContext{GL: gl}, nil
}
