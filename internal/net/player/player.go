package player

import (
	"sync"
	"webgl-app/internal/net/message"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	conn   *websocket.Conn
	id     string
	name   string
	roomID string
	mu     sync.RWMutex
}

func NewPlayer(conn *websocket.Conn, name string) *Player {
	return &Player{
		conn:   conn,
		id:     uuid.New().String(),
		name:   name,
		roomID: "",
	}
}

func (p *Player) ID() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.id
}

func (p *Player) PlayerInfo() *message.PlayerInfo {
	p.mu.Lock()
	defer p.mu.Unlock()

	return &message.PlayerInfo{
		ID:   p.id,
		Name: p.name,
	}
}

func (p *Player) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.name = name
}

func (p *Player) GetName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.name
}

func (p *Player) SetRoomID(roomID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.roomID = roomID
}

func (p *Player) GetRoomID() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.roomID
}

func (p *Player) Send(msg message.Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.conn.WriteJSON(msg)
}
