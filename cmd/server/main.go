package main

import (
	"log"
	"net/http"
	"path/filepath"
	"webgl-app/internal/net/wshandler"
)

func main() {
	ws := wshandler.NewWebSocket()

	http.Handle("/", http.FileServer(http.Dir(filepath.Join("static"))))
	http.HandleFunc("/ws", ws.WebSocketHandler)

	log.Println("Server started at :8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
