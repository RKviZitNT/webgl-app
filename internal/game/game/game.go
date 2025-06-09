//go:build js

package game

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"webgl-app/internal/assetsmanager"
	"webgl-app/internal/game/character"
	"webgl-app/internal/game/fighter"
	"webgl-app/internal/game/level"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/message"
)

type Game struct {
	socket   *js.Value
	glCtx    *webgl.GLContext
	keys     map[string]bool
	assets   *assetsmanager.AssetsManager
	level    *level.Level
	fighters map[string]*fighter.Fighter
}

var (
	PlayerId  string
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

func (g *Game) Start(playerId string, fightersInfo []message.FighterInfo) {
	PlayerId = playerId
	Speed = 2

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

	g.assets = assetsmanager.NewAssetsManager()
	js.Global().Call("showScreen", "loading_screen")
	err := g.assets.Load(g.glCtx, assetsmanager.AssetsSrc)
	if err != nil {
		g.Stop("Failed to load assets!")
	}
	js.Global().Call("showScreen", "game_screen")

	g.level = level.NewLevel(g.assets.GetTexture("background1"), g.assets.GetTexture("background1"))

	character.Characters = make(map[character.CharacterName]*character.Character)
	character.Characters[character.Warrior] = character.NewCharacter(character.Warrior, sprite.NewSprite(g.assets.GetTexture(string(character.Warrior)), primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(69, 44))))

	g.fighters = make(map[string]*fighter.Fighter)
	for _, fighterInfo := range fightersInfo {
		g.fighters[fighterInfo.ID] = fighter.NewFighter(fighterInfo.CharacterName, fighterInfo.Collider)
	}

	g.renderLoop()
}

func (g *Game) Stop(exitMsg string) {
	if exitMsg != "" {
		js.Global().Call("log", exitMsg)
	}
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
	g.fighters[PlayerId].Collider.Move(Direction.Normalize())
}

func (g *Game) draw() {
	gl := g.glCtx.GL

	gl.Call("viewport", 0, 0, g.glCtx.CanvasRect.Width(), g.glCtx.CanvasRect.Height())
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	gl.Call("useProgram", g.glCtx.Program)

	g.glCtx.DrawTexture(g.level.Background, g.glCtx.CanvasRect)

	for _, f := range g.fighters {
		g.glCtx.DrawSprite(f.Character.Sprite, f.Collider.Pos, 5)
	}
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.FighterInfo{
			ID:            PlayerId,
			CharacterName: string(character.Warrior),
			Collider:      g.fighters[PlayerId].Collider,
		},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		js.Global().Get("console").Call("error", fmt.Sprint("JSON error:", err.Error()))
		return
	}
	g.socket.Call("send", string(jsonData))
}

func (g *Game) UpdatePlayersData(fighterInfo message.FighterInfo) {
	g.fighters[fighterInfo.ID].Collider = fighterInfo.Collider
}
