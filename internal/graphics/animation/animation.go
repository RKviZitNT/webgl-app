//go:build js

package animation

import (
	"fmt"
	"webgl-app/internal/assetsmanager"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/utils"
)

const (
	TypeAnimationsSpritesheet string = "animation-spritesheet"
	TypeAnimationsSet         string = "animation-set"
)

type AnimationsSpritesheetCutArea struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type AnimationTextureData struct {
	Texture          string `json:"texture"`
	FrameWidthCount  int    `json:"frame_width_count"`
	FrameHeigntCount int    `json:"frame_height_count"`
	AllFrameCount    int    `json:"all_frame_count"`
}

type AnimationsSpritesheetData struct {
	SpritesheetTexture string                       `json:"spritesheet_texture"`
	FrameWidthCount    int                          `json:"frame_width_count"`
	FrameHeigntCount   int                          `json:"frame_height_count"`
	AllFrameCount      int                          `json:"all_frame_count"`
	CutArea            AnimationsSpritesheetCutArea `json:"cut_area"`
}

type AnimationsParameters struct {
	Type            string                    `json:"type"`
	SpritesheetData AnimationsSpritesheetData `json:"spritesheet_data"`
}

type AnimationData struct {
	TextureData AnimationTextureData `json:"texture_data"`
	FrameTime   float64              `json:"frame_time"`
	FirstFrame  int                  `json:"first_frame"`
	FrameCount  int                  `json:"frame_count"`
}

type AnimationsMeta struct {
	Parameters AnimationsParameters     `json:"parameters"`
	Animations map[string]AnimationData `json:"animations"`
}

type CreateAnimationInfo struct {
	Texture         *webgl.Texture
	AnimationData   AnimationData
	SpritesheetData *AnimationsSpritesheetData
	Scale           float64
	Offset          primitives.Vec2
	SpecularOffset  primitives.Vec2
}

type Animation struct {
	Frames            []*webgl.Sprite
	FrameTime         float64
	CurrentFrameIndex int
	IsEnd             bool
	timer             float64
}

func NewAnimation(info CreateAnimationInfo) *Animation {
	var frames []*webgl.Sprite

	aData := info.AnimationData

	if info.SpritesheetData != nil {
		frameWidth := info.Texture.Width / float64(info.SpritesheetData.FrameWidthCount)
		frameHeight := info.Texture.Height / float64(info.SpritesheetData.FrameHeigntCount)

		shData := info.SpritesheetData

		frames = make([]*webgl.Sprite, 0)

		for i := aData.FirstFrame; i < aData.FirstFrame+aData.FrameCount; i++ {
			if i > shData.AllFrameCount {
				break
			}

			col := float64((i - 1) % shData.FrameWidthCount)
			row := float64((i - 1) / shData.FrameWidthCount)

			rect := primitives.NewRect(float64(col*frameWidth), float64(row*frameHeight), frameWidth, frameHeight)

			frames = append(frames, webgl.NewSprite(info.Texture, &rect, info.Scale, &info.Offset, &info.SpecularOffset))
		}
	} else {
		// texSize := info.Texture.Size()
		// frameWidth := texSize.X / float64(info.AnimationData.FrameCount)
		// frameHeight := texSize.Y / float64(info.AnimationData.FrameCount)
	}

	return &Animation{
		Frames:            frames,
		FrameTime:         aData.FrameTime,
		CurrentFrameIndex: 0,
		timer:             0,
	}
}

func NewAnimationsSet(metaName string, assets *assetsmanager.AssetsManager, scale float64, offset primitives.Vec2, specularOffset primitives.Vec2) (map[string]*Animation, error) {
	var meta AnimationsMeta
	metadataSrc := assets.GetMetadata(metaName).String()
	if metadataSrc == "" {
		return nil, fmt.Errorf("Metadata not found for %s", metaName)
	}
	err := utils.ParseStringToJSON(metadataSrc, &meta)
	if err != nil {
		return nil, err
	}

	animations := make(map[string]*Animation)

	switch meta.Parameters.Type {
	case TypeAnimationsSpritesheet:
		spritesheetTex := assets.GetTexture(meta.Parameters.SpritesheetData.SpritesheetTexture)

		animationInfo := CreateAnimationInfo{
			Texture:         spritesheetTex,
			SpritesheetData: &meta.Parameters.SpritesheetData,
			Scale:           scale,
			Offset:          offset,
			SpecularOffset:  specularOffset,
		}

		for aName, _ := range meta.Animations {
			animationInfo.AnimationData = meta.Animations[aName]
			animations[aName] = NewAnimation(animationInfo)
		}

	case TypeAnimationsSet:
		return nil, fmt.Errorf("Incorrect type")

	default:
		return nil, fmt.Errorf("Incorrect type")
	}

	return animations, nil
}

func (a *Animation) Reset() {
	if a == nil {
		return
	}
	a.timer = 0
	a.CurrentFrameIndex = 0
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
		a.CurrentFrameIndex = (a.CurrentFrameIndex + 1) % len(a.Frames)
		if a.CurrentFrameIndex == len(a.Frames)-1 {
			a.IsEnd = true
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

	if a.CurrentFrameIndex >= len(a.Frames) {
		jsfunc.LogError(fmt.Sprintf("Animation.GetCurrentFrame: invalid frame index %d/%d", a.CurrentFrameIndex, len(a.Frames)))
		return nil
	}

	frame := a.Frames[a.CurrentFrameIndex]
	if frame == nil {
		jsfunc.LogError(fmt.Sprintf("Animation.GetCurrentFrame: nil frame at index %d", a.CurrentFrameIndex))
	}

	return frame
}

func (a *Animation) GetFrame(index int) *webgl.Sprite {
	frame := a.Frames[index%len(a.Frames)]
	if frame == nil {
		jsfunc.LogError(fmt.Sprintf("Animation.GetFrame: nil frame at index %d", a.CurrentFrameIndex))
	}

	return frame
}
