//go:build js

package fighter

import (
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
)

type Fighter struct {
	Character *character.Character
	Collider  primitives.Rect
	State     animation.AnimationType
}

func NewFighter(character *character.Character, collider primitives.Rect) *Fighter {
	return &Fighter{
		Character: character,
		Collider:  collider,
		State:     animation.Idle,
	}
}
