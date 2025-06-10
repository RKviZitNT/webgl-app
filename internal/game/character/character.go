//go:build js

package character

import (
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/sprite"
)

type CharacterName string

const (
	Warrior CharacterName = "warrior"
)

type Character struct {
	Name       CharacterName
	Sprite     *sprite.Sprite
	Animations map[string]*animation.Animation
}

func NewCharacter(name CharacterName, sprite *sprite.Sprite) *Character {
	return &Character{
		Name:       name,
		Sprite:     sprite,
		Animations: make(map[string]*animation.Animation),
	}
}

func (c *Character) AddAnimation(name string, animation *animation.Animation) {
	c.Animations[name] = animation
}
