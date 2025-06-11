package message

import "webgl-app/internal/graphics/primitives"

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

type PlayerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomInfo struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	OwnerId      string `json:"owner_id"`
	PlayersCount int    `json:"players_count"`
	MaxPlayers   int    `json:"max_players"`
	NeedPlayers  int    `json:"need_players"`
}

type FighterInfo struct {
	ID            string
	CharacterName string
	Collider      primitives.Rect
	State         string
}

type StartGameData struct {
	FightersInfo []FighterInfo
}
