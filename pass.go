package tienlen_bot

type Pass struct {
}

func NewPass() *Pass {
	return &Pass{}
}

func (p *Pass) kind() CombinationKind {
	return CombinationPass
}

func (p *Pass) equals(combination Combination) bool {
	return combination.kind() == CombinationPass
}

func (p *Pass) cards() []*Card {
	return []*Card{}
}

func (p *Pass) defeats(combination Combination) bool {
	return false
}

func (p *Pass) copy() Combination {
	return p
}

func (p *Pass) String() string {
	return "pass"
}

