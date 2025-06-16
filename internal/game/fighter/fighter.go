//go:build js

package fighter

import (
	"time"
	"webgl-app/internal/config"
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/net/message"
)

type FighterState string

const (
	Idle    FighterState = "idle"
	Run     FighterState = "run"
	Jump    FighterState = "jump"
	Attack1 FighterState = "attack1"
	Attack2 FighterState = "attack2"
	Hit     FighterState = "hit"
	Death   FighterState = "death"
)

type Colliders struct {
	HitBox primitives.Rect
	Attack primitives.Rect
}

type HoldingKeys struct {
	MoveLeft  bool
	MoveRight bool
	Jump      bool
	Attack    bool
}

type FighterProperties struct {
	Specular             bool
	HealthPoints         float64
	dir                  primitives.Vec2
	speed                float64
	jumpSpeed            float64
	velocity             float64
	gravity              float64
	invulnerableCooldown float64
	attackCooldown       float64
	move                 bool
	jump                 bool
	attack               bool
	comboAttack          bool
	hit                  bool
	death                bool
}

type Fighter struct {
	State          FighterState
	AnimationState FighterState
	Colliders      Colliders
	Character      *character.Character
	Animation      animation.Animation
	Properties     FighterProperties
	Control        message.FighterControl
	HoldingKeys    HoldingKeys
}

func NewFighter(char *character.Character, posX float64) *Fighter {
	colliders := Colliders{
		HitBox: primitives.NewRect(0, 0, 100, 160),
	}
	colliders.HitBox.SetCenter(primitives.NewVec2(posX, config.ProgramConf.Window.Height))

	Properties := FighterProperties{
		HealthPoints: char.Properies.HealthPoints,
		dir:          primitives.Vec2{},
		jumpSpeed:    1500,
		speed:        500,
		velocity:     0,
		gravity:      5000,
	}

	fighter := Fighter{
		State:      Idle,
		Colliders:  colliders,
		Character:  char,
		Animation:  *char.Animations[string(Idle)],
		Properties: Properties,
	}

	return &fighter
}

func (f *Fighter) updateAnimationState() {
	if f.AnimationState != f.State {
		f.AnimationState = f.State
		f.Animation = *f.Character.Animations[string(f.AnimationState)]
	}
}

func (f *Fighter) Draw(glCtx *webgl.GLContext) {
	glCtx.RenderSprite(f.Animation.GetCurrentFrame(), f.Colliders.HitBox, f.Properties.Specular)

	glCtx.RenderRect(f.Colliders.HitBox, webgl.ColorRed(1.0))
	glCtx.RenderRect(f.Colliders.Attack, webgl.ColorGreen(1.0))
}

func (f *Fighter) Update(deltaTime time.Duration, enemyFighter *Fighter) {
	if f.State == Death && f.Animation.IsEnd {
		return
	}

	var dx, dy float64

	f.handleCooldowns(deltaTime.Seconds())
	f.handleEnemyAttack(enemyFighter.Colliders.Attack, enemyFighter)

	f.Properties.move = false
	if !f.Properties.death && !f.Properties.hit {
		if !f.Properties.attack {
			if f.Control.MoveLeft {
				dx += -1
				f.Properties.Specular = true
			}
			if f.Control.MoveRight {
				dx += 1
				f.Properties.Specular = false
			}
			if f.Control.Jump {
				if !f.HoldingKeys.Jump {
					f.jump()
					f.HoldingKeys.Jump = true
				}
			} else {
				f.HoldingKeys.Jump = false
			}
		}
		if f.Control.Attack {
			if !f.HoldingKeys.Attack {
				f.attack()
				f.HoldingKeys.Attack = true
			}
		} else {
			f.HoldingKeys.Attack = false
		}
	}

	f.move(&dx, enemyFighter.Colliders.HitBox.Center())
	f.gravity(&dy, deltaTime.Seconds())
	f.updatePos(dx, dy, deltaTime.Seconds())

	f.handleState()
	f.updateAnimationState()

	f.Animation.Update(float64(deltaTime.Milliseconds()))
}

func (f *Fighter) move(dx *float64, efc primitives.Vec2) {
	if *dx == 0 {
		if !f.Properties.attack && !f.Properties.death {
			f.handleSpecular(efc)
		}
		f.Properties.move = false
	} else {
		f.Properties.move = true
	}

	*dx *= f.Properties.speed

}

func (f *Fighter) gravity(dy *float64, deltaTime float64) {
	f.Properties.velocity += f.Properties.gravity * deltaTime
	*dy += f.Properties.velocity
}

func (f *Fighter) updatePos(dx, dy, deltaTime float64) {
	dir := primitives.NewVec2(dx, dy)
	offset := dir.MulValue(deltaTime)

	newPos := f.Colliders.HitBox.Move(offset)
	f.Colliders.HitBox.Pos = newPos

	f.handleWorldCollision()

	f.Colliders.Attack = primitives.NewRect(0, 0, 0, 0)

	if f.Properties.attack {
		switch f.State {
		case Attack1:
			if f.Animation.CurrentFrameIndex >= f.Character.Properies.Attack1FrameIndex {
				if !f.Properties.Specular {
					f.Colliders.Attack = primitives.NewRect(
						f.Colliders.HitBox.Right(),
						f.Colliders.HitBox.Top()-f.Character.Properies.Attack1Up,
						f.Character.Properies.Attack1Range,
						f.Character.Properies.Attack1Height,
					)
				} else {
					f.Colliders.Attack = primitives.NewRect(
						f.Colliders.HitBox.Left()-f.Character.Properies.Attack1Range,
						f.Colliders.HitBox.Top()-f.Character.Properies.Attack1Up,
						f.Character.Properies.Attack1Range,
						f.Character.Properies.Attack1Height,
					)
				}
			}

		case Attack2:
			if f.Animation.CurrentFrameIndex >= f.Character.Properies.Attack2FrameIndex {
				if !f.Properties.Specular {
					f.Colliders.Attack = primitives.NewRect(
						f.Colliders.HitBox.Right(),
						f.Colliders.HitBox.Top()-f.Character.Properies.Attack2Up,
						f.Character.Properies.Attack2Range,
						f.Character.Properies.Attack2Height,
					)
				} else {
					f.Colliders.Attack = primitives.NewRect(
						f.Colliders.HitBox.Left()-f.Character.Properies.Attack2Range,
						f.Colliders.HitBox.Top()-f.Character.Properies.Attack2Up,
						f.Character.Properies.Attack2Range,
						f.Character.Properies.Attack2Height,
					)
				}
			}
		}
	}
}

func (f *Fighter) jump() {
	if f.Properties.jump {
		return
	}
	f.Properties.velocity = -f.Properties.jumpSpeed
	f.Properties.jump = true
}

func (f *Fighter) attack() {
	if f.Properties.attackCooldown > 0 || f.Properties.jump {
		return
	}

	if f.Properties.attack {
		f.Properties.comboAttack = true
	} else {
		f.Properties.attack = true
	}
}

func (f *Fighter) handleCooldowns(deltaTime float64) {
	if f.Properties.invulnerableCooldown > 0 {
		f.Properties.invulnerableCooldown -= deltaTime
	}
	if f.Properties.attackCooldown > 0 {
		f.Properties.attackCooldown -= deltaTime
	}
}

func (f *Fighter) handleEnemyAttack(attackCollider primitives.Rect, enemyFighter *Fighter) {
	if f.Properties.invulnerableCooldown > 0 || f.Properties.death || f.Properties.hit || (attackCollider.Width() <= 0 && attackCollider.Height() <= 0) {
		return
	}

	if f.Colliders.HitBox.Intersection(attackCollider) {
		if enemyFighter.State == Attack1 {
			f.Properties.HealthPoints -= enemyFighter.Character.Properies.Attack1Damage
		}
		if enemyFighter.State == Attack2 {
			f.Properties.HealthPoints -= enemyFighter.Character.Properies.Attack2Damage
		}

		if f.Properties.HealthPoints <= 0 {
			f.Properties.HealthPoints = 0
			f.Properties.death = true
		} else {
			f.Properties.hit = true
		}
	}
}

func (f *Fighter) handleState() {
	if f.Properties.death {
		f.State = Death
	} else if f.Properties.hit {
		if f.State != Hit {
			f.State = Hit
		} else if f.Animation.IsEnd {
			f.Properties.hit = false
		}
	} else if f.Properties.attack {
		if f.State != Attack1 && f.State != Attack2 {
			f.State = Attack1
		} else if f.Animation.IsEnd {
			if f.Properties.comboAttack {
				f.State = Attack2
				if f.Animation.IsEnd {
					f.Properties.comboAttack = false
				}
			} else {
				f.Properties.attackCooldown = 0.3
				f.Properties.attack = false
			}
		}
	} else if f.Properties.jump {
		f.State = Jump
	} else if f.Properties.move {
		f.State = Run
	} else {
		f.State = Idle
	}
}

func (f *Fighter) handleWorldCollision() {
	leftWall := 0.0
	rightWall := config.ProgramConf.Window.Width
	floor := config.ProgramConf.Window.Height - 250

	if f.Colliders.HitBox.Left() < leftWall {
		f.Colliders.HitBox.SetLeft(0)
		f.Properties.move = false
	}
	if f.Colliders.HitBox.Right() > rightWall {
		f.Colliders.HitBox.SetRight(rightWall)
		f.Properties.move = false
	}
	if f.Colliders.HitBox.Bottom() > floor {
		f.Properties.jump = false
		f.Properties.velocity = 0
		f.Colliders.HitBox.SetBottom(floor)
	}
}

func (f *Fighter) handleSpecular(center primitives.Vec2) {
	if f.Colliders.HitBox.Center().X > center.X {
		f.Properties.Specular = true
	} else {
		f.Properties.Specular = false
	}
}
