//go:build js

package clienthandler

import (
	"encoding/json"
	"fmt"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/message"
	"webgl-app/internal/utils"
)

func handleServerMessage(raw string) {
	var msg message.Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		jsfunc.LogError(fmt.Sprint("Failed to parse message: ", err.Error()))
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
	case message.EndGameMsg:
		handleEndGame(msg.Data)
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
		jsfunc.LogError(fmt.Sprint("Unknown message type: ", msg.Type))
	}
}

func handleCreateRoom(data interface{}) {
	sendUpdateRoomInfoMsg()
	sendUpdatePlayerInfoMsg()
	jsfunc.ShowScreen(jsfunc.LobbyScreen)
}

func handleJoinRoom(data interface{}) {
	sendUpdateRoomInfoMsg()
	sendUpdatePlayerInfoMsg()
	jsfunc.ShowScreen(jsfunc.LobbyScreen)
}

func handleLeaveRoom(data interface{}) {
	jsfunc.ShowScreen(jsfunc.MainMenuScreen)
}

func handleStartGame(data interface{}) {
	sendUpdateRoomInfoMsg()
	jsfunc.ShowScreen(jsfunc.GameScreenScreen)

	var gameData message.StartGameData
	utils.ParseInterfaceToJSON(data, &gameData)

	gm.Stop()
	go gm.Start(playerInfo.ID, gameData.FightersInfo)
}

func handleEndGame(data interface{}) {
	gm.Stop()
	sendUpdateRoomInfoMsg()
	jsfunc.ShowScreen(jsfunc.LobbyScreen)
}

func handleUpdateRoomInfo(data interface{}) {
	if err := utils.ParseInterfaceToJSON(data, &roomInfo); err != nil {
		jsfunc.LogError(err.Error())
		return
	}

	updateUi()
}

func handleUpdatePlayerInfo(data interface{}) {
	if err := utils.ParseInterfaceToJSON(data, &playerInfo); err != nil {
		jsfunc.LogError(err.Error())
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
	gm.Stop()
	jsfunc.ShowScreen(jsfunc.MainMenuScreen)
}

func handleGameState(data interface{}) {
	var playerState message.FighterInfo
	utils.ParseInterfaceToJSON(data, &playerState)
	gm.UpdatePlayersData(playerState)
}

func handleError(data interface{}) {
	gm.Stop()
	jsfunc.ShowScreen(jsfunc.MainMenuScreen)
}
