//go:build js

package main

import (
	"webgl-app/internal/config"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/clienthandler"
)

func main() {
	var err error

	jsfunc.LogInfo(" ----- Loading configs ----- ")
	err = config.LoadConfigs("config_shaders.json", "config_assets.json")
	if err != nil {
		jsfunc.LogError(err.Error())
	}

	jsfunc.LogInfo(" ----- Init WebGL ----- ")
	GLContext, err := webgl.NewWebGLContext("game_canvas")
	if err != nil {
		jsfunc.LogError(err.Error())
		return
	}
	err = GLContext.InitWebGL()
	if err != nil {
		jsfunc.LogError(err.Error())
		return
	}

	err = clienthandler.InitGame(GLContext)
	if err != nil {
		jsfunc.LogError(err.Error())
		return
	}

	clienthandler.RegisterCallbacks()

	select {}
}
