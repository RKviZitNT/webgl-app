//go:build js

package game

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"
	"webgl-app/internal/assetsmanager"
	"webgl-app/internal/game/character"
	"webgl-app/internal/game/fighter"
	"webgl-app/internal/game/level"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/message"
)

type Game struct {
	playerId     string
	running      bool
	socket       *js.Value
	glCtx        *webgl.GLContext
	keys         map[string]bool
	assets       *assetsmanager.AssetsManager
	fighters     map[string]*fighter.Fighter
	characters   map[character.CharacterName]*character.Character
	levels       map[level.LevelName]*level.Level
	currentLevel *level.Level
}

var (
	Direction primitives.Vec2
	Speed     float64
)

func NewGame(socket *js.Value, glCtx *webgl.GLContext) (*Game, error) {
	assets := assetsmanager.NewAssetsManager()
	err := assets.Load(glCtx, assetsmanager.ASrc)
	if err != nil {
		return nil, err
	}

	characters := make(map[character.CharacterName]*character.Character)

	aTypes := []animation.AnimationType{animation.Idle, animation.Walk, animation.Attack1, animation.Attack2, animation.Death, animation.Hurt, animation.Jump}
	warriorAnim, err := animation.CreateAnimations(aTypes, assets.GetMetadata("warrior_anim").String(), assets.GetTexture("warrior"))
	if err != nil {
		return nil, err
	}
	warriorChar := character.NewCharacter(character.Warrior, sprite.NewSprite(assets.GetTexture(string(character.Warrior)), primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(69, 44))))
	warriorChar.SetAnimations(warriorAnim)
	characters[character.Warrior] = warriorChar

	levels := make(map[level.LevelName]*level.Level)
	levels[level.DefaultLevel] = level.NewLevel(level.DefaultLevel, assets.GetTexture("background1"), assets.GetTexture("background1"))

	return &Game{
		socket:     socket,
		glCtx:      glCtx,
		keys:       make(map[string]bool),
		assets:     assets,
		characters: characters,
		levels:     levels,
	}, nil
}

func (g *Game) Start(playerId string, fightersInfo []message.FighterInfo) {
	g.playerId = playerId
	g.running = true

	g.currentLevel = g.levels[level.DefaultLevel]

	Speed = 200.0

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

	g.fighters = make(map[string]*fighter.Fighter)
	for _, fighterInfo := range fightersInfo {
		g.fighters[fighterInfo.ID] = fighter.NewFighter(g.characters[character.CharacterName(fighterInfo.CharacterName)], fighterInfo.Collider)
	}

	g.renderLoop()
}

func (g *Game) Stop() {
	g.running = false
	g.fighters = make(map[string]*fighter.Fighter)
	g.currentLevel = nil
}

func (g *Game) renderLoop() {
	var (
		renderFrame   js.Func
		frameTime     = time.Second / 240
		lastFrameTime time.Time
		elapsedTime   time.Duration
	)

	lastFrameTime = time.Now()

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if g.running {
			currentTime := time.Now()
			deltaTime := currentTime.Sub(lastFrameTime)
			lastFrameTime = currentTime

			g.update(deltaTime)
			g.sendPlayerState()
			g.draw()

			elapsedTime += deltaTime
			if elapsedTime < frameTime {
				time.Sleep(frameTime - elapsedTime)
			}

			js.Global().Call("requestAnimationFrame", renderFrame)
		}
		return nil
	})

	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (g *Game) update(deltaTime time.Duration) {
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
	if g.keys["Escape"] {
		g.sendEndGameMsg()
	}

	Direction.Normalize()
	Direction.MulValue(Speed * deltaTime.Seconds())
	g.fighters[g.playerId].Collider.Move(Direction)

	for _, f := range g.fighters {
		f.Character.Animations[animation.Idle].Update(float64(deltaTime.Milliseconds()))
	}
}

func (g *Game) draw() {
	gl := g.glCtx.GL

	gl.Call("viewport", 0, 0, g.glCtx.CanvasRect.Width(), g.glCtx.CanvasRect.Height())
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	gl.Call("useProgram", g.glCtx.Program)

	g.glCtx.DrawTexture(g.currentLevel.Background, g.glCtx.CanvasRect)

	for _, f := range g.fighters {
		g.glCtx.DrawSprite(f.Character.Animations[f.State].GetCurrentFrame(), f.Collider.Pos, 5)
	}
}

func (g *Game) UpdatePlayersData(fighterInfo message.FighterInfo) {
	fighter := g.fighters[fighterInfo.ID]
	fighter.Collider = fighterInfo.Collider
	fighter.State = animation.AnimationType(fighterInfo.State)
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.FighterInfo{
			ID:            g.playerId,
			CharacterName: string(character.Warrior),
			Collider:      g.fighters[g.playerId].Collider,
		},
	}

	g.sendMessage(msg)
}

func (g *Game) sendEndGameMsg() {
	msg := message.Message{
		Type: message.EndGameMsg,
		Data: nil,
	}

	g.sendMessage(msg)
}

func (g *Game) sendMessage(msg message.Message) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		js.Global().Get("console").Call("error", fmt.Sprint("JSON error:", err.Error()))
		return
	}
	g.socket.Call("send", string(jsonData))
}
