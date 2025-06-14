//go:build js

package game

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"
	"webgl-app/internal/assetsmanager"
	"webgl-app/internal/config"
	"webgl-app/internal/game/character"
	"webgl-app/internal/game/fighter"
	"webgl-app/internal/game/level"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/message"
)

type Game struct {
	players      []string
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
	jsfunc.LogInfo(" ----- Loading assets ----- ")
	assets := assetsmanager.NewAssetsManager()
	err := assets.Load(glCtx, "assets_config.json")
	if err != nil {
		return nil, err
	}

	characters := make(map[character.CharacterName]*character.Character)

	warriorAnim, err := animation.NewAnimationsSet(assets.GetMetadata("warrior_anim").String(), assets.GetTexture("warrior"), 5.5, primitives.NewVec2(-98, -75), primitives.NewVec2(-180, -75))
	if err != nil {
		return nil, err
	}
	warriorChar := character.NewCharacter(character.Warrior, warriorAnim["idle"].GetCurrentFrame())
	warriorChar.SetAnimations(warriorAnim)
	characters[character.Warrior] = warriorChar

	levels := make(map[level.LevelName]*level.Level)
	levels[level.DefaultLevel] = level.NewLevel(level.DefaultLevel, webgl.NewSprite(assets.GetTexture("background2"), nil, 1, primitives.NewVec2(0, 0), primitives.NewVec2(0, 0)))

	return &Game{
		players:    make([]string, 2),
		socket:     socket,
		glCtx:      glCtx,
		keys:       make(map[string]bool),
		assets:     assets,
		characters: characters,
		levels:     levels,
	}, nil
}

func (g *Game) Start(playerId string, fightersInfo []message.FighterInfo) {
	g.players[0] = playerId

	g.currentLevel = g.levels[level.DefaultLevel]

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
		if g.players[0] != fighterInfo.ID {
			g.players[1] = fighterInfo.ID
		}
		g.fighters[fighterInfo.ID] = fighter.NewFighter(g.characters[character.CharacterName(fighterInfo.CharacterName)], fighterInfo.Collider)
	}

	g.renderLoop()
}

func (g *Game) Stop() {
	g.running = false
	g.keys = make(map[string]bool)
}

func (g *Game) renderLoop() {
	var (
		renderFrame   js.Func
		frameTime     = time.Second / time.Duration(config.ProgramConf.Window.FrameRate)
		lastFrameTime time.Time
		elapsedTime   time.Duration
	)

	lastFrameTime = time.Now()

	g.running = true
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if g.running {
			currentTime := time.Now()
			deltaTime := currentTime.Sub(lastFrameTime)
			lastFrameTime = currentTime

			g.update(deltaTime)
			g.sendPlayerState()
			g.draw()

			elapsedTime = time.Now().Sub(currentTime)
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
	if g.keys["Escape"] {
		g.sendEndGameMsg()
	}

	if g.players[1] != "" {
		g.fighters[g.players[0]].HandleSpecular(g.fighters[g.players[1]].Collider.Center())
		g.fighters[g.players[1]].HandleSpecular(g.fighters[g.players[0]].Collider.Center())
	}

	g.fighters[g.players[0]].Update(g.keys, deltaTime)
}

func (g *Game) draw() {
	g.glCtx.RenderSprite(g.currentLevel.Background, g.glCtx.Screen.BaseScreenRect, false)

	f1, f2 := g.fighters[g.players[0]], g.fighters[g.players[1]]

	if g.players[1] != "" {
		f2.Draw(g.glCtx)
	}
	f1.Draw(g.glCtx)

	g.glCtx.DrawQueue()
}

func (g *Game) UpdatePlayersData(fighterInfo message.FighterInfo) {
	_fighter := g.fighters[fighterInfo.ID]
	_fighter.Collider = fighterInfo.Collider
	_fighter.State = fighter.FighterState(fighterInfo.State)
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.FighterInfo{
			ID:            g.players[0],
			CharacterName: string(character.Warrior),
			Collider:      g.fighters[g.players[0]].Collider,
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
		jsfunc.LogError(fmt.Sprint("JSON error:", err.Error()))
		return
	}
	g.socket.Call("send", string(jsonData))
}
