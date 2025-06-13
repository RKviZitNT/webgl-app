//go:build js

package character

import (
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/webgl"
)

type CharacterName string

const (
	Warrior CharacterName = "warrior"
)

type Character struct {
	Name       CharacterName
	Sprite     *webgl.Sprite
	Animations map[animation.AnimationType]*animation.Animation
}

func NewCharacter(name CharacterName, sprite *webgl.Sprite) *Character {
	return &Character{
		Name:       name,
		Sprite:     sprite,
		Animations: make(map[animation.AnimationType]*animation.Animation),
	}
}

func (c *Character) AddAnimation(aType animation.AnimationType, animation *animation.Animation) {
	c.Animations[aType] = animation
}

func (c *Character) SetAnimations(animations map[animation.AnimationType]*animation.Animation) {
	c.Animations = animations
}
