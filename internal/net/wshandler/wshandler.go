package wshandler

import (
	"encoding/json"
	"log"
	"net/http"
	"webgl-app/internal/net/message"
	"webgl-app/internal/net/player"
	"webgl-app/internal/net/roommanager"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocket struct {
	upgrader websocket.Upgrader
	rm       roommanager.RoomManager
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		rm: *roommanager.NewRoomManager(),
	}
}

func (ws *WebSocket) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	player := player.NewPlayer(conn, "Player")
	log.Printf("Player %s connected", player.ID())

	defer func() {
		ws.handleEndGame(player)
		ws.handleLeaveRoom(player)
		conn.Close()
	}()

	for {
		_, rdmsg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Player %s disconnected", player.ID())
			break
		}

		var msg message.Message
		if err := json.Unmarshal(rdmsg, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		go ws.handleMessage(player, msg)
	}
}
