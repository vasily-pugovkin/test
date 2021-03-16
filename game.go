package tienlen_bot

import (
	"math/rand"
)

type GameConfiguration struct {
	Passed               []bool
	MaxPlayer            int
	PreviousPlayerIndex  int
	CurrentPlayerIndex   int
	LastDealtCombination Combination
	Rape                 bool
	IsFirstTurn          bool
	UseHeuristic         bool
}

func NewDefaultGameConfig(maxPlayers int) *GameConfiguration {
	return &GameConfiguration{
		Passed:               make([]bool, maxPlayers),
		MaxPlayer:            maxPlayers,
		PreviousPlayerIndex:  0,
		CurrentPlayerIndex:   0,
		LastDealtCombination: nil,
		Rape:                 false,
		IsFirstTurn:          true,
		UseHeuristic:         true,
	}
}

type Game interface {
	Move(combination Combination)
	GetCurrentPlayer() Player
	GetCurrentPlayerIndex() int
	GetPreviousPlayerIndex() int
	GetWinnerIndex() int
	GetMaxPlayerNumber() int
	Copy() Game
	AllAvailableCombinations() []Combination
	IsEnd() bool
	PlayRandomUntilEnd()
	GetReward() Reward
	AddPlayer(player Player)
	CurrentNumberOfPlayers() int
	Validate()
	NextTurn()
	GetConfig() *GameConfiguration
	SetConfig(configuration *GameConfiguration)
	GetPlayerAt(index int) Player
	PlayerPassed(index int) bool
	GetLastDealtCombination() Combination
	HasNoLastDealtCombination() bool
}

type LocalGame struct {
	players              []Player
	reward               Reward
	maxNumberOfPlayers   int
	size                 int
	currentPlayerIndex   int
	previousPlayerIndex  int
	lastDealtCombination Combination
	passedPlayersCheck   []bool
	config               *GameConfiguration
	isFirstTurn          bool
	isEnd                bool
	ply                  int
}

func NewGame(config *GameConfiguration) Game {
	return &LocalGame{
		players:              make([]Player, config.MaxPlayer),
		reward:               nil,
		maxNumberOfPlayers:   config.MaxPlayer,
		size:                 0,
		currentPlayerIndex:   config.CurrentPlayerIndex,
		previousPlayerIndex:  config.PreviousPlayerIndex,
		lastDealtCombination: config.LastDealtCombination,
		passedPlayersCheck:   config.Passed,
		config:               config,
		isFirstTurn:          config.IsFirstTurn,
		isEnd:                false,
		ply:                  0,
	}
}

func (l *LocalGame) Move(combination Combination) {
	l.ply++
	if l.isFirstTurn {
		l.isFirstTurn = false
	}
	if combination.Kind() == CombinationPass {
		if l.currentPlayerIndex == l.previousPlayerIndex {
			panic("current player can not pass (must move)")
		}
		l.passedPlayersCheck[l.currentPlayerIndex] = true
		l.NextTurn()
		return
	}
	l.lastDealtCombination = combination
	l.GetCurrentPlayer().Remove(combination)
	l.previousPlayerIndex = l.currentPlayerIndex
	l.NextTurn()
}

func (l *LocalGame) GetCurrentPlayer() Player {
	return l.players[l.currentPlayerIndex]
}

func (l *LocalGame) GetCurrentPlayerIndex() int {
	return l.currentPlayerIndex
}

func (l *LocalGame) GetPreviousPlayerIndex() int {
	return l.previousPlayerIndex
}

func (l *LocalGame) GetWinnerIndex() int {
	for i := 0; i < l.maxNumberOfPlayers; i++ {
		if len(l.players[i].AllAvailableCombinations()) == 0 {
			return i
		}
	}
	return 0
}

func (l *LocalGame) GetMaxPlayerNumber() int {
	return l.maxNumberOfPlayers
}

func (l *LocalGame) Copy() Game {
	game := &LocalGame{
		players:              make([]Player, len(l.players)),
		reward:               nil,
		maxNumberOfPlayers:   l.maxNumberOfPlayers,
		size:                 l.size,
		currentPlayerIndex:   l.currentPlayerIndex,
		previousPlayerIndex:  l.previousPlayerIndex,
		lastDealtCombination: l.lastDealtCombination,
		passedPlayersCheck:   make([]bool, len(l.passedPlayersCheck)),
		config:               l.config,
		isFirstTurn:          l.isFirstTurn,
		isEnd:                l.isEnd,
		ply:                  l.ply,
	}
	for i := 0; i < l.maxNumberOfPlayers; i++ {
		game.players[i] = l.players[i].Copy()
	}
	if notNil(l.reward) {
		game.reward = l.reward.Copy()
	}
	copy(game.passedPlayersCheck, l.passedPlayersCheck)
	return game
}

func (l *LocalGame) AllAvailableCombinations() []Combination {
	player := l.GetCurrentPlayer()
	if l.previousPlayerIndex == l.currentPlayerIndex {
		list := player.AllAvailableCombinations()
		if l.isFirstTurn {
			card := player.GetSmallestCard()
			availableCards := []Combination{}
			for i := 0; i < len(list); i++ {
				if containsCard(list[i].Cards(), card) {
					availableCards = append(availableCards, list[i])
				}
			}
			return availableCards
		}
		return list
	} else {
		return player.AllAvailableCombinationsDefeat(l.lastDealtCombination)
	}
}

func (l *LocalGame) IsEnd() bool {
	return l.isEnd
}

func (l *LocalGame) PlayRandomUntilEnd() {
	for {
		if l.IsEnd() {
			break
		}
		list := l.AllAvailableCombinations();
		if len(list) > 0 &&
			len(list[len(list)-1].Cards()) == l.GetCurrentPlayer().GetCardsLength() {
			l.Move(list[len(list) - 1])
			break
		}
		if len(list) == 0 || l.currentPlayerIndex != l.previousPlayerIndex {
			list = append(list, NewPass())
		}
		combination := list[rand.Intn(len(list))]
		l.Move(combination)
	}
	l.reward = NewReward(l.maxNumberOfPlayers)
	if l.config.Rape {
		winner := l.GetWinnerIndex()
		if l.players[winner].IsBot() {
			for i := 0; i < l.maxNumberOfPlayers; i++ {
				if l.players[i].IsBot() {
					l.reward.SetScore(i, 1)
				} else {
					l.reward.SetScore(i, 0)
				}
			}
		} else {
			for i := 0; i < l.maxNumberOfPlayers; i++ {
				l.reward.SetScore(i, ifThen(len(l.players[i].AllAvailableCombinations()) < 0, float64(1), float64(0)).(float64))
			}
		}
	} else {
		if l.config.UseHeuristic {
			winner := l.GetWinnerIndex()
			total := - float64(l.ply) * FactorPly
			for i := 0; i < l.maxNumberOfPlayers; i++ {
				total += l.players[i].GetScore()
			}
			for i := 0; i < l.maxNumberOfPlayers; i++ {
				if i == winner {
					l.reward.SetScore(i, 1 + total)
				} else {
					l.reward.SetScore(i, - l.GetPlayerAt(i).GetScore())
				}
			}
		} else {
			for i := 0; i < l.maxNumberOfPlayers; i++ {
				l.reward.SetScore(i, ifThen(len(l.GetPlayerAt(i).AllAvailableCombinations()) <= 0, float64(1), float64(0)).(float64))
			}
		}
	}
}

func (l *LocalGame) GetReward() Reward {
	return l.reward
}

func (l *LocalGame) AddPlayer(player Player) {
	for i := 0; i < l.maxNumberOfPlayers; i++ {
		if isNil(l.players[i]) {
			player.SetIndex(i)
			player.Validate()
			l.players[i] = player
			l.size++
			break
		}
	}
	if l.size == l.maxNumberOfPlayers {
		l.Validate()
	}
}

func (l *LocalGame) CurrentNumberOfPlayers() int {
	return l.size
}

func (l *LocalGame) Validate() {
	for _, player := range l.players {
		if isNil(player) {
			panic("nil player")
		}
	}
}

func (l *LocalGame) NextTurn() {
	l.isEnd = len(l.GetCurrentPlayer().AllAvailableCombinations()) == 0
	l.increaseIndex()
	for l.passedPlayersCheck[l.currentPlayerIndex] {
		l.increaseIndex()
	}
	if l.currentPlayerIndex == l.previousPlayerIndex {
		for i := range l.passedPlayersCheck {
			l.passedPlayersCheck[i] = false
		}
	}
}

func (l *LocalGame) GetConfig() *GameConfiguration {
	return l.config
}

func (l *LocalGame) SetConfig(configuration *GameConfiguration) {
	l.config = configuration
}

func (l *LocalGame) GetPlayerAt(index int) Player {
	return l.players[index]
}

func (l *LocalGame) PlayerPassed(index int) bool {
	return l.passedPlayersCheck[index]
}

func (l *LocalGame) GetLastDealtCombination() Combination {
	return l.lastDealtCombination
}

func (l *LocalGame) HasNoLastDealtCombination() bool {
	return isNil(l.lastDealtCombination)
}

func (l *LocalGame) increaseIndex() {
	if l.currentPlayerIndex == l.maxNumberOfPlayers - 1 {
		l.currentPlayerIndex = 0
	} else {
		l.currentPlayerIndex++
	}
}
