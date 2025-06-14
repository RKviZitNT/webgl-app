package message

import (
	"webgl-app/internal/graphics/primitives"
)

type MessageType string

const (
	ErrorMsg            MessageType = "error"
	CreateRoomMsg       MessageType = "create_room"
	JoinRoomMsg         MessageType = "join_room"
	LeaveRoomMsg        MessageType = "leave_room"
	StartGameMsg        MessageType = "start_game"
	EndGameMsg          MessageType = "end_game"
	UpdateRoomInfoMsg   MessageType = "update_room_info"
	UpdatePlayerInfoMsg MessageType = "update_player_info"
	PlayerLeftMsg       MessageType = "player_left"
	PlayerJoinMsg       MessageType = "player_join"
	RoomClosedMsg       MessageType = "room_closed"
	GameStateMsg        MessageType = "game_state"
)

type Message struct {
	Type MessageType
	Data interface{}
}

type PlayerInfo struct {
	ID   string
	Name string
}

type RoomInfo struct {
	ID           string
	Status       string
	OwnerId      string
	PlayersCount int
	MaxPlayers   int
	NeedPlayers  int
}

type FighterInfo struct {
	ID            string
	CharacterName string
	State         string
	Collider      primitives.Rect
}

type StartGameData struct {
	FightersInfo []FighterInfo
}
