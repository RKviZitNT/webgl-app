//go:build js

package fighter

import (
	"time"
	"webgl-app/internal/config"
	"webgl-app/internal/game/character"
	"webgl-app/internal/graphics/animation"
	"webgl-app/internal/graphics/primitives"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/net/message"
)

type FighterState string

const (
	Idle    FighterState = "idle"
	Walk    FighterState = "walk"
	Attack1 FighterState = "attack1"
	Attack2 FighterState = "attack2"
	Death   FighterState = "death"
	Hurt    FighterState = "hurt"
	Jump    FighterState = "jump"
)

type FighterPropertys struct {
	Specular    bool
	dir         primitives.Vec2
	speed       float64
	jumpSpeed   float64
	velocity    float64
	gravity     float64
	attackRange float64
	jump        bool
	move        bool
	attack1     bool
	attack2     bool
}

type Fighter struct {
	State          FighterState
	Collider       primitives.Rect
	AttackCollider primitives.Rect
	Character      *character.Character
	Animation      animation.Animation
	Propertys      FighterPropertys
	Control        *message.FighterControl
}

func NewFighter(char *character.Character, collider primitives.Rect) *Fighter {
	propertys := FighterPropertys{
		dir:         primitives.Vec2{},
		jumpSpeed:   1500,
		speed:       400,
		velocity:    0,
		gravity:     5000,
		attackRange: 100,
		jump:        false,
		move:        false,
		attack1:     false,
		attack2:     false,
	}

	fighter := Fighter{
		Collider:       collider,
		AttackCollider: primitives.NewRect(0, 0, collider.Width()/2+propertys.attackRange, collider.Height()),
		Character:      char,
		Propertys:      propertys,
		Control:        nil,
	}
	fighter.SetState(Idle)

	return &fighter
}

func (f *Fighter) SetState(state FighterState) {
	if f.State != state {
		f.State = state
		jsfunc.LogInfo(f.Collider.Pos.X)
		f.Animation = *f.Character.Animations[string(state)]
	}
}

func (f *Fighter) Draw(glCtx *webgl.GLContext) {
	glCtx.RenderSprite(f.Animation.GetCurrentFrame(), f.Collider, f.Propertys.Specular)

	glCtx.RenderRect(f.Collider, webgl.ColorRed(1.0))
	glCtx.RenderRect(f.AttackCollider, webgl.ColorGreen(1.0))
}

func (f *Fighter) Update(keys map[string]bool, deltaTime time.Duration) {
	f.handleState()

	var dx, dy float64

	moveLeft := func() {
		dx = 1
		f.Propertys.Specular = false
		f.Propertys.move = true
	}
	moveRight := func() {
		dx = -1
		f.Propertys.Specular = true
		f.Propertys.move = true
	}
	jump := func() {
		f.jump()
	}

	f.Propertys.move = false
	if f.Control == nil {
		if keys["KeyD"] {
			moveLeft()
		}
		if keys["KeyA"] {
			moveRight()
		}
		if keys["Space"] {
			jump()
		}
	} else {
		if f.Control.MoveLeft {
			moveLeft()
		}
		if f.Control.MoveRight {
			moveRight()
		}
		if f.Control.Jump {
			jump()
		}
	}

	f.move(dx, dy, deltaTime.Seconds())

	f.Animation.Update(float64(deltaTime.Milliseconds()))

	f.handleCollision()
}

func (f *Fighter) move(dx, dy, deltaTime float64) {
	f.Propertys.velocity += f.Propertys.gravity * deltaTime

	dx *= f.Propertys.speed
	dy += f.Propertys.velocity

	dir := primitives.NewVec2(dx, dy)
	offset := dir.MulValue(deltaTime)

	newPos := f.Collider.Move(offset)
	f.Collider.Pos = newPos

	f.AttackCollider.Pos = primitives.NewVec2(f.Collider.Center().X, f.Collider.Top())
}

func (f *Fighter) jump() {
	if !f.Propertys.jump {
		f.Propertys.velocity = -f.Propertys.jumpSpeed
		f.Propertys.jump = true
	}
}

func (f *Fighter) attack() {

}

func (f *Fighter) handleState() {
	if f.Propertys.jump {
		f.SetState(Jump)
	} else if f.Propertys.attack1 {
		f.SetState(Attack1)
	} else if f.Propertys.attack2 {
		f.SetState(Attack2)
	} else if f.Propertys.move {
		f.SetState(Walk)
	} else {
		f.SetState(Idle)
	}
}

func (f *Fighter) handleCollision() {
	leftWall := 0.0
	rightWall := config.ProgramConf.Window.Width
	floor := config.ProgramConf.Window.Height - 250

	if f.Collider.Left() < leftWall {
		f.Collider.SetLeft(0)
	}
	if f.Collider.Right() > rightWall {
		f.Collider.SetRight(rightWall)
	}
	if f.Collider.Bottom() > floor {
		f.Propertys.jump = false
		f.Propertys.velocity = 0
		f.Collider.SetBottom(floor)
	}
}

func (f *Fighter) HandleSpecular(center primitives.Vec2) {
	if f.Collider.Center().X > center.X {
		f.Propertys.Specular = true
	} else {
		f.Propertys.Specular = false
	}
}
