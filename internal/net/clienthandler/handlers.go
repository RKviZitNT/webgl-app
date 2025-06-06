//go:build js

package clienthandler

import (
	"encoding/json"
	"syscall/js"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/utils"
)

func handleServerMessage(raw string) {
	var msg message.Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		println("Failed to parse message: ", err.Error())
		return
	}

	switch msg.Type {
	case message.CreateRoomMsg:
		handleCreateRoom(msg.Data)
	case message.JoinRoomMsg:
		handleJoinRoom(msg.Data)
	case message.LeaveRoomMsg:
		handleLeaveRoom(msg.Data)
	case message.StartGameMsg:
		handleStartGame(msg.Data)
	case message.UpdateRoomInfoMsg:
		handleUpdateRoomInfo(msg.Data)
	case message.UpdatePlayerInfoMsg:
		handleUpdatePlayerInfo(msg.Data)
	case message.PlayerJoinMsg:
		handlePlayerJoin(msg.Data)
	case message.PlayerLeftMsg:
		handlePlayerLeft(msg.Data)
	case message.RoomClosedMsg:
		handleRoomClosed(msg.Data)
	case message.GameStateMsg:
		handleGameState(msg.Data)
	case message.ErrorMsg:
		handleError(msg.Data)
	default:
		println("Unknown message type: ", msg.Type)
	}
}

func handleCreateRoom(data interface{}) {
	sendUpdateRoomInfoMsg()
	sendUpdatePlayerInfoMsg()
	js.Global().Call("showScreen", "lobby")
}

func handleJoinRoom(data interface{}) {
	sendUpdateRoomInfoMsg()
	sendUpdatePlayerInfoMsg()
	js.Global().Call("showScreen", "lobby")
}

func handleLeaveRoom(data interface{}) {
	js.Global().Call("showScreen", "main_menu")
}

func handleStartGame(data interface{}) {
	js.Global().Call("showScreen", "game_screen")

	var gameData message.StartGameData
	utils.ReadStruct(data, &gameData)

	gm.Start(playerInfo.ID, gameData.Data)
}

func handleUpdateRoomInfo(data interface{}) {
	if err := utils.ReadStruct(data, &roomInfo); err != nil {
		println(err)
		return
	}

	updateUi()
}

func handleUpdatePlayerInfo(data interface{}) {
	if err := utils.ReadStruct(data, &playerInfo); err != nil {
		println(err)
		return
	}

	updateUi()
}

func handlePlayerLeft(data interface{}) {
	sendUpdateRoomInfoMsg()
}

func handlePlayerJoin(data interface{}) {
	sendUpdateRoomInfoMsg()
}

func handleRoomClosed(data interface{}) {
	js.Global().Call("showScreen", "main_menu")
}

func handleGameState(data interface{}) {
	var playerState message.PlayerState
	utils.ReadStruct(data, &playerState)
	gm.UpdatePlayersData(playerState)
}

func handleError(data interface{}) {

}
