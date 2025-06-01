//go:build js

package main

import (
	"webgl-app/internal/net/clienthandler"
)

func main() {
	clienthandler.RegisterCallbacks()
	select {}
}
