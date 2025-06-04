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
	gl         js.Value
	socket     *js.Value
	canvas     js.Value
	keys       map[string]bool
	characters map[string]*character.Character
	program    js.Value
	position   js.Value
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

	gl := canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		panic("WebGL not supported")
	}

	vertexShaderSource := `
        attribute vec2 a_position;
        uniform vec2 u_resolution;

        void main() {
            vec2 zeroToOne = a_position / u_resolution;
            vec2 zeroToTwo = zeroToOne * 2.0;
            vec2 clipSpace = zeroToTwo - 1.0;

            gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
        }
    `

	fragmentShaderSource := `
        precision mediump float;
        void main() {
            gl_FragColor = vec4(0, 0, 1, 1); // blue
        }
    `

	program := createProgram(gl, vertexShaderSource, fragmentShaderSource)
	gl.Call("useProgram", program)

	positionLocation := gl.Call("getAttribLocation", program, "a_position")
	return &Game{
		gl:         gl,
		socket:     sk,
		canvas:     canvas,
		keys:       make(map[string]bool),
		characters: chars,
		program:    program,
		position:   positionLocation,
	}
}

func (g *Game) Start(playerId string) {
	PlayerId = playerId

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

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Update()
		g.sendPlayerState()
		g.Draw()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (g *Game) Update() {
	if g.keys["ArrowRight"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 2, Y: 0})
	}
	if g.keys["ArrowLeft"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: -2, Y: 0})
	}
	if g.keys["ArrowUp"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 0, Y: -2})
	}
	if g.keys["ArrowDown"] {
		g.characters[PlayerId].HitBox.Move(primitives.Vec2{X: 0, Y: 2})
	}
}

func (g *Game) Draw() {
	gl := g.gl
	gl.Call("viewport", 0, 0, int(Width), int(Height))
	gl.Call("clearColor", 0.9, 0.9, 0.9, 1.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))

	resLoc := gl.Call("getUniformLocation", g.program, "u_resolution")
	gl.Call("uniform2f", resLoc, Width, Height)

	for _, c := range g.characters {
		x := float32(c.HitBox.Pos.X)
		y := float32(c.HitBox.Pos.Y)
		size := float32(50.0)

		vertices := []float32{
			x, y,
			x + size, y,
			x, y + size,
			x, y + size,
			x + size, y,
			x + size, y + size,
		}

		buffer := gl.Call("createBuffer")
		gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buffer)

		jsVertices := js.Global().Get("Float32Array").New(len(vertices))
		for i, v := range vertices {
			jsVertices.SetIndex(i, v)
		}
		gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), jsVertices, gl.Get("STATIC_DRAW"))

		gl.Call("enableVertexAttribArray", g.position)
		gl.Call("vertexAttribPointer", g.position, 2, gl.Get("FLOAT"), false, 0, 0)

		gl.Call("drawArrays", gl.Get("TRIANGLES"), 0, 6)
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

func createProgram(gl js.Value, vertexSource, fragmentSource string) js.Value {
	vertexShader := compileShader(gl, gl.Get("VERTEX_SHADER"), vertexSource)
	fragmentShader := compileShader(gl, gl.Get("FRAGMENT_SHADER"), fragmentSource)

	program := gl.Call("createProgram")
	gl.Call("attachShader", program, vertexShader)
	gl.Call("attachShader", program, fragmentShader)
	gl.Call("linkProgram", program)

	if !gl.Call("getProgramParameter", program, gl.Get("LINK_STATUS")).Bool() {
		println("Could not link shader program")
	}

	return program
}

func compileShader(gl js.Value, shaderType js.Value, source string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)

	if !gl.Call("getShaderParameter", shader, gl.Get("COMPILE_STATUS")).Bool() {
		println("Shader compile failed:", gl.Call("getShaderInfoLog", shader).String())
	}
	return shader
}
