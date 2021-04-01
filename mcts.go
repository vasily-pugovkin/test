package tienlen_bot

import (
	"fmt"
	"math"
	"time"
)

type MctsConfig struct {
	Interactions int
	C float64
	Debug bool
	MinThinkingTime int64
	MaxThinkingTime int64
	K float64
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
	pairs := getSmallestPairsInPairsList(game)
	if notNil(pairs) {
		return pairs
	}
	// monte carlo tree search algorithm
	root := NewNode(nil, nil, -1, game)
	root.SetCFactor(config.C)
	root.SetKFactor(config.K)
	startThinkingTime := currentTimeMillis()
	for interactions > 0 && currentTimeMillis() - startThinkingTime < config.MaxThinkingTime {
		interactions--
		/* keep playing while the ratio of winning is less than 50% */
		if currentTimeMillis() - startThinkingTime > config.MinThinkingTime {
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
			config.Interactions - interactions, root.GetReward(), root.GetVisit(), currentTimeMillis() - startThinkingTime))
		root.PrintAllChildren()
	}

	return root.GetMostVisitedChildCombination()
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/*
   KK 33
   2. Nếu đôi K không phải là lớn nhất bài dựa trên các lá đã đánh ra thì xử lí tiếp:
       2.1. Nếu bài user kia còn > 4 lá thì cứ đánh đôi bé trước
       2.2. Nếu user kia còn <= 4 lá thì đánh đôi K trước
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

// nếu bộ còn toàn cóc lẻ và bên kia không chặn được con cóc lẻ nào thì đánh lần lượt từ thấp lên cao
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
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() || game.PlayerPassed(i) {
			continue
		}
		if len(game.GetPlayerAt(i).AllAvailableCombinationsDefeat(list[1])) != 0 {
			return nil
		}
	}
	return list[1]
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