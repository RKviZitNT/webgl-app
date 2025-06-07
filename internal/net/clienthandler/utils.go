//go:build js

package clienthandler

import (
	"encoding/json"
	"fmt"
	"syscall/js"
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
		println("JSON error:", err.Error())
		return
	}
	socket.Call("send", string(jsonData))
}

func updateUi() {
	js.Global().Get("document").Call("getElementById", "lobby_code").Set("textContent", roomInfo.ID)
	js.Global().Get("document").Call("getElementById", "room_status").Set("textContent", fmt.Sprintf("Status: %s", roomInfo.Status))
	if playerInfo.ID == roomInfo.OwnerId {
		js.Global().Call("updateOwnerControls", true)
		if roomInfo.Status == room.Ready {
			js.Global().Call("switchStartButtonState", true)
		} else {
			js.Global().Call("switchStartButtonState", false)
		}
	} else {
		js.Global().Call("updateOwnerControls", false)
	}
}
