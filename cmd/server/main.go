package main

import (
	"log"
	"net/http"
	"path/filepath"
	"webgl-app/internal/net/wshandler"
)

func main() {
	ws := wshandler.NewWebSocket()

	http.Handle("/", http.FileServer(http.Dir(filepath.Join("build", "static"))))
	http.HandleFunc("/ws", ws.WebSocketHandler)

	log.Println("Server startes at http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
