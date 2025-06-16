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
	HealthPoints      float64
	Attack1Damage     float64
	Attack2Damage     float64
	Attack1Range      float64
	Attack2Range      float64
	Attack1Height     float64
	Attack2Height     float64
	Attack1Up         float64
	Attack2Up         float64
	Attack1FrameIndex int
	Attack2FrameIndex int
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
