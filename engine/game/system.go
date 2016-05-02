package game

import (
	"engo.io/ecs"
	"engo.io/engo"
)

var (
	AnimationRate float32 = 0.1
	WalkAction    *engo.Animation
	StopAction    *engo.Animation
	SkillAction   *engo.Animation
	Actions       []*engo.Animation
)

func Preload() {
	// animation for chicken
	StopAction = &engo.Animation{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	WalkAction = &engo.Animation{Name: "move", Frames: []int{11, 12, 13, 14, 15}, Loop: true}
	SkillAction = &engo.Animation{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	Actions = []*engo.Animation{SkillAction, StopAction, WalkAction}
}

type controlEntity struct {
	*ecs.BasicEntity
	*engo.AnimationComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, anim *engo.AnimationComponent) {
	c.entities = append(c.entities, controlEntity{basic, anim})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if engo.Keys.Get(engo.ArrowRight).Down() {
			e.AnimationComponent.SelectAnimationByAction(WalkAction)
		} else if engo.Keys.Get(engo.Space).Down() {
			e.AnimationComponent.SelectAnimationByAction(SkillAction)
		}
	}
}

// ecrire un move system