package componant

// Tower to defence a base
type Tower struct {
	life                                 Life
	position                             Position
	id, price, defence, damage, distance int
}

func (t Tower) ID() int {
	return t.id
}

func (t Tower) Position() Position {
	return t.position
}

func (t Tower) Price() int {
	return t.price
}

func (t Tower) Life() Life {
	return t.life
}

func (t Tower) Defence() int {
	return t.defence
}

func (t Tower) Attack() (int, int) {
	return t.damage, t.distance
}
