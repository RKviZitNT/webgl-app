//go:build js

package animation

import (
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/sprite"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/utils"
)

type AnimationType string

const (
	Idle    AnimationType = "idle"
	Walk    AnimationType = "walk"
	Run     AnimationType = "run"
	Attack1 AnimationType = "attack1"
	Attack2 AnimationType = "attack2"
	Death   AnimationType = "death"
	Hurt    AnimationType = "hurt"
	Jump    AnimationType = "jump"
)

type AnimationsParameters struct {
	Width            int `json:"width"`
	Height           int `json:"height"`
	FrameWidthCount  int `json:"frame_width_count"`
	FrameHeigntCount int `json:"frame_height_count"`
	AllFrameCount    int `json:"all_frame_count"`
}

type AnimationData struct {
	FrameTime  float64 `json:"frame_time"`
	FirstFrame int     `json:"first_frame"`
	FrameCount int     `json:"frame_count"`
}

type AnimationsData struct {
	Parameters AnimationsParameters     `json:"parameters"`
	Animations map[string]AnimationData `json:"animations"`
}

type Animation struct {
	Frames          []*sprite.Sprite
	FrameTime       float64
	timer           float64
	currentFrameIdx int
}

func NewAnimation(aType AnimationType, data *AnimationsData, texture *texture.Texture) *Animation {
	aData := data.Animations[string(aType)]

	frameWidth := data.Parameters.Width / data.Parameters.FrameWidthCount
	frameHeight := data.Parameters.Height / data.Parameters.FrameHeigntCount
	frames := make([]*sprite.Sprite, 0, aData.FrameCount)

	for i := aData.FirstFrame; i < aData.FirstFrame+aData.FrameCount; i++ {
		if i > data.Parameters.AllFrameCount {
			break
		}

		col := (i - 1) % data.Parameters.FrameWidthCount
		row := (i - 1) / data.Parameters.FrameWidthCount

		pos := primitives.NewVec2(
			float64(col*frameWidth),
			float64(row*frameHeight),
		)
		size := primitives.NewVec2(
			float64(frameWidth),
			float64(frameHeight),
		)
		rect := primitives.NewRect(pos, size)

		frames = append(frames, sprite.NewSprite(texture, rect))
	}

	return &Animation{
		Frames:          frames,
		FrameTime:       aData.FrameTime,
		timer:           0,
		currentFrameIdx: 0,
	}
}

func CreateAnimations(aTypes []AnimationType, metadata string, texture *texture.Texture) (map[AnimationType]*Animation, error) {
	var data *AnimationsData
	err := utils.ParseStringToJSON(metadata, &data)
	if err != nil {
		return nil, err
	}

	animations := make(map[AnimationType]*Animation)
	for _, aType := range aTypes {
		animations[aType] = NewAnimation(aType, data, texture)
	}
	return animations, nil
}

func (a *Animation) Update(deltaTime float64) {
	a.timer += deltaTime
	if a.timer > a.FrameTime {
		a.timer = 0
		a.currentFrameIdx = (a.currentFrameIdx + 1) % len(a.Frames)
	}
}

func (a *Animation) GetCurrentFrame() *sprite.Sprite {
	return a.Frames[a.currentFrameIdx]
}
