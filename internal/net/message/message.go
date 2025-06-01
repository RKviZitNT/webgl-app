package message

import "webgl-app/internal/game/character"

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
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type StartGameData struct {
	Data map[string]*character.Character
}

type PlayerState struct {
	Id   string
	Data character.Character
}
