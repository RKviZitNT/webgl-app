//go:build js

package fighter

import (
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/primitives"
)

type Fighter struct {
	Character *character.Character
	Collider  primitives.Rect
}

func NewFighter(characterName string, collider primitives.Rect) *Fighter {
	return &Fighter{
		Character: character.Characters[character.CharacterName(characterName)],
		Collider:  collider,
	}
}
