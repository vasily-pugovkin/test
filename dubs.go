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

func (d *Dubs) kind() CombinationKind {
	return CombinationDubs
}

func (d *Dubs) equals(combination Combination) bool {
	if combination.kind() != CombinationDubs {
		return false
	}
	dub := combination.(*Dubs)
	return dub.card1.equals(d.card1) && dub.card2.equals(d.card2)
}

func (d *Dubs) cards() []*Card {
	return []*Card{d.card1, d.card2}
}

func (d *Dubs) defeats(combination Combination) bool {
	if combination.kind() != CombinationDubs {
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

func (d *Dubs) copy() Combination {
	return NewDubs(d.card1, d.card2)
}

func (d *Dubs) String() string {
	return fmt.Sprintf("[%s %s]", d.card1, d.card2)
}