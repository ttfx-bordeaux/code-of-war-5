package game

//GamePosition game position
type PositionComponent struct {
	Abs int
	Ord int
}

func NewPositionComponent(abs, ord int) *PositionComponent {
	return &PositionComponent{
		Abs: abs,
		Ord: ord,
	}
}

func (p PositionComponent) Type() string {
	return "PositionComponenet"
}
