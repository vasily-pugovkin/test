package tienlen_bot

import "fmt"

type FourConsecutivePairs struct {
	dubs1   *Dubs
	dubs2   *Dubs
	dubs3   *Dubs
	dubs4   *Dubs
	maxSuit Suit
	minRank Rank
}

func NewFourConsecutivePairs(dubs1, dubs2, dubs3, dubs4 *Dubs) *FourConsecutivePairs {
	if !isFourConsecutivePairs(dubs1, dubs2, dubs3, dubs4) {
		panic("invalid four consecutive pairs")
	}
	return &FourConsecutivePairs{
		dubs1:   dubs1,
		dubs2:   dubs2,
		dubs3:   dubs3,
		dubs4:   dubs4,
		maxSuit: dubs4.maxSuit,
		minRank: dubs1.rank,
	}
}

func (f *FourConsecutivePairs) kind() CombinationKind {
	return CombinationFourConsecutivePairs
}

func (f *FourConsecutivePairs) equals(combination Combination) bool {
	if combination.kind() != CombinationFourConsecutivePairs {
		return false
	}
	o := combination.(*FourConsecutivePairs)
	return f.dubs1.equals(o.dubs1) &&
		f.dubs2.equals(o.dubs2) &&
		f.dubs3.equals(o.dubs3) &&
		f.dubs4.equals(o.dubs4)
}

func (f *FourConsecutivePairs) cards() []*Card {
	return []*Card{
		f.dubs1.card1, f.dubs1.card2,
		f.dubs2.card1, f.dubs2.card2,
		f.dubs3.card1, f.dubs3.card2,
		f.dubs4.card1, f.dubs4.card2,
	}
}

func (f *FourConsecutivePairs) defeats(combination Combination) bool {
	switch combination.kind() {
	case CombinationFourConsecutivePairs:
		o := combination.(*FourConsecutivePairs)
		if f.minRank == o.minRank {
			return f.maxSuit > o.maxSuit
		}
		return f.minRank > o.minRank
	case CombinationQuads:
		return true
	case CombinationThreeConsecutivePairs:
		return true
	case CombinationDubs:
		return combination.(*Dubs).rank == Two
	case CombinationSingle:
		return combination.(*SingleCard).card.rank == Two
	}
	return false
}

func (f *FourConsecutivePairs) copy() Combination {
	return NewFourConsecutivePairs(f.dubs1.copy().(*Dubs), f.dubs2.copy().(*Dubs),
		f.dubs3.copy().(*Dubs), f.dubs4.copy().(*Dubs))
}

func (f *FourConsecutivePairs) String() string {
	return fmt.Sprintf("{%s %s %s %s", f.dubs1, f.dubs2, f.dubs3, f.dubs4)
}
