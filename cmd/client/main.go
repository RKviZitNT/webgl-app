//go:build js

package main

import (
	"syscall/js"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/clienthandler"
)

func main() {
	GLContext, err := webgl.NewWebGLCtx("game_canvas")
	if err != nil {
		js.Global().Get("console").Call("error", err.Error())
		return
	}
	err = GLContext.InitWebGL()
	if err != nil {
		js.Global().Get("console").Call("error", err.Error())
		return
	}

	clienthandler.InitGame(GLContext)
	clienthandler.RegisterCallbacks()

	select {}
}
