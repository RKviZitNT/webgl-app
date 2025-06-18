//go:build js

package game

import (
	"encoding/json"
	"fmt"
	"math"
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
	"webgl-app/internal/utils"
)

type GameState struct {
	isStart       bool
	isEnd         bool
	startCooldown float64
	endCooldown   float64
}

type Game struct {
	gameState    GameState
	fighters     []*fighter.Fighter
	running      bool
	socket       *js.Value
	glCtx        *webgl.GLContext
	keys         map[string]bool
	assets       *assetsmanager.AssetsManager
	levels       map[string]*level.Level
	titles       map[string]*webgl.Sprite
	characters   map[string]*character.Character
	healthBar    *animation.Animation
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
		socket: socket,
		glCtx:  glCtx,
		keys:   make(map[string]bool),
		assets: assets,
	}

	err = game.createLevels(assets)
	if err != nil {
		return nil, err
	}

	err = game.createTitles(assets)
	if err != nil {
		return nil, err
	}

	err = game.createCharacters(assets)
	if err != nil {
		return nil, err
	}

	err = game.createHealthBar(assets)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (g *Game) createLevels(assets *assetsmanager.AssetsManager) error {
	g.levels = make(map[string]*level.Level)

	var err error
	g.levels["level_1"], err = level.NewLevel("level_1", assets.GetTexture("background1"))
	if err != nil {
		return err
	}
	g.levels["level_2"], err = level.NewLevel("level_2", assets.GetTexture("background2"))
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) createTitles(assets *assetsmanager.AssetsManager) error {
	g.titles = make(map[string]*webgl.Sprite)

	texFlags := assets.GetTexture("flags")
	g.titles["victory"] = webgl.NewSprite(texFlags, &primitives.Rect{Pos: primitives.NewVec2(0, 0), Size: primitives.NewVec2(texFlags.Width, texFlags.Height/2)}, 2, nil, nil)
	g.titles["defeat"] = webgl.NewSprite(texFlags, &primitives.Rect{Pos: primitives.NewVec2(0, texFlags.Height/2), Size: primitives.NewVec2(texFlags.Width, texFlags.Height/2)}, 2, nil, nil)

	texStart := assets.GetTexture("start")
	g.titles["start"] = webgl.NewSprite(texStart, &primitives.Rect{Pos: primitives.NewVec2(0, 0), Size: primitives.NewVec2(texStart.Width, texStart.Height)}, 1.5, nil, nil)
	return nil
}

func (g *Game) createCharacters(assets *assetsmanager.AssetsManager) error {
	g.characters = make(map[string]*character.Character)

	warriorAnim, err := animation.NewAnimationsSet("warrior_meta", assets, 5.5, primitives.NewVec2(-98, -75), primitives.NewVec2(-181, -75))
	if err != nil {
		return err
	}
	warriorCharProperties := character.CharacterProperties{
		HealthPoints: 100,
		Attack1: character.AttackProperties{
			FrameIndex: 5,
			Damage:     9,
			Range:      110,
			Height:     160,
			Up:         40,
		},
		Attack2: character.AttackProperties{
			FrameIndex: 1,
			Damage:     6,
			Range:      80,
			Height:     120,
			Up:         30,
		},
	}
	warriorChar := character.NewCharacter("warrior", warriorAnim["idle"].GetCurrentFrame(), warriorCharProperties)
	warriorChar.SetAnimations(warriorAnim)

	g.characters["warrior"] = warriorChar

	return nil
}

func (g *Game) createHealthBar(assets *assetsmanager.AssetsManager) error {
	tex := assets.GetTexture("health-bar")
	animationInfo := animation.CreateAnimationInfo{
		Texture: tex,
		AnimationData: animation.AnimationData{
			FirstFrame: 1,
			FrameCount: 7,
		},
		SpritesheetData: &animation.AnimationsSpritesheetData{
			FrameWidthCount:  7,
			FrameHeigntCount: 1,
			AllFrameCount:    7,
		},
		Scale:          5,
		Offset:         primitives.Vec2{},
		SpecularOffset: primitives.Vec2{},
	}

	g.healthBar = animation.NewAnimation(animationInfo)

	return nil
}

func (g *Game) Start(playerId string, fightersPositions map[string]int) {
	g.gameState = GameState{
		isStart:       false,
		isEnd:         false,
		startCooldown: 2,
		endCooldown:   3,
	}

	g.fighters = make([]*fighter.Fighter, 2)

	g.currentLevel = g.levels["level_1"]

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
		config.ProgramConfig.Window.Width / 4,
		config.ProgramConfig.Window.Width - config.ProgramConfig.Window.Width/4,
	}

	for id, fighterPos := range fightersPositions {
		if playerId == id {
			g.fighters[0] = fighter.NewFighter(g.characters["warrior"], positions[fighterPos])
		} else {
			g.fighters[1] = fighter.NewFighter(g.characters["warrior"], positions[fighterPos])
		}
	}

	if config.ProgramConfig.Debug {
		if g.fighters[1] == nil {
			g.fighters[1] = fighter.NewFighter(g.characters["warrior"], positions[1])
		}
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
		frameTime     = time.Second / time.Duration(config.ProgramConfig.Window.FrameRate)
		lastFrameTime time.Time
		elapsedTime   time.Duration
	)

	lastFrameTime = time.Now()

	g.running = true
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !g.running {
			return nil
		}

		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastFrameTime)
		lastFrameTime = currentTime

		g.update(deltaTime)
		g.draw()

		g.sendPlayerState()

		elapsedTime = time.Now().Sub(currentTime)
		if elapsedTime < frameTime {
			time.Sleep(frameTime - elapsedTime)
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (g *Game) update(deltaTime time.Duration) {
	if !g.running || g.fighters[0] == nil || g.fighters[1] == nil {
		return
	}

	if !g.gameState.isStart {
		g.gameState.startCooldown -= deltaTime.Seconds()
		if g.gameState.startCooldown <= 0 {
			g.gameState.isStart = true
		}
	}
	if g.gameState.isEnd {
		g.gameState.endCooldown -= deltaTime.Seconds()
		if g.gameState.endCooldown <= 0 {
			g.sendEndGameMsg()
		}
	}

	if g.keys["Escape"] {
		g.sendEndGameMsg()
	}

	if g.gameState.isStart {
		g.fighters[0].Control = message.FighterControl{
			MoveLeft:  g.keys["KeyA"],
			MoveRight: g.keys["KeyD"],
			Jump:      g.keys["Space"],
			Attack:    g.keys["KeyJ"],
		}
	}

	g.fighters[0].Update(deltaTime, g.fighters[1])
	g.fighters[1].Update(deltaTime, g.fighters[0])

	if g.fighters[0].State == fighter.Death || g.fighters[1].State == fighter.Death {
		g.gameState.isEnd = true
	}
}

func (g *Game) draw() {
	if !g.running {
		return
	}

	g.currentLevel.Draw(g.glCtx)

	g.healthBarsDraw()

	if g.fighters[1] != nil {
		g.fighters[1].Draw(g.glCtx)
	}
	g.fighters[0].Draw(g.glCtx)

	g.titleDraw()

	g.glCtx.DrawQueue()
}

func (g *Game) titleDraw() {
	if !g.gameState.isStart {
		g.glCtx.RenderSprite(g.titles["start"], primitives.NewRect(530, -50, 0, 0), false)
	}
	if g.gameState.isEnd {
		if g.fighters[0].State == fighter.Death {
			g.glCtx.RenderSprite(g.titles["defeat"], primitives.NewRect(440, 75, 0, 0), false)
		} else {
			g.glCtx.RenderSprite(g.titles["victory"], primitives.NewRect(440, 0, 0, 0), false)
		}
	}
}

func (g *Game) healthBarsDraw() {
	g.glCtx.RenderSprite(g.healthBar.GetFrame(0), primitives.NewRect(50, 30, 0, 0), false)
	g.glCtx.RenderSprite(g.healthBar.GetFrame(0), primitives.NewRect(1310, 30, 0, 0), true)

	g.glCtx.RenderSprite(g.healthBar.GetFrame(2), primitives.NewRect(50, 30, 0, 0), false)
	g.glCtx.RenderSprite(g.healthBar.GetFrame(2), primitives.NewRect(1310, 30, 0, 0), true)

	if g.fighters[0] != nil {
		hpPercent := g.fighters[0].Properties.HealthPoints / 100.0
		frameIndex := int(math.Ceil(hpPercent * 6))
		frameIndex = utils.Clamp(frameIndex, 0, 6)

		g.glCtx.RenderSprite(g.healthBar.GetFrame(7-frameIndex), primitives.NewRect(50, 30, 0, 0), false)
	}

	if g.fighters[1] != nil {
		hpPercent := g.fighters[1].Properties.HealthPoints / 100.0
		frameIndex := int(math.Ceil(hpPercent * 6))
		frameIndex = utils.Clamp(frameIndex, 0, 6)

		g.glCtx.RenderSprite(g.healthBar.GetFrame(7-frameIndex), primitives.NewRect(1310, 30, 0, 0), true)
	}
}

func (g *Game) sendPlayerState() {
	msg := message.Message{
		Type: message.GameStateMsg,
		Data: message.FighterInfo{
			CharacterName: g.fighters[0].Character.Name,
			HealthPoints:  g.fighters[0].Properties.HealthPoints,
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
	if g.fighters[1].Character.Name != fighterInfo.CharacterName {
		g.fighters[1].Character = g.characters[fighterInfo.CharacterName]
	}
	g.fighters[1].Properties.HealthPoints = fighterInfo.HealthPoints
	g.fighters[1].Colliders.HitBox = fighterInfo.HitBox
	if g.gameState.isStart {
		g.fighters[1].Control = fighterInfo.Control
	}
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
