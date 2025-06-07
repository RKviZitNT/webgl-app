//go:build js

package game

import (
	"encoding/json"
	"syscall/js"
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/message"
)

type Game struct {
	canvas       js.Value
	socket       *js.Value
	glCtx        *webgl.GLContext
	gl           js.Value
	program      js.Value
	vertexBuffer js.Value
	keys         map[string]bool
	characters   map[string]*character.Character
}

var (
	PlayerId string
	Width    float64
	Height   float64

	Direction primitives.Vec2
	Speed     float64
)

var vertices = []float32{
	-0.5, -0.5,
	0.5, -0.5,
	0.5, 0.5,
	-0.5, 0.5,
}

var indeces = []uint16{
	0, 1, 2,
	2, 3, 0,
}

func NewGame(socket *js.Value, glCtx *webgl.GLContext) *Game {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "game_canvas")

	Width = js.Global().Get("innerWidth").Float()
	Height = js.Global().Get("innerHeight").Float()
	canvas.Set("width", Width)
	canvas.Set("height", Height)

	return &Game{
		canvas: canvas,
		socket: socket,
		glCtx:  glCtx,
		gl:     glCtx.GL,
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

	gl := g.gl

	vShader, err := webgl.CompileShader(gl, webgl.VertexShaderSrc, gl.Get("VERTEX_SHADER"))
	if err != nil {
		return err
	}
	fShader, err := webgl.CompileShader(gl, webgl.FragmentShaderSrc, gl.Get("FRAGMENT_SHADER"))
	if err != nil {
		return err
	}

	program, err := webgl.CreateProgram(gl, vShader, fShader)
	if err != nil {
		return err
	}
	g.program = program

	g.vertexBuffer = webgl.CreateBuffer(gl, vertices, gl.Get("STATIC_DRAW"))

	popsitionAttrib := gl.Call("getAttribLocation", program, "a_position")
	gl.Call("enableVertexAttribArray", popsitionAttrib)
	gl.Call("vertexAttribPointer", popsitionAttrib, 2, gl.Get("FLOAT"), false, 0, 0)

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
	Direction = *primitives.NewVec2(0, 0)

	if g.keys["ArrowRight"] || g.keys["KeyD"] {
		Direction.AddVec2(*primitives.NewVec2(1, 0))
	}
	if g.keys["ArrowLeft"] || g.keys["KeyA"] {
		Direction.AddVec2(*primitives.NewVec2(-1, 0))
	}
	if g.keys["ArrowUp"] || g.keys["KeyW"] {
		Direction.AddVec2(*primitives.NewVec2(0, -1))
	}
	if g.keys["ArrowDown"] || g.keys["KeyS"] {
		Direction.AddVec2(*primitives.NewVec2(0, 1))
	}

	Direction.MulValue(Speed)
	g.characters[PlayerId].HitBox.Move(Direction.Normalize())
}

func (g *Game) draw() {
	gl := g.gl

	gl.Call("viewport", 0, 0, int(Width), int(Height))
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	gl.Call("useProgram", g.program)
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), g.vertexBuffer)

	for _, char := range g.characters {
		offsetLoc := gl.Call("getUniformLocation", g.program, "u_offset")
		gl.Call("uniform2f", offsetLoc, char.HitBox.Pos.X/Width*2-1, -(char.HitBox.Pos.Y/Height*2 - 1))
		gl.Call("drawArrays", gl.Get("TRIANGLE_FAN"), 0, 4)
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
