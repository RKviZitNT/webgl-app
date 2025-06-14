//go:build js

package animation

import (
	"fmt"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/utils"
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
	Frames          []*webgl.Sprite
	FrameTime       float64
	timer           float64
	currentFrameIdx int
	isEnd           bool
}

func NewAnimation(aName string, data AnimationsData, texture *webgl.Texture, scale float64, offset primitives.Vec2, specularOffset primitives.Vec2) *Animation {
	if texture == nil {
		return nil
	}

	aData := data.Animations[aName]

	frameWidth := data.Parameters.Width / data.Parameters.FrameWidthCount
	frameHeight := data.Parameters.Height / data.Parameters.FrameHeigntCount
	frames := make([]*webgl.Sprite, 0, aData.FrameCount)

	for i := aData.FirstFrame; i < aData.FirstFrame+aData.FrameCount; i++ {
		if i > data.Parameters.AllFrameCount {
			break
		}

		col := (i - 1) % data.Parameters.FrameWidthCount
		row := (i - 1) / data.Parameters.FrameWidthCount

		rect := primitives.NewRect(float64(col*frameWidth), float64(row*frameHeight), float64(frameWidth), float64(frameHeight))

		frames = append(frames, webgl.NewSprite(texture, &rect, scale, offset, specularOffset))
	}

	return &Animation{
		Frames:          frames,
		FrameTime:       aData.FrameTime,
		timer:           0,
		currentFrameIdx: 0,
	}
}

func NewAnimationsSet(metadata string, texture *webgl.Texture, scale float64, offset primitives.Vec2, specularOffset primitives.Vec2) (map[string]*Animation, error) {
	var data AnimationsData
	err := utils.ParseStringToJSON(metadata, &data)
	if err != nil {
		return nil, err
	}

	animations := make(map[string]*Animation)
	for aName, _ := range data.Animations {
		animations[aName] = NewAnimation(aName, data, texture, scale, offset, specularOffset)
	}
	return animations, nil
}

func (a *Animation) Reset() {
	if a == nil {
		return
	}
	a.timer = 0
	a.currentFrameIdx = 0
}

func (a *Animation) Update(deltaTime float64) {
	if a == nil {
		jsfunc.LogError("Animation.Update: nil animation")
		return
	}

	if len(a.Frames) == 0 {
		jsfunc.LogError("Animation.Update: empty frames")
		return
	}

	a.timer += deltaTime
	if a.timer > a.FrameTime {
		a.timer = 0
		a.currentFrameIdx = (a.currentFrameIdx + 1) % len(a.Frames)
		if a.currentFrameIdx == 0 {
			a.isEnd = true
		}
	}
}

func (a *Animation) GetCurrentFrame() *webgl.Sprite {
	if a == nil {
		jsfunc.LogError("Animation.Update: nil animation")
		return nil
	}

	if len(a.Frames) == 0 {
		jsfunc.LogError("Animation.Update: empty frames")
		return nil
	}

	if a.currentFrameIdx >= len(a.Frames) {
		jsfunc.LogError(fmt.Sprintf("Animation.GetCurrentFrame: invalid frame index %d/%d", a.currentFrameIdx, len(a.Frames)))
		return nil
	}

	frame := a.Frames[a.currentFrameIdx]
	if frame == nil {
		jsfunc.LogError(fmt.Sprintf("Animation.GetCurrentFrame: nil frame at index %d", a.currentFrameIdx))
	}

	return frame
}

func (a *Animation) IsEnd() bool {
	return a.isEnd
}
