package componant

// Tower to defence a base
type Tower struct {
	Identifier
	Positioner
	Pricer
	Healther
	Defencer
	Attaquer
	id, x, y, price, currentLife, maxLife, defence, damage, distance int
}

func (t Tower) ID() int {
	return t.id
}

func (t Tower) Position() (int, int) {
	return t.x, t.y
}

func (t Tower) Price() int {
	return t.price
}

func (t Tower) Health() (int, int) {
	return t.currentLife, t.maxLife
}

func (t Tower) Defence() int {
	return t.defence
}

func (t Tower) Attack() (int, int) {
	return t.damage, t.distance
}
