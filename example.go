package tienlen_bot

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func createGame(config *GameConfiguration) Game {
	game := NewGame(config)
	deck := NewDeck()
	for i := 0; i < config.MaxPlayer; i++ {
		player := NewPlayer()
		player.SetBot(false)
		player.SetCards(deck.randomCards(13))
		game.AddPlayer(player)
	}
	return game
}

func getCards(player Player) []*Card {
	cards := []*Card{}
	list := player.AllAvailableCombinations()
	for i := range list {
		tmpCards := list[i].cards()
		for j := range tmpCards {
			if containsCard(cards, tmpCards[j]) {
				continue
			}
			cards = append(cards, tmpCards[j])
		}
	}
	return cards
}

func StartNewExampleGame() {
	gameConfig := NewDefaultGameConfig(4)
	game := createGame(gameConfig)
	mctsConfig := NewDefaultMctsConfig()
	reader := bufio.NewReader(os.Stdin)

	for !game.IsEnd() {
		println(fmt.Sprintf("Cards: %+v", getCards(game.GetCurrentPlayer())))
		if game.GetCurrentPlayerIndex() == 0 {
			println("Your turn")
			combinations := game.AllAvailableCombinations()
			if game.GetCurrentPlayerIndex() != game.GetPreviousPlayerIndex() {
				combinations = append(combinations, NewPass())
			}
			for i := range combinations {
				println(fmt.Sprintf("%d. %s", i+1, combinations[i]))
			}
			println("your selection:")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			input, err := strconv.ParseInt(strings.ReplaceAll(text, "\n", ""), 10, 64)
			if err != nil {
				panic(err)
			}
			println(fmt.Sprintf("[PLAYER] %s", combinations[input - 1]))
			game.Move(combinations[input - 1])
		} else {
			bestCombination := SelectBestCombination(game, mctsConfig)
			println(fmt.Sprintf("Bot %d dropped %s", game.GetCurrentPlayerIndex(), bestCombination))
			game.Move(bestCombination)
		}
	}
	println(fmt.Sprintf("****** %s IS WINNER ******", ifThen(game.GetWinnerIndex() == 0, "PLAYER", "BOT")))
}