//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
)

type Screen struct {
	Canvas         js.Value
	BaseScreenRect *primitives.Rect
	ScreenRect     *primitives.Rect
	AspectRatio    float64
	Scale          *primitives.Vec2
	Offset         *primitives.Vec2
}

func NewScreen(canvasId string, baseScreenRect *primitives.Rect) (*Screen, error) {
	canvas := js.Global().Get("document").Call("getElementById", canvasId)
	if canvas.IsUndefined() || canvas.IsNull() {
		return nil, fmt.Errorf("NewScreen: canvas is nil")
	}

	screen := &Screen{
		Canvas:         canvas,
		BaseScreenRect: baseScreenRect,
		ScreenRect:     baseScreenRect,
		AspectRatio:    float64(baseScreenRect.Width()) / float64(baseScreenRect.Height()),
	}
	screen.Update()

	return screen, nil
}

func (s *Screen) Update() {
	windowWidth := js.Global().Get("innerWidth").Float()
	windowHeight := js.Global().Get("innerHeight").Float()

	targetWidth := windowWidth
	targetHeight := windowWidth / s.AspectRatio

	if targetHeight > s.ScreenRect.Height() {
		targetWidth = windowHeight * s.AspectRatio
		targetHeight = windowHeight
	}

	s.Canvas.Set("width", targetWidth)
	s.Canvas.Set("height", targetHeight)
	s.ScreenRect = primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(targetWidth, targetHeight))
}
