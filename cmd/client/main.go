//go:build js

package main

import (
	"webgl-app/internal/game/graphics"
	"webgl-app/internal/net/clienthandler"
)

func main() {
	GLContext, err := graphics.InitWebGL("game_canvas")
	if err != nil {
		println(err)
		return
	}

	clienthandler.InitGame(GLContext)
	clienthandler.RegisterCallbacks()

	select {}
}
