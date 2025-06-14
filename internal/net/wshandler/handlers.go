package wshandler

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/player"
	"webgl-app/internal/net/room"
	"webgl-app/internal/utils"
)

func (ws *WebSocket) handleMessage(_player *player.Player, msg message.Message) {
	switch msg.Type {
	case message.CreateRoomMsg:
		ws.handleCreateRoom(_player, msg)
	case message.JoinRoomMsg:
		ws.handleJoinRoom(_player, msg)
	case message.LeaveRoomMsg:
		ws.handleLeaveRoom(_player)
	case message.StartGameMsg:
		ws.handleStartGame(_player)
	case message.EndGameMsg:
		ws.handleEndGame(_player)
	case message.UpdateRoomInfoMsg:
		ws.handleUpdateRoomInfo(_player)
	case message.UpdatePlayerInfoMsg:
		ws.handleUpdatePlayerInfo(_player)
	case message.GameStateMsg:
		ws.handlerGameState(_player, msg)
	default:
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "unknown message type",
		})
	}
}

func (ws *WebSocket) handleCreateRoom(_player *player.Player, msg message.Message) {
	var settings room.RoomSettings
	utils.ParseInterfaceToJSON(msg.Data, &settings)

	roomCode, err := ws.rm.CreateRoom(_player.ID(), settings)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	err = ws.rm.JoinRoom(_player, roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		ws.rm.DeleteRoom(roomCode)
		return
	}

	_player.Send(message.Message{
		Type: message.CreateRoomMsg,
		Data: nil,
	})
}

func (ws *WebSocket) handleJoinRoom(_player *player.Player, msg message.Message) {
	roomCode := msg.Data.(string)

	err := ws.rm.JoinRoom(_player, roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	_room, _ := ws.rm.GetRoom(roomCode)
	_player.Send(message.Message{
		Type: message.JoinRoomMsg,
		Data: nil,
	})

	_room.Broadcast(message.Message{
		Type: message.PlayerJoinMsg,
		Data: _player.GetName(),
	}, _player.ID())
}

func (ws *WebSocket) handleLeaveRoom(_player *player.Player) {
	roomCode := _player.GetRoomID()
	if roomCode == "" {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "player is not in any room",
		})
		return
	}

	_room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	if err := ws.rm.KickFromRoom(_player, roomCode); err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}
	_player.Send(message.Message{
		Type: message.LeaveRoomMsg,
		Data: nil,
	})

	if _room.GetOwnerID() == _player.ID() {
		ws.rm.DeleteRoom(roomCode)
		_room.Broadcast(message.Message{
			Type: message.RoomClosedMsg,
			Data: "the owner has closed the room",
		}, _player.ID())
	} else {
		_room.Broadcast(message.Message{
			Type: message.PlayerLeftMsg,
			Data: _player.ID(),
		}, nil)
	}
}

func (ws *WebSocket) handleStartGame(_player *player.Player) {
	roomCode := _player.GetRoomID()
	if roomCode == "" {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	_room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	positions := []primitives.Rect{
		primitives.NewRect(100, 550, 100, 160),
		primitives.NewRect(500, 550, 100, 160),
	}
	ids := make([]string, 0, 2)
	for id := range _room.GetPlayers() {
		ids = append(ids, id)
	}

	fightersInfo := make([]message.FighterInfo, 0, len(ids))
	for i, id := range ids {
		fightersInfo = append(fightersInfo, message.FighterInfo{
			ID:            id,
			CharacterName: "warrior",
			Collider:      positions[i],
		})
	}

	_room.Broadcast(message.Message{
		Type: message.StartGameMsg,
		Data: message.StartGameData{
			FightersInfo: fightersInfo,
		},
	}, nil)
}

func (ws *WebSocket) handleEndGame(_player *player.Player) {
	roomCode := _player.GetRoomID()
	if roomCode == "" {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	_room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	_room.Broadcast(message.Message{
		Type: message.EndGameMsg,
		Data: nil,
	}, nil)
}

func (ws *WebSocket) handleUpdateRoomInfo(_player *player.Player) {
	roomCode := _player.GetRoomID()
	if roomCode == "" {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	_room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	_player.Send(message.Message{
		Type: message.UpdateRoomInfoMsg,
		Data: _room.RoomInfo(),
	})
}

func (ws *WebSocket) handleUpdatePlayerInfo(_player *player.Player) {
	_player.Send(message.Message{
		Type: message.UpdatePlayerInfoMsg,
		Data: _player.PlayerInfo(),
	})
}

func (ws *WebSocket) handlerGameState(_player *player.Player, msg message.Message) {
	roomCode := _player.GetRoomID()
	if roomCode == "" {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	_room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		_player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	_room.Broadcast(msg, _player.ID())
}
