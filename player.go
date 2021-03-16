package tienlen_bot

import (
	"sort"
)

type Player interface {
	GetOriginalCards() []*Card
	AllAvailableCombinations() []Combination
	AllAvailableCombinationsDefeat(combination Combination) []Combination
	Remove(combination Combination)
	SetCards(cards []*Card)
	WithCards(cards []*Card) Player
	SetIndex(index int)
	Validate()
	Copy() Player
	IsBot() bool
	SetBot(bot bool)
	GetCardsLength() int
	GetSmallestCard() *Card
	GetScore() float64
	GetAllCombinationsHasSameAtLeastOneCardWith(combination Combination) []Combination
}

type LocalPlayer struct {
	index        int
	isBot        bool
	cards        []*Card
	combinations []Combination
	connectors   map[Combination][]Combination
	cardsLength  int
	score  float64
}

func NewPlayer() Player {
	return &LocalPlayer{
		index:        -1,
		isBot:        false,
		cards:        nil,
		combinations: make([]Combination, 0),
		connectors:   make(map[Combination][]Combination),
		cardsLength:  13,
	}
}

func (l *LocalPlayer) GetOriginalCards() []*Card {
	return l.cards
}

func (l *LocalPlayer) AllAvailableCombinations() []Combination {
	return l.combinations
}

func (l *LocalPlayer) AllAvailableCombinationsDefeat(combination Combination) []Combination {
	combinations := []Combination{}
	for i := 0; i < len(l.combinations); i++ {
		if l.combinations[i].Defeats(combination) {
			combinations = append(combinations, l.combinations[i])
		}
	}
	return combinations
}

func (l *LocalPlayer) Remove(combination Combination) {
	connector, ok := l.connectors[combination]
	if !ok {
		panic("invalid input")
	}
	l.removeCombination(combination)
	l.cardsLength -= len(combination.Cards())
	if l.cardsLength < 0 {
		panic("length of card can not be less than zero")
	}
	for _, c := range connector {
		index := l.removeCombination(c)
		if index >= 0 {
			l.computeScore(c)
		}
	}
	l.computeScore(combination)
	if l.score < -0.00001 {
		panic("score must be bigger than zero")
	}
}

func (l *LocalPlayer) SetCards(cards []*Card) {
	sort.Slice(cards, func(i, j int) bool {
		return compareCard(cards[i], cards[j]) < 0
	})
	l.cards = cards
	l.cardsLength = len(cards)
}

func (l *LocalPlayer) WithCards(cards []*Card) Player {
	l.SetCards(cards)
	return l
}

func (l *LocalPlayer) SetIndex(index int) {
	l.index = index
}

func (l *LocalPlayer) Validate() {
	if l.cards == nil {
		panic("invalid Cards")
	}
	for _, card := range l.cards {
		l.combinations = append(l.combinations, NewSingleCard(card))
	}
	dubs := GetDubs(l.cards)
	for i := 0; i < len(dubs); i++ {
		l.combinations = append(l.combinations, dubs[i])
	}
	sequences := GetSequence(l.cards)
	for i := 0; i < len(sequences); i++ {
		l.combinations = append(l.combinations, sequences[i])
	}
	twoPairs := GetTwoConsecutivePairs(l.cards)
	for i := 0; i < len(twoPairs); i++ {
		l.combinations = append(l.combinations, twoPairs[i])
	}
	trips := GetTrips(l.cards)
	for i := 0; i < len(trips); i++ {
		l.combinations = append(l.combinations, trips[i])
	}
	threePairs := GetThreeConsecutivePairs(l.cards)
	for i := 0; i < len(threePairs); i++ {
		l.combinations = append(l.combinations, threePairs[i])
	}
	quads := GetQuads(l.cards)
	for i := 0; i < len(quads); i++ {
		l.combinations = append(l.combinations, quads[i])
	}
	fourPairs := GetFourConsecutivePairs(l.cards)
	for i := 0; i < len(fourPairs); i++ {
		l.combinations = append(l.combinations, fourPairs[i])
	}
	for i := 0; i < len(l.combinations); i++ {
		combination := l.combinations[i]
		connector := []Combination{}
		for j := 0; j < len(l.combinations); j++ {
			if i != j && hasAtLeastSameOneCard(combination.Cards(), l.combinations[j].Cards()) {
				connector = append(connector, l.combinations[j])
			}
		}
		l.connectors[combination] = connector
		l.computeScore(l.combinations[i])
	}
	l.score = -l.score
	sort.Slice(l.combinations, func(i, j int) bool {
		return len(l.combinations[i].Cards()) < len(l.combinations[j].Cards())
	})
}

func (l *LocalPlayer) Copy() Player {
	player :=  &LocalPlayer{
		index:        l.index,
		isBot:        l.isBot,
		cards:        l.cards,
		combinations: make([]Combination, len(l.combinations)),
		connectors:   l.connectors,
		cardsLength:  l.cardsLength,
		score:        l.score,
	}
	copy(player.combinations, l.combinations)
	return player
}

func (l *LocalPlayer) IsBot() bool {
	return l.isBot
}

func (l *LocalPlayer) SetBot(bot bool) {
	l.isBot =  bot
}

func (l *LocalPlayer) GetCardsLength() int {
	return l.cardsLength
}

func (l *LocalPlayer) removeCombination(combination Combination) int {
	for i := 0; i < len(l.combinations); i++ {
		if l.combinations[i].Equals(combination) {
			l.combinations[i] = l.combinations[len(l.combinations) - 1]
			l.combinations = l.combinations[:len(l.combinations) - 1]
			return i
		}
	}
	return -1
}

func (l *LocalPlayer) GetSmallestCard() *Card {
	return l.cards[0]
}

func hasAtLeastSameOneCard(cards1, cards2 []*Card) bool {
	for _, c1 := range cards1 {
		for _, c2 := range cards2 {
			if c1.equals(c2) {
				return true
			}
		}
	}
	return false
}

func (l *LocalPlayer) GetScore() float64 {
	return l.score
}

func (l *LocalPlayer) GetAllCombinationsHasSameAtLeastOneCardWith(combination Combination) []Combination {
	return l.connectors[combination]
}

func (l *LocalPlayer) computeScore(combination Combination) {
	switch combination.Kind() {
	case CombinationSingle:
		c := combination.(*SingleCard)
		if c.card.rank == Two {
			if c.card.suit < Diamond {
				l.score -= FactorBlack2sCard
			} else {
				l.score -= FactorRed2sCard
			}
		} else {
			l.score -= FactorNormalSingleCard
		}
	case CombinationThreeConsecutivePairs:
		l.score -= FactorThreePairs
	case CombinationQuads:
		l.score -= FactorQuads
	case CombinationFourConsecutivePairs:
		l.score -= FactorFourPairs
	}
}