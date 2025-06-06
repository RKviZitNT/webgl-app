//go:build js

package clienthandler

import (
	"syscall/js"
	"webgl-app/internal/game/game"
	"webgl-app/internal/game/graphics"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/player"
	"webgl-app/internal/net/room"
)

var (
	socket     js.Value
	roomInfo   room.RoomInfo
	playerInfo player.PlayerInfo
	gm         *game.Game
)

func RegisterCallbacks() {
	c := make(chan struct{}, 0)

	js.Global().Set("createLobby", js.FuncOf(createLobby))
	js.Global().Set("joinLobby", js.FuncOf(joinLobby))
	js.Global().Set("leaveLobby", js.FuncOf(leaveLobby))
	js.Global().Set("startGame", js.FuncOf(startGame))

	connectWebSocket()

	<-c
}

func InitGame(GLCtx *graphics.GLContext) {
	gm = game.NewGame(&socket, GLCtx)
}

func connectWebSocket() {
	socket = js.Global().Get("WebSocket").New("ws://localhost:8080/ws")

	socket.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("WebSocket connected")
		return nil
	}))

	socket.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msg := args[0].Get("data").String()
		// println("Received:", msg)
		handleServerMessage(msg)
		return nil
	}))

	socket.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("WebSocket error")
		return nil
	}))
}

func createLobby(this js.Value, args []js.Value) interface{} {
	msg := message.Message{
		Type: message.CreateRoomMsg,
		Data: room.RoomSettings{
			MaxPlayers:  2,
			NeedPlayers: 2,
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
