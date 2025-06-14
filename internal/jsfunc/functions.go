//go:build js

package jsfunc

import (
	"syscall/js"
)

type Screen string

const (
	MainMenuScreen      Screen = "main_menu"
	LobbyScreen         Screen = "lobby"
	LobbyConnectScreen  Screen = "lobby_connect"
	LoadingScreenScreen Screen = "loading_screen"
	GameScreenScreen    Screen = "game_screen"
)

func ShowScreen(screen Screen) {
	js.Global().Call("showScreen", string(screen))
}

func UpdateOwnerControls(isOwner bool) {
	js.Global().Call("updateOwnerControls", isOwner)
}

func SwitchStartButtonState(isEnabled bool) {
	js.Global().Call("switchStartButtonState", isEnabled)
}

func SetLoadingProgress(progress float64, message string) {
	js.Global().Call("setLoadingProgress", progress, message)
}

func LoadFile(path string) js.Value {
	return js.Global().Call("loadFile", path)
}

func LoadImage(path string) js.Value {
	return js.Global().Call("loadImage", path)
}
