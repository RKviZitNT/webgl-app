//go:build js

package webgl

import (
	"fmt"
	"syscall/js"
	"webgl-app/internal/graphics/primitives"
)

type Screen struct {
	BaseScreenRect *primitives.Rect
	ScreenRect     *primitives.Rect
	canvas         js.Value
	aspectRatio    float64
}

func newScreen(canvasId string, baseScreenRect *primitives.Rect) (*Screen, error) {
	canvas := js.Global().Get("document").Call("getElementById", canvasId)
	if canvas.IsUndefined() || canvas.IsNull() {
		return nil, fmt.Errorf("NewScreen: canvas is nil")
	}

	screen := &Screen{
		BaseScreenRect: baseScreenRect,
		ScreenRect:     baseScreenRect,
		canvas:         canvas,
		aspectRatio:    float64(baseScreenRect.Width()) / float64(baseScreenRect.Height()),
	}
	screen.Update()

	return screen, nil
}

func (s *Screen) Update() {
	windowWidth := js.Global().Get("innerWidth").Float()
	windowHeight := js.Global().Get("innerHeight").Float()

	targetWidth := windowWidth
	targetHeight := windowWidth / s.aspectRatio

	if targetHeight > s.ScreenRect.Height() {
		targetWidth = windowHeight * s.aspectRatio
		targetHeight = windowHeight
	}

	s.canvas.Set("width", targetWidth)
	s.canvas.Set("height", targetHeight)
	s.ScreenRect = primitives.NewRect(primitives.NewVec2(0, 0), primitives.NewVec2(targetWidth, targetHeight))
}
