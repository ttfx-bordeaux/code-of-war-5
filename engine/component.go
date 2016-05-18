package main

import (
	"engo.io/engo"
)

//GamePosition game position
type PositionComponent struct {
	Abs int
	Ord int
	idGround int
	Rate float32   // How often entity should move, in seconds.
	nbIteForCpltMove float32
	Change float32 // the time since the last incrementation
	DeltaX float32  // is a split movement regarding rate
	DeltaY float32
	Completed bool // true when position is reached
	NbIteration float32
}

func NewPositionComponent(abs, ord, idGround int) PositionComponent {
	return PositionComponent{
		Abs: abs,
		Ord: ord,
		idGround: idGround,
		Rate: 0.1,
		nbIteForCpltMove: 100,
		NbIteration: -1,
	}
}

func (p *PositionComponent) changePositionTo(abs, ord int) {
	p.Abs = abs
	p.Ord = ord
	p.Completed = false
	p.NbIteration = -1
}

func (p PositionComponent) toPoint() engo.Point {
	padding := p.idGround * (nbTilesAbs*TileWidth + padRight)
	return engo.Point{X: float32(p.Abs*TileWidth + padding), Y: float32(p.Ord * TileHeight)}
}

func (p *PositionComponent) InitDelta(from engo.Point) {
	to := p.toPoint()
	to.Subtract(from)
	p.DeltaX = to.X / p.nbIteForCpltMove
	p.DeltaY = to.Y / p.nbIteForCpltMove
	p.NbIteration = p.nbIteForCpltMove
}

func (p PositionComponent) Type() string {
	return "PositionComponenet"
}
