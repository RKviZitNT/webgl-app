package room

import (
	"fmt"
	"sync"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/player"
)

type RoomStatus string

const (
	Waiting RoomStatus = "Waiting"
	Ready   RoomStatus = "Ready"
	InGame  RoomStatus = "In game"
)

type RoomSettings struct {
	MaxPlayers  int `json:"max_players"`
	NeedPlayers int `json:"need_players"`
}

type Room struct {
	id       string
	status   RoomStatus
	settings RoomSettings
	players  map[string]*player.Player
	ownerID  string
	mu       sync.Mutex
}

type RoomInfo struct {
	ID       string       `json:"id"`
	Status   RoomStatus   `json:"status"`
	OwnerId  string       `json:"owner_id"`
	Settings RoomSettings `json:"settings"`
}

func NewRoom(roomCode string, settings RoomSettings) *Room {
	return &Room{
		id:       roomCode,
		status:   Waiting,
		settings: settings,
		players:  make(map[string]*player.Player),
		ownerID:  "",
	}
}

func (r *Room) ID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.id
}

func (r *Room) RoomInfo() *RoomInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	return &RoomInfo{
		ID:       r.id,
		Status:   r.status,
		OwnerId:  r.ownerID,
		Settings: r.settings,
	}
}

func (r *Room) SetStatus(status RoomStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status = status
}

func (r *Room) GetStatus() RoomStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.status
}

func (r *Room) GetSettings() RoomSettings {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.settings
}

func (r *Room) AddPlayer(player *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status == InGame {
		return fmt.Errorf("there is a game going on in the room now")
	}
	if len(r.players) > r.settings.MaxPlayers {
		return fmt.Errorf("room is full")
	}

	r.players[player.ID()] = player
	r.updateStatus()

	return nil
}

func (r *Room) RemovePlayer(player *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.players[player.ID()]
	if !exists {
		return fmt.Errorf("player not found in room")
	}

	delete(r.players, player.ID())
	r.updateStatus()

	return nil
}

func (r *Room) GetPlayers() map[string]*player.Player {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.players
}

func (r *Room) GetPlayer(id string) (*player.Player, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	player, exists := r.players[id]
	if !exists {
		return nil, fmt.Errorf("player not found in room")
	}

	return player, nil
}

func (r *Room) GetPlayersCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.players)
}

func (r *Room) SetOwnerID(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ownerID = id
}

func (r *Room) GetOwnerID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.ownerID
}

func (r *Room) Broadcast(msg message.Message, excludedPlayerId interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.players {
		if excludedPlayerId == nil || excludedPlayerId.(string) != p.ID() {
			p.Send(msg)
		}
	}
}

func (r *Room) updateStatus() {
	if len(r.players) >= r.settings.NeedPlayers {
		r.status = Ready
	} else {
		r.status = Waiting
	}
}
