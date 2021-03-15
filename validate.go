package tienlen_bot

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func isDubs(card1, card2 *Card) bool {
	return card1.rank == card2.rank
}

func isTrips(card1, card2, card3 *Card) bool {
	return card1.rank == card2.rank && card2.rank == card3.rank
}

func isQuads(card1, card2, card3, card4 *Card) bool {
	return isTrips(card1, card2, card3) && card3.rank == card4.rank
}

func isSequence(cards []*Card) bool {
	if len(cards) <= 2 {
		return false
	}
	if cards[len(cards) - 1].rank == Two {
		return false
	}
	for i := 0; i < len(cards) - 1; i++ {
		if cards[i].rank + 1 != cards[i + 1].rank {
			return false
		}
	}
	return true
}

func isHomogeneitySequence(cards []*Card) bool {
	for _, card := range cards {
		if card.suit != cards[0].suit {
			return false
		}
	}
	return true
}

func isTwoConsecutivePairs(dubs1, dubs2 *Dubs) bool {
	if dubs2.rank == Two {
		return false
	}
	return dubs1.rank + 1 == dubs2.rank
}

func isThreeConsecutivePairs(dubs1, dubs2, dubs3 *Dubs) bool {
	if dubs3.rank == Two {
		return false
	}
	return dubs1.rank + 1 == dubs2.rank && dubs2.rank + 1 == dubs3.rank
}

func isFourConsecutivePairs(dubs1, dubs2, dubs3, dubs4 *Dubs) bool {
	if dubs4.rank == Two {
		return false
	}
	return dubs1.rank + 1 == dubs2.rank && dubs2.rank + 1 == dubs3.rank && dubs3.rank + 1 == dubs4.rank
}

func compareRank(r1, r2 Rank) int {
	if r1 > r2 {
		return 1
	}
	if r1 < r2 {
		return -1
	}
	return 0
}

func compareSuit(s1, s2 Suit) int {
	if s1 > s2 {
		return 1
	}
	if s1 < s2 {
		return -1
	}
	return 0
}

func compareCard(c1, c2 *Card) int {
	if c1.rank == c2.rank {
		return compareSuit(c1.suit, c2.suit)
	}
	return compareRank(c1.rank, c2.rank)
}

func compareSingleCard(c1, c2 *SingleCard) int {
	return compareCard(c1.card, c2.card)
}

func compareDubs(c1, c2 *Dubs) int {
	if c1.rank != c2.rank {
		return compareRank(c1.rank, c2.rank)
	}
	if c1.maxSuit != c2.maxSuit {
		return compareSuit(c1.maxSuit, c2.maxSuit)
	}
	return compareSuit(c1.minSuit, c2.minSuit)
}

func compareTrips(c1, c2 *Trips) int {
	if c1.rank != c2.rank {
		return compareRank(c1.rank, c2.rank)
	}
	c := compareCard(c1.card3, c2.card3)
	if c != 0 {
		return c
	}
	c = compareCard(c1.card2, c2.card2)
	if c != 0 {
		return c
	}
	return compareCard(c1.card1, c2.card1)
}

func compareQuads(c1, c2 *Quads) int {
	return compareRank(c1.rank, c2.rank)
}

func compareSequence(c1, c2 *Sequence) int {
	if len(c1.cardList) > len(c2.cardList) {
		return 1
	}
	if len(c1.cardList) < len(c2.cardList) {
		return -1
	}
	if c1.minRank != c2.minRank {
		return compareRank(c1.minRank, c2.minRank)
	}
	for i := len(c1.cardList) - 1; i >= 0; i-- {
		if c1.cardList[i].suit != c2.cardList[i].suit {
			return compareSuit(c1.cardList[i].suit, c2.cardList[i].suit)
		}
	}
	return 0
}

func compareTwoConsecutivePairs(c1, c2 *TwoConsecutivePairs) int {
	if c1.minRank != c2.minRank {
		return compareRank(c1.minRank, c2.minRank)
	}
	if c1.dubs2.equals(c2.dubs2) {
		return compareDubs(c1.dubs1, c2.dubs1)
	}
	return compareDubs(c1.dubs2, c2.dubs2)
}

func compareThreeConsecutivePairs(c1, c2 *ThreeConsecutivePairs) int {
	if c1.minRank != c2.minRank {
		return compareRank(c1.minRank, c2.minRank)
	}
	c := compareDubs(c1.dubs3, c2.dubs3)
	if c != 0 {
		return c
	}
	c = compareDubs(c1.dubs2, c2.dubs2)
	if c != 0 {
		return c
	}
	return compareDubs(c1.dubs1, c2.dubs1)
}

func compareFourConsecutivePairs(c1, c2 *FourConsecutivePairs) int {
	if c1.minRank != c2.minRank {
		return compareRank(c1.minRank, c2.minRank)
	}
	c := compareDubs(c1.dubs4, c2.dubs4)
	if c != 0 {
		return c
	}
	c = compareDubs(c1.dubs3, c2.dubs3)
	if c != 0 {
		return c
	}
	c = compareDubs(c1.dubs2, c2.dubs2)
	if c != 0 {
		return c
	}
	return compareDubs(c1.dubs1, c2.dubs1)
}

func isNil(i interface{}) bool {
	if i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()) {
		return true
	}
	return false
}

func notNil(i interface{}) bool {
	return !isNil(i)
}