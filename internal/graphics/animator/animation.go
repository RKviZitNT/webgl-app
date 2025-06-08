//go:build js

package animation

import "webgl-app/internal/graphics/sprite"

type Animation struct {
	Frames          []*sprite.Sprite
	FrameTime       float64
	timer           float64
	currentFrameIdx int
}

func NewAnimation(frames []*sprite.Sprite, frameTime float64) *Animation {
	return &Animation{
		Frames:          frames,
		FrameTime:       frameTime,
		timer:           0,
		currentFrameIdx: 0,
	}
}

func (a *Animation) Update(deltaTime float64) {
	a.timer += deltaTime
	if a.timer >= a.FrameTime {
		a.timer = 0 //-= deltaTime
		a.currentFrameIdx = (a.currentFrameIdx + 1) % len(a.Frames)
	}
}

func (a *Animation) GetCurrentFrame() *sprite.Sprite {
	return a.Frames[a.currentFrameIdx]
}
