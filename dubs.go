package tienlen_bot

import "fmt"

type Dubs struct {
	card1 *Card
	card2 *Card
	maxSuit Suit
	minSuit Suit
	rank Rank
}

func NewDubs(card1, card2 *Card) *Dubs {
	if !isDubs(card1, card2) {
		panic("invalid dubs")
	}
	return &Dubs{
		card1:   card1,
		card2:   card2,
		maxSuit: MaxSuit(card1.suit, card2.suit),
		minSuit: MinSuit(card1.suit, card2.suit),
		rank:    card1.rank,
	}
}

func (d *Dubs) Kind() CombinationKind {
	return CombinationDubs
}

func (d *Dubs) Equals(combination Combination) bool {
	if combination.Kind() != CombinationDubs {
		return false
	}
	dub := combination.(*Dubs)
	return dub.card1.equals(d.card1) && dub.card2.equals(d.card2)
}

func (d *Dubs) Cards() []*Card {
	return []*Card{d.card1, d.card2}
}

func (d *Dubs) Defeats(combination Combination) bool {
	if combination.Kind() != CombinationDubs {
		return false
	}
	dub := combination.(*Dubs)
	if d.rank > dub.rank {
		return true
	}
	if d.rank < dub.rank {
		return false
	}
	if d.rank >= Six {
		return d.maxSuit > dub.maxSuit
	}
	return true
}

func (d *Dubs) Copy() Combination {
	return NewDubs(d.card1, d.card2)
}

func (d *Dubs) String() string {
	return fmt.Sprintf("[%s %s]", d.card1, d.card2)
}