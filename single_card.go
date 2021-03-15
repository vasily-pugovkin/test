package tienlen_bot

type SingleCard struct {
	card *Card
}

func NewSingleCard(card *Card) *SingleCard {
	return &SingleCard{card: card}
}

func (s *SingleCard) kind() CombinationKind {
	return CombinationSingle
}

func (s *SingleCard) equals(combination Combination) bool {
	return combination.kind() == CombinationSingle && combination.(*SingleCard).card.equals(s.card)
}

func (s *SingleCard) cards() []*Card {
	return []*Card{s.card}
}

func (s *SingleCard) defeats(combination Combination) bool {
	if combination.kind() != CombinationSingle {
		return false
	}
	c := combination.(*SingleCard)
	if c.card.rank > s.card.rank {
		return false
	}
	if c.card.rank < s.card.rank {
		return true
	}
	if c.card.rank >= Six {
		return s.card.suit > c.card.suit
	}
	return true
}

func (s *SingleCard) copy() Combination {
	return NewSingleCard(&Card{
		rank: s.card.rank,
		suit: s.card.suit,
	})
}

func (s *SingleCard) String() string {
	return s.card.String()
}