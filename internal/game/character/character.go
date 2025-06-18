//go:build js

package character

import (
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/webgl"
)

type AttackProperties struct {
	FrameIndex int
	Damage     float64
	Range      float64
	Height     float64
	Up         float64
}

type CharacterProperties struct {
	HealthPoints float64
	Attack1      AttackProperties
	Attack2      AttackProperties
}

type Character struct {
	Name       string
	Sprite     *webgl.Sprite
	Animations map[string]*animation.Animation
	Properies  CharacterProperties
}

func NewCharacter(name string, sprite *webgl.Sprite, properies CharacterProperties) *Character {
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
