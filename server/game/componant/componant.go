package componant

// Identifier componant
type Identifier interface {
	ID() int
}

// Position of an entity
type Position struct{ x, y int }

// Positioner componant
type Positioner interface {
	Position() Position
}

// Velociter componant
type Velociter interface {
	Velocity() int
}

// Pricer componant
type Pricer interface {
	Price() int
}

// Life of an entity
type Life struct{ current, max int }

// Lifer componant
type Lifer interface {
	Life() Life
}

// Defencer componant
type Defencer interface {
	Defence() int
}

// Attaquer componant
type Attaquer interface {
	Attack() (damage, distance int)
}
