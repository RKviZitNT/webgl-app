package roommanager

import (
	"fmt"
	"sync"
	"webgl-app/internal/net/player"
	"webgl-app/internal/net/room"
	"webgl-app/internal/utils"
)

type RoomManager struct {
	rooms map[string]*room.Room
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*room.Room),
	}
}

func (rm *RoomManager) CreateRoom(ownerID string, settings room.RoomSettings) (string, error) {
	if settings.MaxPlayers <= 0 {
		return "", fmt.Errorf("invalid max players count")
	}

	roomCode, err := rm.generateRoomCode(6)
	if err != nil {
		return "", err
	}

	room := room.NewRoom(roomCode, settings)
	room.SetOwnerID(ownerID)

	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.rooms[roomCode] = room

	return roomCode, nil
}

func (rm *RoomManager) DeleteRoom(roomCode string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomCode]; !exists {
		return fmt.Errorf("room not found")
	}

	for _, p := range rm.rooms[roomCode].GetPlayers() {
		p.SetRoomID("")
	}
	delete(rm.rooms, roomCode)

	return nil
}

func (rm *RoomManager) JoinRoom(player *player.Player, roomCode string) error {
	rm.mu.Lock()
	room, exists := rm.rooms[roomCode]
	rm.mu.Unlock()

	if !exists {
		return fmt.Errorf("room  not found")
	}

	if err := room.AddPlayer(player); err != nil {
		return err
	}
	player.SetRoomID(roomCode)

	return nil
}

func (rm *RoomManager) KickFromRoom(player *player.Player, roomCode string) error {
	rm.mu.Lock()
	room, exists := rm.rooms[roomCode]
	rm.mu.Unlock()

	if !exists {
		return fmt.Errorf("room not found")
	}

	if err := room.RemovePlayer(player); err != nil {
		return err
	}
	player.SetRoomID("")

	return nil
}

func (rm *RoomManager) GetRoom(roomCode string) (*room.Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomCode]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}

	return room, nil
}

func (rm *RoomManager) generateRoomCode(lenght int) (string, error) {
	const maxAttempts = 100
	for i := 0; i < maxAttempts; i++ {
		code, err := utils.GenerateRandomCode(lenght)
		if err != nil {
			return "", err
		}

		rm.mu.Lock()
		_, exists := rm.rooms[code]
		rm.mu.Unlock()

		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique code")
}
