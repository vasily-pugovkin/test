package tienlen_bot

import "fmt"

type Quads struct {
	card1 *Card
	card2 *Card
	card3 *Card
	card4 *Card
	rank  Rank
}

func NewQuads(card1, card2, card3, card4 *Card) *Quads {
	if !isQuads(card1, card2, card3, card4) {
		panic("invalid quads")
	}
	return &Quads{
		card1: card1,
		card2: card2,
		card3: card3,
		card4: card4,
		rank:  card1.rank,
	}
}

func (q *Quads) Kind() CombinationKind {
	return CombinationQuads
}

func (q *Quads) Equals(combination Combination) bool {
	if combination.Kind() != CombinationQuads {
		return false
	}
	return q.rank == combination.(*Quads).rank
}

func (q *Quads) Cards() []*Card {
	return []*Card{q.card1, q.card2, q.card3, q.card4}
}

func (q *Quads) Defeats(combination Combination) bool {
	switch combination.Kind() {
	case CombinationQuads:
		return q.rank > combination.(*Quads).rank
	case CombinationThreeConsecutivePairs:
		return true
	case CombinationDubs:
		return combination.(*Dubs).rank == Two
	case CombinationSingle:
		return combination.(*SingleCard).card.rank == Two
	}
	return false
}

func (q *Quads) Copy() Combination {
	return NewQuads(q.card1, q.card2, q.card3, q.card4)
}

func (q *Quads) String() string {
	return fmt.Sprintf("[%s %s %s %s]", q.card1, q.card2, q.card3, q.card4)
}