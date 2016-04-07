package componant

// Identifier componant
type Identifier interface {
	ID() int
}

// Positioner componant
type Positioner interface {
	Position() (x, y int)
}

// Velociter componant
type Velociter interface {
	Velocity() int
}

// Pricer componant
type Pricer interface {
	Price() int
}

// Healther componant
type Healther interface {
	Health() (currentLife, maxLife int)
}

// Defencer componant
type Defencer interface {
	Defence() int
}

// Attaquer componant
type Attaquer interface {
	Attack() (damage, distance int)
}
