package tienlen_bot

import "fmt"

type TwoConsecutivePairs struct {
	dubs1 *Dubs
	dubs2 *Dubs
	minRank Rank
	maxSuit Suit
}

func NewTwoConsecutivePairs(dubs1, dubs2 *Dubs) *TwoConsecutivePairs {
	if !isTwoConsecutivePairs(dubs1, dubs2) {
		panic("invalid two consecutive pairs")
	}
	return &TwoConsecutivePairs{
		dubs1:   dubs1,
		dubs2:   dubs2,
		minRank: dubs1.rank,
		maxSuit: dubs2.maxSuit,
	}
}

func (t *TwoConsecutivePairs) Kind() CombinationKind {
	return CombinationTwoConsecutivePairs
}

func (t *TwoConsecutivePairs) Equals(combination Combination) bool {
	if combination.Kind() != CombinationTwoConsecutivePairs {
		return false
	}
	o := combination.(*TwoConsecutivePairs)
	return t.dubs1.Equals(o.dubs1) && t.dubs2.Equals(o.dubs2)
}

func (t *TwoConsecutivePairs) Cards() []*Card {
	return []*Card{t.dubs1.card1, t.dubs1.card2, t.dubs2.card1, t.dubs2.card2}
}

func (t *TwoConsecutivePairs) Defeats(combination Combination) bool {
	if combination.Kind() != CombinationTwoConsecutivePairs {
		return false
	}
	o := combination.(*TwoConsecutivePairs)
	if t.minRank == o.minRank {
		return t.maxSuit > o.maxSuit
	}
	return t.minRank > o.minRank
}

func (t *TwoConsecutivePairs) Copy() Combination {
	return NewTwoConsecutivePairs(t.dubs1.Copy().(*Dubs), t.dubs2.Copy().(*Dubs))
}

func (t *TwoConsecutivePairs) String() string {
	return fmt.Sprintf("{%s %s}", t.dubs1, t.dubs2)
}

