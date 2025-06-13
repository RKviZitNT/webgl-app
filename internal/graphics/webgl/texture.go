//go:build js

package webgl

import (
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/jsfunc"
)

type Texture struct {
	texture js.Value
	width   int
	height  int
}

func NewTexture(gl js.Value, image js.Value) *Texture {
	if image.IsUndefined() || image.IsNull() {
		jsfunc.LogError("NewTexture: image is undefined or null")
		return nil
	}

	texture := gl.Call("createTexture")
	gl.Call("bindTexture", gl.Get("TEXTURE_2D"), texture)

	gl.Call("texImage2D", gl.Get("TEXTURE_2D"), 0, gl.Get("RGBA"), gl.Get("RGBA"), gl.Get("UNSIGNED_BYTE"), image)

	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_WRAP_S"), gl.Get("CLAMP_TO_EDGE"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_WRAP_T"), gl.Get("CLAMP_TO_EDGE"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_MIN_FILTER"), gl.Get("NEAREST"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_MAG_FILTER"), gl.Get("NEAREST"))

	return &Texture{
		texture: texture,
		width:   image.Get("width").Int(),
		height:  image.Get("height").Int(),
	}
}

func (t *Texture) GetTexture() *js.Value {
	return &t.texture
}

func (t *Texture) Size() *primitives.Vec2 {
	return primitives.NewVec2(float64(t.width), float64(t.height))
}
