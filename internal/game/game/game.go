//go:build js

package game

import (
	"encoding/json"
	"syscall/js"
	"webgl-app/internal/game/character"
	"webgl-app/internal/game/primitives"
	"webgl-app/internal/net/message"
)

type Game struct {
	ctx        js.Value
	socket     *js.Value
	canvas     js.Value
	keys       map[string]bool
	characters map[string]*character.Character
}

var (
	PlayerId string
	Width    float64
	Height   float64
)

func NewGame(sk *js.Value, chars map[string]*character.Character) *Game {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "game_canvas")

	Width = js.Global().Get("innerWidth").Float()
	Height = js.Global().Get("innerHeight").Float()
	canvas.Set("width", Width)
	canvas.Set("height", Height)

	ctx := canvas.Call("getContext", "2d")

	return &Game{
		ctx:        ctx,
		socket:     sk,
		canvas:     canvas,
		keys:       make(map[string]bool),
		characters: chars,
	}
}

func (g *Game) Start(playerId string) {
	PlayerId = playerId
	// Подписка на клавиши (пример)
	js.Global().Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		g.keys[key] = true
		return nil
	}))

	js.Global().Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		g.keys[key] = false
		return nil
	}))

	// Игровой цикл
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("resizeCanvas")
		g.Update()
		g.sendPlayerState()
		g.Draw()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (g *Game) Update() {
	// Движение игрока (пример)
	if g.keys["ArrowRight"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 1, Y: 0})
	}
	if g.keys["ArrowLeft"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: -1, Y: 0})
	}
	if g.keys["ArrowUp"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 0, Y: -1})
	}
	if g.keys["ArrowDown"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 0, Y: 1})
	}
}

func (g *Game) Draw() {
	// Очистка экрана
	g.ctx.Call("clearRect", 0, 0, g.canvas.Get("width"), g.canvas.Get("height"))

	// Рисуем игрока (квадрат)
	g.ctx.Set("fillStyle", "blue")
	for id, _ := range g.characters {
		g.ctx.Call("fillRect", g.characters[id].HitBox.Pos.X, g.characters[id].HitBox.Pos.Y, 50, 50)
	}
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.PlayerState{
			Id:   PlayerId,
			Data: *g.characters[PlayerId],
		},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		println("JSON error:", err.Error())
		return
	}
	g.socket.Call("send", string(jsonData))
}

func (g *Game) UpdatePlayersData(playerState message.PlayerState) {
	g.characters[playerState.Id] = &playerState.Data
}
