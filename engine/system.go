package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"log"
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

type moveEntity struct {
	*ecs.BasicEntity
	*PositionComponent
	*engo.SpaceComponent
}

type MoveSystem struct {
	entities []moveEntity
}

func (m *MoveSystem) Add(basic *ecs.BasicEntity, pos *PositionComponent, space *engo.SpaceComponent) {
	m.entities = append(m.entities, moveEntity{basic, pos, space})
}

func (m *MoveSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range m.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		m.entities = append(m.entities[:delete], m.entities[delete+1:]...)
	}
}

func (m *MoveSystem) Update(dt float32) {
	for _, e := range m.entities {
		if e.PositionComponent.Completed == false {
			if e.PositionComponent.NbIteration == -1 {
				log.Printf("init movement for id %v", e.BasicEntity.ID())
				e.PositionComponent.InitDelta(e.SpaceComponent.Position)
			}
			e.PositionComponent.Change += dt
			if e.PositionComponent.Change >= e.PositionComponent.Rate {
				e.SpaceComponent.Position.X += e.PositionComponent.DeltaX
				e.SpaceComponent.Position.Y += e.PositionComponent.DeltaY
				e.PositionComponent.Change = 0
				e.PositionComponent.NbIteration -= 1
			}
			if e.PositionComponent.NbIteration <= 0 {
				log.Printf("stop movement for id %v", e.BasicEntity.ID())
				e.PositionComponent.Completed = true
			}
		}
	}
}
