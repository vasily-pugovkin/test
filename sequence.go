package tienlen_bot

type Sequence struct {
	cardList    []*Card
	minRank     Rank
	suit        Suit
	maxRank     Rank
	homogeneity bool
}

func NewSequence(cards []*Card) *Sequence {
	if !isSequence(cards) {
		panic("invalid sequence")
	}
	return &Sequence{
		cardList:    cards,
		minRank:     cards[0].rank,
		suit:        cards[len(cards)-1].suit,
		maxRank:     cards[len(cards)-1].rank,
		homogeneity: isHomogeneitySequence(cards),
	}
}

func (s *Sequence) Kind() CombinationKind {
	return CombinationSequence
}

func (s *Sequence) Equals(combination Combination) bool {
	if combination.Kind() != CombinationSequence {
		return false
	}
	sequence := combination.(*Sequence)
	if len(sequence.cardList) != len(s.cardList) {
		return false
	}
	for i := 0; i < len(s.cardList); i++ {
		if !s.cardList[i].equals(sequence.cardList[i]) {
			return false
		}
	}
	return true
}

func (s *Sequence) Cards() []*Card {
	return s.cardList
}

func (s *Sequence) Defeats(combination Combination) bool {
	if combination.Kind() != CombinationSequence{
		return false
	}
	sequence := combination.(*Sequence)
	if len(s.cardList) != len(sequence.cardList) {
		return false
	}
	if !s.homogeneity && sequence.homogeneity {
		return false
	}
	if s.minRank > sequence.minRank {
		return true
	}
	if s.minRank < sequence.minRank {
		return false
	}
	if s.maxRank >= Six {
		return s.suit > sequence.suit
	}
	return true
}

func (s *Sequence) Copy() Combination {
	return NewSequence(s.cardList)
}

func (s *Sequence) String() string {
	tmp := "["
	for _, card := range s.cardList {
		tmp += card.String() + " "
	}
	tmp += "]"
	return tmp
}