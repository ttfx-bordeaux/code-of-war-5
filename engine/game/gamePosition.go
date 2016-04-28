package game

//GamePosition game position
type GamePosition struct {
	Abs int
	Ord int
}

func NewGamePositionComponent(abs, ord int) *GamePosition {
	return &GamePosition{
		Abs: abs,
		Ord: ord,
	}
}

func (p GamePosition) Type() string {
	return "GamePositionComponenet"
}
