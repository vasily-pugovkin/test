package tienlen_bot

import (
	"fmt"
	"math"
	"time"
)

type MctsConfig struct {
	Interactions    int
	C               float64
	Debug           bool
	MinThinkingTime int64
	MaxThinkingTime int64
	K               float64
}

func NewDefaultMctsConfig() *MctsConfig {
	return &MctsConfig{
		Interactions:    1000000000,
		C:               math.Sqrt(2),
		Debug:           false,
		MinThinkingTime: 1000,
		MaxThinkingTime: 2000,
		K:               500,
	}
}

func SelectBestCombination(game Game, config *MctsConfig) Combination {
	interactions := config.Interactions
	list := game.AllAvailableCombinations()
	if len(list) == 0 {
		return NewPass()
	}
	// person knowledge to make bot looks similar to a real person
	singleCard := getBestMoveForDefeatingSingleCard(game)
	if notNil(singleCard) {
		return singleCard
	}
	singleCard = getBestMoveIfAllOtherPeopleHasOnlyOneCard(game)
	if notNil(singleCard) {
		return singleCard
	}
	pairs := getSmallestPairsInPairsList(game)
	if notNil(pairs) {
		return pairs
	}
	// monte carlo tree search algorithm
	root := NewNode(nil, nil, -1, game)
	if len(root.(*LocalNode).unexploredCombinations) == 1 {
		return root.(*LocalNode).unexploredCombinations[0]
	}
	root.SetCFactor(config.C)
	root.SetKFactor(config.K)
	startThinkingTime := currentTimeMillis()
	for interactions > 0 && currentTimeMillis()-startThinkingTime < config.MaxThinkingTime {
		interactions--
		/* keep playing while the ratio of winning is less than 50% */
		if currentTimeMillis()-startThinkingTime > config.MinThinkingTime {
			var x, y float64
			for i := 0; i < game.GetMaxPlayerNumber(); i++ {
				if i == game.GetCurrentPlayerIndex() {
					x = root.GetReward().GetScoreOfPlayer(i)
				} else {
					y = math.Max(y, root.GetReward().GetScoreOfPlayer(i))
				}
			}
			if x > y {
				break
			}
		}
		/* continue loop */
		gameCopy := game.Copy()
		node := root.Select(gameCopy)
		node = node.Expand(gameCopy)
		reward := node.Simulate(gameCopy)
		node.BackPropagation(reward)
	}

	if config.Debug {
		println(fmt.Sprintf("MCTS %d interactions, reward: %+v, visit: %d, thinking time: %d",
			config.Interactions-interactions, root.GetReward(), root.GetVisit(), currentTimeMillis()-startThinkingTime))
		root.PrintAllChildren()
	}

	return root.GetMostVisitedChildCombination()
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/*
   KK 33
   2. N???u ????i K kh??ng ph???i l?? l???n nh???t b??i d???a tr??n c??c l?? ???? ????nh ra th?? x??? l?? ti???p:
       2.1. N???u b??i user kia c??n > 4 l?? th?? c??? ????nh ????i b?? tr?????c
       2.2. N???u user kia c??n <= 4 l?? th?? ????nh ????i K tr?????c
*/
func getSmallestPairsInPairsList(game Game) Combination {
	list := game.AllAvailableCombinations()
	if len(list) != 6 {
		return nil
	}
	cardCount := 0
	dubsCount := 0
	dubs := make([]Combination, 2)
	for i := range list {
		if list[i].Kind() == CombinationSingle {
			cardCount++
		} else if list[i].Kind() == CombinationDubs {
			if dubsCount >= 2 {
				return nil
			}
			dubs[dubsCount] = list[i]
			dubsCount++
		}
	}
	if cardCount != 4 || dubsCount != 2 {
		return nil
	}
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		combinations := game.GetPlayerAt(i).AllAvailableCombinations()
		count := 0
		for j := range combinations {
			if combinations[j].Kind() == CombinationSingle {
				count++
			}
			if combinations[j].Defeats(dubs[1]) {
				return nil
			}
		}
		if count > 4 {
			return dubs[0]
		}
	}
	return dubs[1]
}

// n???u b??? c??n to??n c??c l??? v?? b??n kia kh??ng ch???n ???????c con c??c l??? n??o th?? ????nh l???n l?????t t??? th???p g???n nh???t l??n cao
// ch??? t??nh n???u c?? ??t nh???t 1 ng?????i c??n 2 l?? tr??? l??n
// c??n n???u t???t c??? m???i ng?????i ?????u c??n 1 l?? th?? ????nh t??? cao xu???ng th???t
func getBestMoveForDefeatingSingleCard(game Game) Combination {
	list := game.AllAvailableCombinations()
	if !isAllSingleCard(list) {
		return nil
	}

	if len(game.GetCurrentPlayer().AllAvailableCombinations()) != len(list) {
		return nil
	}
	if len(list) == 1 {
		return list[0]
	}
	// ki???m tra ??i???u ki???n c?? ai ch???t ???????c qu??n nh??? g???n nh???t kh??ng
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() || game.PlayerPassed(i) {
			continue
		}
		if len(game.GetPlayerAt(i).AllAvailableCombinationsDefeat(list[1])) != 0 {
			return nil
		}
	}
	// check tr?????ng h???p c?? ??t nh???t 1 th???ng tr??n b??n c??n nhi???u h??n 1 l?? th?? ????nh t??? l?? nh??? g???n nh???t n???u c??n nhi???u
	// h??n 2 l??
	// ho???c ????nh l?? nh??? nh???t n???u c??n 2 l??
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		if game.GetPlayerAt(i).GetCardsLength() > 1 {
			// n???u ??t h??n 3 l?? th?? ????nh l?? nh??? nh???t
			if game.GetPlayerAt(game.GetCurrentPlayerIndex()).GetCardsLength() <= 2 {
				return list[0]
			}
			// n???u nhi???u h??n 2 l?? th?? ????nh l?? g???n nh??? nh???t
			return list[1]
		}
	}
	// n???u m???i ng?????i ch??? c??n 1 l?? th?? ????nh l?? to nh???t tr??? xu???ng
	return list[len(list)-1]
}

// n???u m???i ng?????i ch??? c??n 1 l?? v?? m??nh kh??ng c??n b??? n??o (ch??? c??n qu??n l???)
// th?? ????nh t??? to t???i nh??? trong tr?????ng h???p l?? turn t??? ????nh ko ph???i turn ch???t
func getBestMoveIfAllOtherPeopleHasOnlyOneCard(game Game) Combination {
	if game.GetCurrentPlayerIndex() != game.GetPreviousPlayerIndex() {
		return nil
	}
	if !allOtherPeopleHasOnlySingleCardLeft(game) {
		return nil
	}

	list := game.GetCurrentPlayer().AllAvailableCombinations()
	return list[len(list)-1]
}

// ki???m tra t???t c??? ng?????i kh??c c?? c??n 1 l?? ko
// v?? m??nh ch??? c??n to??n con l??? kh??ng
func allOtherPeopleHasOnlySingleCardLeft(game Game) bool {
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			if !isAllSingleCard(game.GetPlayerAt(i).AllAvailableCombinations()) {
				return false
			}
			continue
		}
		if game.GetPlayerAt(i).GetCardsLength() != 1 {
			return false
		}
	}
	return true
}

// ????nh qu??n nh??? nh???t, n???u b??i m??nh ch???c th???ng v?? to??n qu??n l???
// n???u ?????i ph????ng kh??ng c??n qu??n
func keepSmallestCardIfSureWinAndHaveNoCombination() {

}

func isAllSingleCard(combinations []Combination) bool {
	for _, combination := range combinations {
		if combination.Kind() == CombinationSingle {
			continue
		}
		return false
	}
	return true
}
