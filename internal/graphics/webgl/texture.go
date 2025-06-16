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

func NewTexture(glCtx *GLContext, img js.Value) *Texture {
	if img.IsUndefined() || img.IsNull() {
		jsfunc.LogError("NewTexture: image is undefined or null")
		return nil
	}

	gl := glCtx.GL

	tex := gl.Call("createTexture")
	gl.Call("bindTexture", gl.Get("TEXTURE_2D"), tex)

	gl.Call("texImage2D", gl.Get("TEXTURE_2D"), 0, gl.Get("RGBA"), gl.Get("RGBA"), gl.Get("UNSIGNED_BYTE"), img)

	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_WRAP_S"), gl.Get("CLAMP_TO_EDGE"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_WRAP_T"), gl.Get("CLAMP_TO_EDGE"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_MIN_FILTER"), gl.Get("NEAREST"))
	gl.Call("texParameteri", gl.Get("TEXTURE_2D"), gl.Get("TEXTURE_MAG_FILTER"), gl.Get("NEAREST"))

	return &Texture{
		texture: tex,
		width:   img.Get("width").Int(),
		height:  img.Get("height").Int(),
	}
}

func (t *Texture) GetTexture() *js.Value {
	return &t.texture
}

func (t *Texture) Size() primitives.Vec2 {
	return primitives.NewVec2(float64(t.width), float64(t.height))
}
