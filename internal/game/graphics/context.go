//go:build js

package graphics

import (
	"syscall/js"
)

type GLContext struct {
	GL js.Value
}

func InitWebGL(canvasID string) (*GLContext, error) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", canvasID)

	gl := canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		return nil, js.Error{Value: js.ValueOf("WebGL not supported")}
	}

	LoadShaders()
	return &GLContext{GL: gl}, nil
}
