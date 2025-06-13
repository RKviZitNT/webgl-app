//go:build js

package fighter

import (
	"webgl-app/internal/config"
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
)

type Fighter struct {
	Character *character.Character
	Collider  *primitives.Rect
	State     animation.AnimationType
	Animation animation.Animation
}

func NewFighter(character *character.Character, collider *primitives.Rect) *Fighter {
	return &Fighter{
		Character: character,
		Collider:  collider,
		State:     animation.Idle,
	}
}

func (f *Fighter) SetAnimation(aType animation.AnimationType) {
	f.Animation.Reset()
	f.Animation = *f.Character.Animations[aType]
}

func (f *Fighter) Move(keys map[string]bool, deltaTime float64) {
	speed := 400.0
	var dx, dy float64

	if keys["ArrowRight"] || keys["KeyD"] {
		dx = 1
	}
	if keys["ArrowLeft"] || keys["KeyA"] {
		dx = -1
	}
	if keys["ArrowUp"] || keys["KeyW"] {
		dy = -1
	}
	if keys["ArrowDown"] || keys["KeyS"] {
		dy = 1
	}
	dir := primitives.NewVec2(dx, dy).Normalize()

	offset := dir.MulValue(speed * deltaTime)
	newRect := f.Collider.Move(offset)
	f.Collider.Pos = newRect.Pos

	if f.Collider.Left() < 0 {
		f.Collider.SetLeft(0)
	}
	if f.Collider.Right() > config.ProgramConf.Window.Width {
		f.Collider.SetRight(config.ProgramConf.Window.Width)
	}
	if f.Collider.Top() < 0 {
		f.Collider.SetTop(0)
	}
	if f.Collider.Bottom() > config.ProgramConf.Window.Height {
		f.Collider.SetBottom(config.ProgramConf.Window.Height)
	}
}
