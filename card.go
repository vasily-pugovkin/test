package tienlen_bot

import (
	"github.com/fatih/color"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Suit int

const (
	Spade Suit = iota
	Club
	Diamond
	Heart
)

func MaxSuit(suits... Suit) Suit {
	suit := Spade
	for i := 0; i < len(suits); i++ {
		if suits[i] > suit {
			suit = suits[i]
		}
	}
	return suit
}

func MinSuit(suits... Suit) Suit {
	suit := Heart
	for i := 0; i < len(suits); i++ {
		if suits[i] < suit {
			suit = suits[i]
		}
	}
	return suit
}

type Rank int

const (
	Three Rank = iota
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
	Two
)

type Card struct {
	rank Rank
	suit Suit
}

func NewCard(rank Rank, suit Suit) *Card {
	return &Card{
		rank: rank,
		suit: suit,
	}
}

func (c *Card) Rank() Rank {
	return c.rank
}

func (c *Card) Suit() Suit {
	return c.suit
}

func (c *Card) String() string {
	s := ""
	switch c.rank {
	case Two:
		s += "2"
	case Three:
		s += "3"
	case Four:
		s += "4"
	case Five:
		s += "5"
	case Six:
		s += "6"
	case Seven:
		s += "7"
	case Eight:
		s += "8"
	case Nine:
		s += "9"
	case Ten:
		s += "10"
	case Jack:
		s += "J"
	case Queen:
		s += "Q"
	case King:
		s += "K"
	case Ace:
		s += "A"
	default:
		s += "undefined"
	}
	blackBold := color.New(color.FgBlack, color.Bold)
	redBold := color.New(color.FgRed, color.Bold)
	switch c.suit {
	case Spade:
		s += "♠"
		s = blackBold.Sprintf(s)
	case Club:
		s += "♣"
		s = blackBold.Sprintf(s)
	case Diamond:
		s += "♦"
		s = redBold.Sprintf(s)
	case Heart:
		s += "♥"
		s = redBold.Sprintf(s)
	default:
		s += "💀"
	}
	return s
}

func (c *Card) equals(card *Card) bool {
	return c.rank == card.rank && c.suit == card.suit
}

func parseCard(s string) *Card {
	c := &Card{}
	var rank, suit string
	if s[:2] == "10" {
		rank = "10"
		suit = s[2:]
	} else {
		rank = s[:1]
		suit = s[1:]
	}
	rank = strings.ToLower(rank)

	switch suit {
	case "♥":
		c.suit = Heart
	case "♦":
		c.suit = Diamond
	case "♣":
		c.suit = Club
	case "♠":
		c.suit = Spade
	default:
		panic("invalid card string " + s)
	}
	switch rank {
	case "a":
		c.rank = Ace
	case "2":
		c.rank = Two
	case "3":
		c.rank = Three
	case "4":
		c.rank = Four
	case "5":
		c.rank = Five
	case "6":
		c.rank = Six
	case "7":
		c.rank = Seven
	case "8":
		c.rank = Eight
	case "9":
		c.rank = Nine
	case "10":
		c.rank = Ten
	case "j":
		c.rank = Jack
	case "q":
		c.rank = Queen
	case "k":
		c.rank = King
	default:
		panic("invalid card string " + s)
	}
	return c
}

func parseCards(s string) []*Card {
	list := strings.Split(strings.ReplaceAll(s, " ", ""), ",")
	cards := make([]*Card, len(list))
	for i := 0; i < len(cards); i++ {
		cards[i] = parseCard(list[i])
	}
	return cards
}

func SortCard(cards []*Card) []*Card {
	sort.Slice(cards, func(i, j int) bool {
		return compareCard(cards[i], cards[j]) < 0
	})
	return cards
}

func ifThen (condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

type Deck struct {
	cards []*Card
}

func NewDeck() *Deck {
	d := &Deck{
		cards: make([]*Card, 52),
	}
	for rank := Three; rank <= Two; rank++ {
		for suit := Spade; suit <= Heart; suit++ {
			d.cards[int(rank) * 4 + int(suit)] = &Card{
				rank: rank,
				suit: suit,
			}
		}
	}
	SortCard(d.cards)
	return d
}

func (d *Deck) randomCards(numberOfCards int) []*Card {
	if len(d.cards) < numberOfCards {
		panic("invalid number of card")
	}
	cards := []*Card{}
	for i := 0; i < numberOfCards; i++ {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(d.cards))
		cards = append(cards, d.cards[index])
		d.cards = append(d.cards[:index], d.cards[index+1:]...)
	}
	return cards
}