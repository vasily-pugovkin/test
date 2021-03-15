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

func (t *TwoConsecutivePairs) kind() CombinationKind {
	return CombinationTwoConsecutivePairs
}

func (t *TwoConsecutivePairs) equals(combination Combination) bool {
	if combination.kind() != CombinationTwoConsecutivePairs {
		return false
	}
	o := combination.(*TwoConsecutivePairs)
	return t.dubs1.equals(o.dubs1) && t.dubs2.equals(o.dubs2)
}

func (t *TwoConsecutivePairs) cards() []*Card {
	return []*Card{t.dubs1.card1, t.dubs1.card2, t.dubs2.card1, t.dubs2.card2}
}

func (t *TwoConsecutivePairs) defeats(combination Combination) bool {
	if combination.kind() != CombinationTwoConsecutivePairs {
		return false
	}
	o := combination.(*TwoConsecutivePairs)
	if t.minRank == o.minRank {
		return t.maxSuit > o.maxSuit
	}
	return t.minRank > o.minRank
}

func (t *TwoConsecutivePairs) copy() Combination {
	return NewTwoConsecutivePairs(t.dubs1.copy().(*Dubs), t.dubs2.copy().(*Dubs))
}

func (t *TwoConsecutivePairs) String() string {
	return fmt.Sprintf("{%s %s}", t.dubs1, t.dubs2)
}

