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
	fighters     []*fighter.Fighter
	running      bool
	socket       *js.Value
	glCtx        *webgl.GLContext
	keys         map[string]bool
	assets       *assetsmanager.AssetsManager
	levels       map[level.LevelName]*level.Level
	characters   map[character.CharacterName]*character.Character
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

	game := Game{
		fighters: make([]*fighter.Fighter, 2),
		socket:   socket,
		glCtx:    glCtx,
		keys:     make(map[string]bool),
		assets:   assets,
	}

	err = game.createLevels(assets)
	if err != nil {
		return nil, err
	}

	err = game.createCharacters(assets)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (g *Game) createLevels(assets *assetsmanager.AssetsManager) error {
	g.levels = make(map[level.LevelName]*level.Level)
	g.levels[level.DefaultLevel] = level.NewLevel(level.DefaultLevel, webgl.NewSprite(assets.GetTexture("background2"), nil, 1, primitives.NewVec2(0, 0), primitives.NewVec2(0, 0)))

	return nil
}

func (g *Game) createCharacters(assets *assetsmanager.AssetsManager) error {
	g.characters = make(map[character.CharacterName]*character.Character)

	warriorAnim, err := animation.NewAnimationsSet(assets.GetMetadata("warrior_anim").String(), assets.GetTexture("warrior"), 5.5, primitives.NewVec2(-98, -75), primitives.NewVec2(-181, -75))
	if err != nil {
		return err
	}
	warriorCharProperties := character.CharacterProperties{
		HealthPoints:      100,
		Attack1Damage:     10,
		Attack2Damage:     5,
		Attack1Range:      110,
		Attack2Range:      70,
		Attack1Height:     160,
		Attack2Height:     60,
		Attack1Up:         50,
		Attack2Up:         40,
		Attack1FrameIndex: 5,
		Attack2FrameIndex: 2,
	}
	warriorChar := character.NewCharacter(character.Warrior, warriorAnim["idle"].GetCurrentFrame(), warriorCharProperties)
	warriorChar.SetAnimations(warriorAnim)

	g.characters[character.Warrior] = warriorChar

	return nil
}

func (g *Game) Start(playerId string, fightersPositions map[string]int) {
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

	positions := []float64{
		config.ProgramConf.Window.Width / 4,
		config.ProgramConf.Window.Width - config.ProgramConf.Window.Width/4,
	}

	for id, fighterPos := range fightersPositions {
		if playerId == id {
			g.fighters[0] = fighter.NewFighter(g.characters[character.Warrior], positions[fighterPos])
		} else {
			g.fighters[1] = fighter.NewFighter(g.characters[character.Warrior], positions[fighterPos])
		}
	}

	if config.ProgramConf.Debug {
		if g.fighters[1] == nil {
			g.fighters[1] = fighter.NewFighter(g.characters[character.Warrior], positions[1])
		}
	}

	g.renderLoop()
}

func (g *Game) Stop() {
	g.running = false
	g.keys = make(map[string]bool)
	g.fighters = make([]*fighter.Fighter, 2)
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
	g.fighters[0].Control = message.FighterControl{
		MoveLeft:  g.keys["KeyA"],
		MoveRight: g.keys["KeyD"],
		Jump:      g.keys["Space"],
		Attack:    g.keys["KeyJ"],
	}

	g.fighters[0].Update(deltaTime, g.fighters[1])
	g.fighters[1].Update(deltaTime, g.fighters[0])
}

func (g *Game) draw() {
	g.glCtx.RenderSprite(g.currentLevel.Background, g.glCtx.Screen.BaseScreenRect, false)

	if g.fighters[1] != nil {
		g.fighters[1].Draw(g.glCtx)
	}
	g.fighters[0].Draw(g.glCtx)

	g.glCtx.DrawQueue()
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.FighterInfo{
			CharacterName: string(character.Warrior),
			HitBox:        g.fighters[0].Colliders.HitBox,
			Control: message.FighterControl{
				MoveLeft:  g.keys["KeyA"],
				MoveRight: g.keys["KeyD"],
				Jump:      g.keys["Space"],
				Attack:    g.keys["KeyJ"],
			},
		},
	}

	g.sendMessage(msg)
}

func (g *Game) UpdatePlayersData(fighterInfo message.FighterInfo) {
	g.fighters[1].Colliders.HitBox = fighterInfo.HitBox
	g.fighters[1].Control = fighterInfo.Control
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
