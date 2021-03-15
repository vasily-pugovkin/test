package tienlen_bot

import (
	"sort"
)

type CombinationKind uint8

const (
	CombinationSingle CombinationKind = iota
	CombinationDubs
	CombinationTrips
	CombinationQuads
	CombinationSequence
	CombinationTwoConsecutivePairs
	CombinationThreeConsecutivePairs
	CombinationFourConsecutivePairs
	CombinationPass
)

type Combination interface {
	kind() CombinationKind
	equals(combination Combination) bool
	cards() []*Card
	defeats(combination Combination) bool
	copy() Combination
	String() string
}

func containsCard(cards []*Card, card *Card) bool {
	for _, c := range cards {
		if c.equals(card) {
			return true
		}
	}
	return false
}

func containsRank(cards []*Card, rank Rank) bool {
	for _, c := range cards {
		if c.rank == rank {
			return true
		}
	}
	return false
}

func GetDubs(cards []*Card) []*Dubs {
	dubs := []*Dubs{}
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if isDubs(cards[i], cards[j]) {
				dubs = append(dubs, NewDubs(cards[i], cards[j]))
			}
		}
	}
	sort.Slice(dubs, func(i, j int) bool {
		return compareDubs(dubs[i], dubs[j]) < 0
	})
	return dubs
}

func GetTrips(cards []*Card) []*Trips {
	trips := []*Trips{}
	count := make([]int, 13)
	for _, card := range cards {
		count[card.rank]++
	}
	for rank := Three; rank <= Two; rank++ {
		if count[rank] < 3 {
			continue
		}
		if count[rank] == 3 {
			list := getAllCardsWithRank(cards, rank)
			trips = append(trips, NewTrips(list[0], list[1], list[2]))
		} else if count[rank] == 4 {
			card1 := &Card{rank: rank, suit: Spade}
			card2 := &Card{rank: rank, suit: Club}
			card3 := &Card{rank: rank, suit: Diamond}
			card4 := &Card{rank: rank, suit: Heart}
			trips = append(trips, NewTrips(card1, card2, card3))
			trips = append(trips, NewTrips(card1, card2, card4))
			trips = append(trips, NewTrips(card1, card3, card4))
			trips = append(trips, NewTrips(card2, card3, card4))
		}
	}
	sort.Slice(trips, func(i, j int) bool {
		return compareTrips(trips[i], trips[j]) < 0
	})
	return trips
}

func GetQuads(cards []*Card) []*Quads {
	quads := []*Quads{}
	count := make([]int, 13)
	for _, card := range cards {
		count[card.rank]++
	}
	for rank := Three; rank <= Two; rank++ {
		if count[rank] == 4 {
			quads = append(quads, NewQuads(&Card{rank: rank, suit: Spade}, &Card{rank: rank,suit: Club},
			&Card{rank: rank, suit: Diamond}, &Card{rank: rank, suit: Heart}))
		}
	}
	sort.Slice(quads, func(i, j int) bool {
		return compareQuads(quads[i], quads[j]) < 0
	})
	return quads
}

func GetSequence(cards []*Card) []*Sequence {
	sequences := []*Sequence{}
	checkList := make([][]*Card, 13)
	for i := 0; i < len(cards); i++ {
		list := checkList[cards[i].rank]
		if list == nil {
			checkList[cards[i].rank] = []*Card{cards[i]}
		} else {
			checkList[cards[i].rank] = append(list, cards[i])
		}
	}
	for i := Three; i < Queen; i++ {
		list := checkList[i]
		if list == nil {
			continue
		}
		var lastSequenceList [][]*Card
		if l := checkList[i + 1]; l == nil {
			continue
		} else {
			lastSequenceList = make([][]*Card, len(list))
			for j := 0; j < len(list); j++ {
				for t := 0; t < len(l); t++ {
					lastSequenceList[j] = []*Card{list[j], l[t]}
				}
			}
		}
		for j := i + 2; j < Ace; j++ {
			l := checkList[j]
			if l == nil {
				break
			}
			currentSequenceList := [][]*Card{}
			for t := 0; t < len(l); t++ {
				for r := 0; r < len(lastSequenceList); r++ {
					sequence := append(lastSequenceList[r], l[t])
					sequences = append(sequences, NewSequence(sequence))
					currentSequenceList = append(currentSequenceList, sequence)
				}
			}
			lastSequenceList = currentSequenceList
		}
	}
	sort.Slice(sequences, func(i, j int) bool {
		return compareSequence(sequences[i], sequences[j]) < 0
	})
	return sequences
}

func GetTwoConsecutivePairs(cards []*Card) []*TwoConsecutivePairs {
	pairs := []*TwoConsecutivePairs{}
	dubs := GetDubs(cards)
	if len(dubs) < 2 {
		return pairs
	}
	for i := 0; i < len(dubs) - 1; i++ {
		for j := i + 1; j < len(dubs); j++ {
			if isTwoConsecutivePairs(dubs[i], dubs[j]) {
				pairs = append(pairs, NewTwoConsecutivePairs(dubs[i], dubs[j]))
			}
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return compareTwoConsecutivePairs(pairs[i], pairs[j]) < 0
	})
	return pairs
}

func GetThreeConsecutivePairs(cards []*Card) []*ThreeConsecutivePairs {
	pairs := []*ThreeConsecutivePairs{}
	dubs := GetDubs(cards)
	if len(dubs) < 3 {
		return pairs
	}
	for i := 0; i < len(dubs) - 2; i++ {
		for j := i + 1; j < len(dubs) - 1; j++ {
			for t := j + 1; t < len(dubs); t++ {
				if isThreeConsecutivePairs(dubs[i], dubs[j], dubs[t]) {
					pairs = append(pairs, NewThreeConsecutivePairs(dubs[i], dubs[j], dubs[t]))
				}
			}
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return compareThreeConsecutivePairs(pairs[i], pairs[j]) < 0
	})
	return pairs
}

func GetFourConsecutivePairs(cards []*Card) []*FourConsecutivePairs {
	pairs := []*FourConsecutivePairs{}
	dubs := GetDubs(cards)
	if len(dubs) < 4 {
		return pairs
	}
	for i := 0; i < len(dubs) - 3; i++ {
		for j := i + 1; j < len(dubs) - 2; j++ {
			for t := j + 1; t < len(dubs) - 1; t++ {
				for k := t + 1; k < len(dubs); k++ {
					if isFourConsecutivePairs(dubs[i], dubs[j], dubs[t], dubs[k]) {
						pairs = append(pairs, NewFourConsecutivePairs(dubs[i], dubs[j], dubs[t], dubs[k]))
					}
				}
			}
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return compareFourConsecutivePairs(pairs[i], pairs[j]) < 0
	})
	return pairs
}

func getAllCardsWithRank(cards []*Card, rank Rank) []*Card {
	list := []*Card{}
	for i := 0; i < len(cards); i++ {
		if cards[i].rank == rank {
			list = append(list, cards[i])
		}
	}
	return list
}