//go:build js

package clienthandler

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/room"
)

func sendUpdateRoomInfoMsg() {
	msg := message.Message{
		Type: message.UpdateRoomInfoMsg,
		Data: nil,
	}

	sendMessage(msg)
}

func sendUpdatePlayerInfoMsg() {
	msg := message.Message{
		Type: message.UpdatePlayerInfoMsg,
		Data: nil,
	}

	sendMessage(msg)
}

func sendMessage(msg message.Message) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		jsfunc.LogError(fmt.Sprint("JSON error:", err.Error()))
		return
	}
	socket.Call("send", string(jsonData))
}

func updateUi() {
	js.Global().Get("document").Call("getElementById", "lobby_code").Set("textContent", roomInfo.ID)
	js.Global().Get("document").Call("getElementById", "room_status").Set("textContent", fmt.Sprintf("Status: %s", roomInfo.Status))
	js.Global().Get("document").Call("getElementById", "current_players").Set("textContent", roomInfo.PlayersCount)
	js.Global().Get("document").Call("getElementById", "max_players").Set("textContent", roomInfo.MaxPlayers)
	if playerInfo.ID == roomInfo.OwnerId {
		jsfunc.UpdateOwnerControls(true)
		if roomInfo.Status == string(room.Ready) {
			jsfunc.SwitchStartButtonState(true)
		} else {
			jsfunc.SwitchStartButtonState(false)
		}
	} else {
		jsfunc.UpdateOwnerControls(false)
	}
}
