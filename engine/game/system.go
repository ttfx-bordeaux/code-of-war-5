package game

import (
	"engo.io/ecs"
	"engo.io/engo"
)

var (
	AnimationRate float32 = 0.1
	WalkAction    *engo.AnimationAction
	StopAction    *engo.AnimationAction
	DieAction     *engo.AnimationAction
	Actions       []*engo.AnimationAction
)

func Preload() {
	// animation for chicken
	StopAction = &engo.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	WalkAction = &engo.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	DieAction = &engo.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	Actions = []*engo.AnimationAction{DieAction, StopAction, WalkAction}
}

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }
func (*ControlSystem) Pre()         {}
func (*ControlSystem) Post()        {}

func (c *ControlSystem) New(*ecs.World) {}

func (c *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var a *engo.AnimationComponent

	if !entity.Component(&a) {
		return
	}

	if engo.Keys.Get(engo.ArrowRight).Down() {
		a.SelectAnimationByAction(WalkAction)
	} else if engo.Keys.Get(engo.Space).Down() {
		a.SelectAnimationByAction(DieAction)
	} else {
		a.SelectAnimationByAction(StopAction)
	}

}

type WhoopSystem struct {
	goingUp bool
}

func (WhoopSystem) Type() string             { return "WhoopSystem" }
func (WhoopSystem) Priority() int            { return 0 }
func (WhoopSystem) New(w *ecs.World)         {}
func (WhoopSystem) AddEntity(*ecs.Entity)    {}
func (WhoopSystem) RemoveEntity(*ecs.Entity) {}

func (ws *WhoopSystem) Update(dt float32) {
	engo.MasterVolume = 1
	ws.goingUp = false
}

// ecrire un move system