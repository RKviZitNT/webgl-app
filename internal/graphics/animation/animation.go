//go:build js

package animation

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/texture"
)

type AnimationType string

const (
	Idle AnimationType = "idle"
)

type AnimationsParameters struct {
	Width            int `json:"width"`
	Height           int `json:"height"`
	FrameWidthCount  int `json:"frame_width_count"`
	FrameHeigntCount int `json:"frame_height_count"`
	AllFrameCount    int `json:"all_frame_count"`
}

type AnimationData struct {
	FirstFrame int `json:"first_frame"`
	FrameCount int `json:"frame_count"`
}

type AnimationsMetaData struct {
	Parameters AnimationsParameters            `json:"parameters"`
	Animations map[AnimationType]AnimationData `json:"animations"`
}

type Animation struct {
	Frames          []*sprite.Sprite
	FrameTime       float64
	timer           float64
	currentFrameIdx int
}

func NewAnimation(aType AnimationType, metaData *AnimationsMetaData, texture *texture.Texture, frameTime float64) *Animation {
	aData := metaData.Animations[aType]

	frameWidth := metaData.Parameters.Width / metaData.Parameters.FrameWidthCount
	frameHeight := metaData.Parameters.Height / metaData.Parameters.FrameHeigntCount
	frames := make([]*sprite.Sprite, 0)

	for i := aData.FirstFrame - 1; i < aData.FirstFrame+aData.FrameCount; i++ {
		frames = append(frames, sprite.NewSprite(texture, primitives.NewRect(primitives.NewVec2(float64(frameWidth*i%metaData.Parameters.Width), float64(frameHeight*i%metaData.Parameters.Height)), primitives.NewVec2(float64(frameWidth), float64(frameHeight)))))
	}

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
