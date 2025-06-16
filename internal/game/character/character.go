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

type CharacterProperties struct {
	HealthPoints float64
	AttackDamage float64
	AttackRange  float64
	AttackHeight float64
	AttackUp     float64
}

type Character struct {
	Name       CharacterName
	Sprite     *webgl.Sprite
	Animations map[string]*animation.Animation
	Properies  CharacterProperties
}

func NewCharacter(name CharacterName, sprite *webgl.Sprite, properies CharacterProperties) *Character {
	return &Character{
		Name:       name,
		Sprite:     sprite,
		Animations: make(map[string]*animation.Animation),
		Properies:  properies,
	}
}

func (c *Character) AddAnimation(aType string, anim *animation.Animation) {
	c.Animations[aType] = anim
}

func (c *Character) SetAnimations(anims map[string]*animation.Animation) {
	c.Animations = anims
}
