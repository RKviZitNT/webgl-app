//go:build js

package clienthandler

import (
	"syscall/js"
	"webgl-app/internal/config"
	"webgl-app/internal/game/game"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/room"
)

var (
	socket     js.Value
	roomInfo   message.RoomInfo
	playerInfo message.PlayerInfo
	gm         *game.Game
)

func RegisterCallbacks() {
	c := make(chan struct{}, 0)

	js.Global().Set("createLobby", js.FuncOf(createLobby))
	js.Global().Set("joinLobby", js.FuncOf(joinLobby))
	js.Global().Set("leaveLobby", js.FuncOf(leaveLobby))
	js.Global().Set("startGame", js.FuncOf(startGame))

	jsfunc.LogInfo(" ----- Connecting to WebSocket ----- ")

	jsfunc.SetLoadingProgress(100, "Initialization...")
	connectWebSocket()

	jsfunc.ShowScreen("main_menu")

	<-c
}

func InitGame(GLCtx *webgl.GLContext) error {
	var err error
	gm, err = game.NewGame(&socket, GLCtx)
	if err != nil {
		return err
	}
	return nil
}

func connectWebSocket() {
	if !socket.IsUndefined() {
		socket.Call("close")
	}

	wsURL := js.Global().Call("getWebSocketURL").String()
	socket = js.Global().Get("WebSocket").New(wsURL)

	socket.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsfunc.LogInfo("WebSocket connected")
		return nil
	}))

	socket.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handleServerMessage(args[0].Get("data").String())
		return nil
	}))

	socket.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsfunc.LogError("WebSocket connection error")
		socket.Call("close")
		gm.Stop()
		return nil
	}))

	socket.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsfunc.LogInfo("WebSocket connection closed")
		return nil
	}))
}

func createLobby(this js.Value, args []js.Value) interface{} {
	var needPlayers int
	if config.ProgramConfig.Debug {
		needPlayers = 1
	} else {
		needPlayers = 2
	}

	msg := message.Message{
		Type: message.CreateRoomMsg,
		Data: room.RoomSettings{
			MaxPlayers:  2,
			NeedPlayers: needPlayers,
		},
	}

	sendMessage(msg)
	return nil
}

func joinLobby(this js.Value, args []js.Value) interface{} {
	roomCode := js.Global().Get("document").Call("getElementById", "room_code").Get("value").String()

	msg := message.Message{
		Type: message.JoinRoomMsg,
		Data: roomCode,
	}

	sendMessage(msg)
	return nil
}

func leaveLobby(this js.Value, args []js.Value) interface{} {
	msg := message.Message{
		Type: message.LeaveRoomMsg,
		Data: nil,
	}

	sendMessage(msg)
	return nil
}

func startGame(this js.Value, args []js.Value) interface{} {
	msg := message.Message{
		Type: message.StartGameMsg,
		Data: nil,
	}

	sendMessage(msg)
	return nil
}
