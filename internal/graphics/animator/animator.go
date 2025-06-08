//go:build js

package animation

type Animator struct {
	Animations map[string]*Animation
}

func (a *Animator) LoadAnimations() {

}

func (a *Animator) GetAnimations(name string) *Animation {
	return a.Animations[name]
}
