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
	Animations map[animation.AnimationType]*animation.Animation
}

func NewCharacter(name CharacterName, sprite *sprite.Sprite) *Character {
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
