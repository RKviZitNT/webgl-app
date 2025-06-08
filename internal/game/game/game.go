//go:build js

package game

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/message"
	"webgl-app/internal/resourcemanager"
)

type Game struct {
	socket     *js.Value
	glCtx      *webgl.GLContext
	keys       map[string]bool
	characters map[string]*character.Character

	backg   *texture.Texture
	texture *texture.Texture
	sprite  *sprite.Sprite
}

var (
	PlayerId string

	Direction primitives.Vec2
	Speed     float64
)

func NewGame(socket *js.Value, glCtx *webgl.GLContext) *Game {
	return &Game{
		socket: socket,
		glCtx:  glCtx,
		keys:   make(map[string]bool),
	}
}

func (g *Game) Start(playerId string, chars map[string]*character.Character) error {
	PlayerId = playerId
	Speed = 2
	g.characters = chars

	js.Global().Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := args[0].Get("code").String()
		g.keys[code] = true

		return nil
	}))
	js.Global().Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := args[0].Get("code").String()
		g.keys[code] = false

		return nil
	}))

	resourcemanager.LoadImage("assets/sprites/warrior/spritesheet.png",
		func(img js.Value) {
			g.texture = texture.NewTexture(g.glCtx.GL, img)
			g.sprite = sprite.NewSprite(g.texture, primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(69, 44)))
		},
		func(err error) {
			js.Global().Get("console").Call("error", fmt.Sprint("Load texture error:", err.Error()))
		})
	resourcemanager.LoadImage("assets/images/backgrounds/background1.jpg",
		func(img js.Value) {
			g.backg = texture.NewTexture(g.glCtx.GL, img)
		},
		func(err error) {
			js.Global().Get("console").Call("error", fmt.Sprint("Load texture error:", err.Error()))
		})

	g.renderLoop()
	return nil
}

func (g *Game) renderLoop() {
	var renderFrame js.Func

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.update()
		g.sendPlayerState()
		g.draw()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (g *Game) update() {
	Direction = primitives.NewVec2(0, 0)

	if g.keys["ArrowRight"] || g.keys["KeyD"] {
		Direction.AddVec2(primitives.NewVec2(1, 0))
	}
	if g.keys["ArrowLeft"] || g.keys["KeyA"] {
		Direction.AddVec2(primitives.NewVec2(-1, 0))
	}
	if g.keys["ArrowUp"] || g.keys["KeyW"] {
		Direction.AddVec2(primitives.NewVec2(0, -1))
	}
	if g.keys["ArrowDown"] || g.keys["KeyS"] {
		Direction.AddVec2(primitives.NewVec2(0, 1))
	}

	Direction.MulValue(Speed)
	g.characters[PlayerId].HitBox.Move(Direction.Normalize())
}

func (g *Game) draw() {
	gl := g.glCtx.GL

	gl.Call("viewport", 0, 0, g.glCtx.CanvasRect.Width(), g.glCtx.CanvasRect.Height())
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	gl.Call("useProgram", g.glCtx.Program)

	if g.backg != nil {
		g.glCtx.DrawTexture(g.backg, g.glCtx.CanvasRect)
	}
	if g.sprite != nil {
		g.glCtx.DrawSprite(g.sprite, primitives.NewVec2(200, 500), 5)
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
		js.Global().Get("console").Call("error", fmt.Sprint("JSON error:", err.Error()))
		return
	}
	g.socket.Call("send", string(jsonData))
}

func (g *Game) UpdatePlayersData(playerState message.PlayerState) {
	g.characters[playerState.Id] = &playerState.Data
}
