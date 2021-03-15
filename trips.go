package tienlen_bot

import "fmt"

type Trips struct {
	card1 *Card
	card2 *Card
	card3 *Card
	rank Rank
}

func NewTrips(card1, card2, card3 *Card) *Trips {
	if !isTrips(card1, card2, card3) {
		panic("invalid trips")
	}
	return &Trips{
		card1: card1,
		card2: card2,
		card3: card3,
		rank:  card1.rank,
	}
}

func (t *Trips) kind() CombinationKind {
	return CombinationTrips
}

func (t *Trips) equals(combination Combination) bool {
	if combination.kind() != CombinationTrips {
		return false
	}
	trips := combination.(*Trips)
	return t.card1.equals(trips.card1) && t.card2.equals(trips.card2) && t.card3.equals(trips.card3)
}

func (t *Trips) cards() []*Card {
	return []*Card{t.card1, t.card2, t.card3}
}

func (t *Trips) defeats(combination Combination) bool {
	if combination.kind() != CombinationTrips {
		return false
	}
	trips := combination.(*Trips)
	return t.rank > trips.rank
}

func (t *Trips) copy() Combination {
	return NewTrips(t.card1, t.card2, t.card3)
}

func (t *Trips) String() string {
	return fmt.Sprintf("[%s %s %s]", t.card1, t.card1, t.card3)
}