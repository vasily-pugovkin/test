package tienlen_bot

import "fmt"

type ThreeConsecutivePairs struct {
	dubs1 *Dubs
	dubs2 *Dubs
	dubs3 *Dubs
	maxSuit Suit
	minRank Rank
}

func NewThreeConsecutivePairs(dubs1, dubs2, dubs3 *Dubs) *ThreeConsecutivePairs {
	if !isThreeConsecutivePairs(dubs1, dubs2, dubs3) {
		panic("invalid three consecutive pairs")
	}
	return &ThreeConsecutivePairs{
		dubs1:   dubs1,
		dubs2:   dubs2,
		dubs3:   dubs3,
		maxSuit: dubs3.maxSuit,
		minRank: dubs1.rank,
	}
}

func (t *ThreeConsecutivePairs) kind() CombinationKind {
	return CombinationThreeConsecutivePairs
}

func (t *ThreeConsecutivePairs) equals(combination Combination) bool {
	if combination.kind() != CombinationThreeConsecutivePairs {
		return false
	}
	o := combination.(*ThreeConsecutivePairs)
	return t.dubs1.equals(o.dubs1) && t.dubs2.equals(o.dubs2) && t.dubs3.equals(o.dubs3)
}

func (t *ThreeConsecutivePairs) cards() []*Card {
	return []*Card {
		t.dubs1.card1, t.dubs1.card2,
		t.dubs2.card1, t.dubs2.card2,
		t.dubs3.card1, t.dubs3.card2,
	}
}

func (t *ThreeConsecutivePairs) defeats(combination Combination) bool {
	if combination.kind() != CombinationThreeConsecutivePairs{
		return false
	}
	o := combination.(*ThreeConsecutivePairs)
	if t.minRank == o.minRank {
		return t.maxSuit > o.maxSuit
	}
	return t.minRank < o.minRank
}

func (t *ThreeConsecutivePairs) copy() Combination {
	return NewThreeConsecutivePairs(t.dubs1.copy().(*Dubs),
		t.dubs2.copy().(*Dubs),
		t.dubs3.copy().(*Dubs))
}

func (t *ThreeConsecutivePairs) String() string {
	return fmt.Sprintf("{%s %s %s}", t.dubs1, t.dubs2, t.dubs3)
}