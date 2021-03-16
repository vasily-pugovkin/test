package tienlen_bot

type Pass struct {
}

func NewPass() *Pass {
	return &Pass{}
}

func (p *Pass) Kind() CombinationKind {
	return CombinationPass
}

func (p *Pass) Equals(combination Combination) bool {
	return combination.Kind() == CombinationPass
}

func (p *Pass) Cards() []*Card {
	return []*Card{}
}

func (p *Pass) Defeats(combination Combination) bool {
	return false
}

func (p *Pass) Copy() Combination {
	return p
}

func (p *Pass) String() string {
	return "pass"
}

