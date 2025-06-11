package wshandler

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/player"
	"webgl-app/internal/net/room"
	"webgl-app/internal/utils"
)

func (ws *WebSocket) handleMessage(player *player.Player, msg *message.Message) {
	switch msg.Type {
	case message.CreateRoomMsg:
		ws.handleCreateRoom(player, msg)
	case message.JoinRoomMsg:
		ws.handleJoinRoom(player, msg)
	case message.LeaveRoomMsg:
		ws.handleLeaveRoom(player)
	case message.StartGameMsg:
		ws.handleStartGame(player)
	case message.EndGameMsg:
		ws.handleEndGame(player)
	case message.UpdateRoomInfoMsg:
		ws.handleUpdateRoomInfo(player)
	case message.UpdatePlayerInfoMsg:
		ws.handleUpdatePlayerInfo(player)
	case message.GameStateMsg:
		ws.handlerGameState(player, msg)
	default:
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "unknown message type",
		})
	}
}

func (ws *WebSocket) handleCreateRoom(player *player.Player, msg *message.Message) {
	var settings room.RoomSettings
	utils.ParseInterfaceToJSON(msg.Data, &settings)

	roomCode, err := ws.rm.CreateRoom(player.ID(), settings)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	err = ws.rm.JoinRoom(player, roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		ws.rm.DeleteRoom(roomCode)
		return
	}

	player.Send(message.Message{
		Type: message.CreateRoomMsg,
		Data: nil,
	})
}

func (ws *WebSocket) handleJoinRoom(player *player.Player, msg *message.Message) {
	roomCode := msg.Data.(string)

	err := ws.rm.JoinRoom(player, roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	room, _ := ws.rm.GetRoom(roomCode)
	player.Send(message.Message{
		Type: message.JoinRoomMsg,
		Data: nil,
	})

	room.Broadcast(message.Message{
		Type: message.PlayerJoinMsg,
		Data: player.GetName(),
	}, player.ID())
}

func (ws *WebSocket) handleLeaveRoom(player *player.Player) {
	roomCode := player.GetRoomID()
	if roomCode == "" {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "player is not in any room",
		})
		return
	}

	room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	if err := ws.rm.KickFromRoom(player, roomCode); err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}
	player.Send(message.Message{
		Type: message.LeaveRoomMsg,
		Data: nil,
	})

	if room.GetOwnerID() == player.ID() {
		ws.rm.DeleteRoom(roomCode)
		room.Broadcast(message.Message{
			Type: message.RoomClosedMsg,
			Data: "the owner has closed the room",
		}, player.ID())
	} else {
		room.Broadcast(message.Message{
			Type: message.PlayerLeftMsg,
			Data: player.ID(),
		}, nil)
	}
}

func (ws *WebSocket) handleStartGame(player *player.Player) {
	roomCode := player.GetRoomID()
	if roomCode == "" {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	ids := make([]string, 0, 2)
	for id := range room.GetPlayers() {
		ids = append(ids, id)
	}

	fightersInfo := make([]message.FighterInfo, 0, len(ids))
	fightersInfo = append(fightersInfo, message.FighterInfo{
		ID:            ids[0],
		CharacterName: "warrior",
		Collider: primitives.NewRect(
			primitives.NewVec2(100, 300),
			primitives.NewVec2(60, 40),
		),
	})
	fightersInfo = append(fightersInfo, message.FighterInfo{
		ID:            ids[1],
		CharacterName: "warrior",
		Collider: primitives.NewRect(
			primitives.NewVec2(500, 300),
			primitives.NewVec2(60, 40),
		),
	})

	room.Broadcast(message.Message{
		Type: message.StartGameMsg,
		Data: message.StartGameData{
			FightersInfo: fightersInfo,
		},
	}, nil)
}

func (ws *WebSocket) handleEndGame(player *player.Player) {
	roomCode := player.GetRoomID()
	if roomCode == "" {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	room.Broadcast(message.Message{
		Type: message.EndGameMsg,
		Data: nil,
	}, nil)
}

func (ws *WebSocket) handleUpdateRoomInfo(player *player.Player) {
	roomCode := player.GetRoomID()
	if roomCode == "" {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	player.Send(message.Message{
		Type: message.UpdateRoomInfoMsg,
		Data: room.RoomInfo(),
	})
}

func (ws *WebSocket) handleUpdatePlayerInfo(player *player.Player) {
	player.Send(message.Message{
		Type: message.UpdatePlayerInfoMsg,
		Data: player.PlayerInfo(),
	})
}

func (ws *WebSocket) handlerGameState(player *player.Player, msg *message.Message) {
	roomCode := player.GetRoomID()
	if roomCode == "" {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: "Player is not in any room",
		})
		return
	}

	room, err := ws.rm.GetRoom(roomCode)
	if err != nil {
		player.Send(message.Message{
			Type: message.ErrorMsg,
			Data: err.Error(),
		})
		return
	}

	room.Broadcast(*msg, player.ID())
}
