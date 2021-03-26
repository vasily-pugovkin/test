package tienlen_bot

import (
	"fmt"
	"math"
	"math/rand"
)

type Node interface {
	Select(game Game) Node
	Expand(game Game) Node
	Simulate(game Game) Reward
	BackPropagation(reward Reward)
	GetMostVisitedChildCombination() Combination
	GetUCT() float64
	GetVisit() int
	GetCombination() Combination
	PrintAsTree(space string)
	PrintAllChildren()
	SetCFactor(c float64)
	GetCFactor() float64
	GetReward() Reward
	SetKFactor(k float64)
}

type LocalNode struct {
	parent                 Node
	combination            Combination
	children               []Node
	unexploredCombinations []Combination
	reward                 Reward
	visit                  int
	currentPlayerIndex     int
	C                      float64
	K                      float64
}

func NewNode(parent *LocalNode, combination Combination, playerIndex int, game Game) Node {
	node := &LocalNode{
		parent:                 parent,
		combination:            combination,
		children:               []Node{},
		unexploredCombinations: []Combination{},
		reward:                 NewReward(game.GetMaxPlayerNumber()),
		visit:                  0,
		currentPlayerIndex:     playerIndex,
		C:                      math.Sqrt(2),
		K:                      0,
	}
	if notNil(parent) {
		node.SetCFactor(parent.GetCFactor())
		node.K = parent.K
	}
	list := game.AllAvailableCombinations()
	if !game.IsEnd() {
		if len(list) > 0 && len(list[len(list)-1].Cards()) == game.GetCurrentPlayer().GetCardsLength() {
			node.unexploredCombinations = append(node.unexploredCombinations, list[len(list)-1])
		} else {
			node.unexploredCombinations = make([]Combination, len(list))
			copy(node.unexploredCombinations, list)
			if isNil(parent) && game.GetCurrentPlayerIndex() == game.GetPreviousPlayerIndex() && !game.GetConfig().IsFirstTurn {
				node.removeStrongCombinationsIfNotNecessary(game)
				println(len(node.unexploredCombinations))
			} else if game.GetCurrentPlayerIndex() != game.GetPreviousPlayerIndex() {
				node.unexploredCombinations = append(node.unexploredCombinations, NewPass())
			}
			if isNil(parent) {
				node.removeStrongCombinationsThan2IfHave2()
				node.keepConsecutivePairsForDefeating2(game)
				node.remove2IfIsFirstTurn(game)
				if node.canDefeatTheirSingleCard(game) {
					node.removePass()
				}
			}
		}
	}

	return node
}

func (l *LocalNode) Select(game Game) Node {
	if game.IsEnd() || len(l.unexploredCombinations) > 0 {
		return l
	}
	var maxScore float64 = -100000000
	var selectedNode Node = l
	for _, child := range l.children {
		uct := child.GetUCT()
		if uct > maxScore {
			maxScore = uct
			selectedNode = child
		}
	}
	if notNil(selectedNode.GetCombination()) {
		game.Move(selectedNode.GetCombination())
	}

	return selectedNode.Select(game)
}

func (l *LocalNode) Expand(game Game) Node {
	if len(l.unexploredCombinations) <= 0 {
		return l
	}
	randomNumber := rand.Intn(len(l.unexploredCombinations))
	combination := l.removeUnexploredCombinationAt(randomNumber)
	player := game.GetCurrentPlayerIndex()
	game.Move(combination)
	node := NewNode(l, combination, player, game)
	l.children = append(l.children, node)
	return node
}

func (l *LocalNode) Simulate(game Game) Reward {
	game.PlayRandomUntilEnd()
	return game.GetReward()
}

func (l *LocalNode) BackPropagation(reward Reward) {
	l.reward.AddReward(reward)
	l.visit++
	if notNil(l.parent) {
		l.parent.BackPropagation(reward)
	}
}

func (l *LocalNode) GetMostVisitedChildCombination() Combination {
	mostVisitCount := 0
	var mostVisitedNode Node = nil
	for i := range l.children {
		if l.children[i].GetVisit() > mostVisitCount {
			mostVisitCount = l.children[i].GetVisit()
			mostVisitedNode = l.children[i]
		}
	}
	if isNil(mostVisitedNode) {
		return nil
	}
	return mostVisitedNode.GetCombination()
}

func (l *LocalNode) GetUCT() float64 {
	exploit := l.reward.GetScoreOfPlayer(l.currentPlayerIndex) / float64(l.visit)
	discover := l.C * math.Sqrt(math.Log(float64(l.parent.GetVisit()))/float64(l.visit))
	balance := l.K / (l.K + float64(l.visit))
	return exploit + discover + balance
}

func (l *LocalNode) GetVisit() int {
	return l.visit
}

func (l *LocalNode) GetCombination() Combination {
	return l.combination
}

func (l *LocalNode) PrintAsTree(space string) {
	for _, node := range l.children {
		node.PrintAsTree(space + "    |")
	}
}

func (l *LocalNode) PrintAllChildren() {
	s := ""
	for _, node := range l.children {
		s += fmt.Sprintf("Node:   %-40s", node.GetCombination())
		s += "|"
		s += fmt.Sprintf("Visit:  %-20d", node.GetVisit())
		s += "|"
		s += fmt.Sprintf("Reward: %-30s", node.GetReward())
		s += "\n"
	}
	println(s)
}

func (l *LocalNode) SetCFactor(c float64) {
	l.C = c
}

func (l *LocalNode) GetCFactor() float64 {
	return l.C
}

func (l *LocalNode) GetReward() Reward {
	return l.reward
}

func (l *LocalNode) SetKFactor(k float64) {
	l.K = k
}

//  xóa 3 đôi thông, 4 đôi thông và 2 đi
func (l *LocalNode) removeStrongCombinationsIfNotNecessary(game Game) {
	conf := game.GetConfig()

	for i := 0; i < conf.MaxPlayer; i++ {
		if game.GetPlayerAt(i).GetCardsLength() <= NumberOfCardsAtLateGame {
			return
		}
	}
	player := game.GetCurrentPlayer()
	temp := make([]Combination, len(l.unexploredCombinations))
	copy(temp, l.unexploredCombinations)

	removedList := []Combination{}
	for i := range l.unexploredCombinations {
		if (l.unexploredCombinations[i].Kind() == CombinationThreeConsecutivePairs &&
			player.GetCardsLength() > 7) ||
			(l.unexploredCombinations[i].Kind() == CombinationFourConsecutivePairs &&
				player.GetCardsLength() > 9) ||
			l.unexploredCombinations[i].Kind() == CombinationQuads ||
			containsRank(l.unexploredCombinations[i].Cards(), Two) {
			removedList = append(removedList, l.unexploredCombinations[i])
		}
	}

	for i := range removedList {
		l.removeUnexploredCombination(removedList[i])
		connectors := game.GetCurrentPlayer().GetAllCombinationsHasSameAtLeastOneCardWith(removedList[i])
		for j := range connectors {
			l.removeUnexploredCombination(connectors[j])
		}
	}

	// nếu vô tình xóa hết con mẹ nó nước đi thì thôi coi như xí xóa
	if len(l.unexploredCombinations) == 0 {
		l.unexploredCombinations = temp
	}
}

// ưu tiên đánh 2 trước khi ra tứ quý, 3 đôi thông hoặc 4 đôi thông
func (l *LocalNode) removeStrongCombinationsThan2IfHave2() {
	removedList1 := []Combination{}
	removedList2 := []Combination{}
	contains2sCard := false
	for i := range l.unexploredCombinations {
		o := l.unexploredCombinations[i]
		if !contains2sCard && o.Kind() == CombinationSingle {
			if o.(*SingleCard).card.rank == Two {
				contains2sCard = true
			}
		}
		if o.Kind() == CombinationQuads ||
			o.Kind() == CombinationFourConsecutivePairs ||
			o.Kind() == CombinationThreeConsecutivePairs {
			removedList1 = append(removedList1, o)
		}

		if containsRank(o.Cards(), Two) && o.Kind() != CombinationSingle {
			removedList2 = append(removedList2, o)
		}
	}

	if len(removedList2) < len(l.unexploredCombinations)-1 {
		for i := range removedList2 {
			l.removeUnexploredCombination(removedList2[i])
		}
	}

	if !contains2sCard {
		return
	}

	for i := range removedList1 {
		l.removeUnexploredCombination(removedList1[i])
	}
}

// luôn dùng tứ quý, 3 đôi thông hoặc 4 đôi thông nếu người trước đánh 2
func (l *LocalNode) keepConsecutivePairsForDefeating2(game Game) {
	if game.HasNoLastDealtCombination() {
		return
	}
	if !containsRank(game.GetLastDealtCombination().Cards(), Two) {
		return
	}
	// nếu bot đánh đôi 2, hoặc tam 2 mà chặn được thì chặn luôn
	if len(game.GetLastDealtCombination().Cards()) >= 2 &&
		len(l.unexploredCombinations) >= 2 {
		l.removePass()
		return
	}

	//nếu bot đánh 1 con 2 lẻ
	if game.GetLastDealtCombination().Kind() == CombinationSingle {
		if !l.hasStrongCombination(l.unexploredCombinations) {
			return
		}
		combinations := game.GetPlayerAt(game.GetPreviousPlayerIndex()).AllAvailableCombinations()
		contains2 := false
		for i := range combinations {
			if combinations[i].Kind() == CombinationSingle &&
				combinations[i].(*SingleCard).card.rank == Two {
				contains2 = true
				break
			}
		}

		// nếu người chơi trước không còn 2 thì chặn luôn
		if !contains2 {
			l.removePass()
			return
		}

		// nếu người chơi trước có 2 thì 100% chặn nếu con vừa đánh là 2 đỏ và 90% chặn nếu là 2 đen
		card := game.GetLastDealtCombination().Cards()[0]
		if card.suit == Heart || card.suit == Diamond {
			l.removePass()
		} else {
			if rand.Intn(10) > 1 {
				l.removePass()
			}
		}
	}
}

// loại con 2 ra nếu turn này mình không phải chặn ai
func (l *LocalNode) remove2IfIsFirstTurn(game Game) {
	if !game.HasNoLastDealtCombination() {
		return
	}
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		if len(game.GetPlayerAt(i).AllAvailableCombinations()) == 1 {
			return
		}
	}
	card := l.getAlmostSmallestSingleCardWhichNotInAnyCombination()
	if isNil(card) || card.rank == Two {
		return
	}
	removedList := []Combination{}
	for i := range l.unexploredCombinations {
		if containsRank(l.unexploredCombinations[i].Cards(), Two) {
			removedList = append(removedList, l.unexploredCombinations[i])
		}
	}

	for i := range removedList {
		l.removeUnexploredCombination(removedList[i])
	}
}

// bắt buộc chặn con lẻ nếu có thể
func (l *LocalNode) canDefeatTheirSingleCard(game Game) bool {
	list := game.GetCurrentPlayer().AllAvailableCombinations()
	if !game.HasNoLastDealtCombination() && game.GetLastDealtCombination().Kind() == CombinationSingle {
	Loop:
		for i := range l.unexploredCombinations {
			if l.unexploredCombinations[i].Kind() != CombinationSingle {
				continue
			}
			for j := range list {
				if list[j].Kind() != CombinationSingle && containsCard(list[j].Cards(), l.unexploredCombinations[i].(*SingleCard).card) {
					continue Loop
				}
			}
			return true
		}
	}
	return false
}

// xóa pass tức là bắt buộc đánh
func (l *LocalNode) removePass() {
	for i := 0; i < len(l.unexploredCombinations); i++ {
		if l.unexploredCombinations[i].Kind() == CombinationPass {
			l.removeUnexploredCombination(l.unexploredCombinations[i])
			break
		}
	}
}

// xóa bộ khỏi list
func (l *LocalNode) removeUnexploredCombination(combination Combination) {
	for i := range l.unexploredCombinations {
		if l.unexploredCombinations[i].Equals(combination) {
			l.unexploredCombinations = append(l.unexploredCombinations[:i], l.unexploredCombinations[i+1:]...)
			return
		}
	}
}

func (l *LocalNode) removeUnexploredCombinationAt(index int) Combination {
	c := l.unexploredCombinations[index]
	l.removeUnexploredCombination(l.unexploredCombinations[index])
	return c
}

// có tứ quý, 3 đôi thông, 4 đôi thông
func (l *LocalNode) hasStrongCombination(combinations []Combination) bool {
	for i := range combinations {
		if combinations[i].Kind() == CombinationQuads || combinations[i].Kind() == CombinationThreeConsecutivePairs ||
			combinations[i].Kind() == CombinationFourConsecutivePairs {
			return true
		}
	}
	return false
}

// lấy lá lẻ có giá trị nhỏ thứ 2 mà ko có trong bất kì bộ nào
func (l *LocalNode) getAlmostSmallestSingleCardWhichNotInAnyCombination() *Card {
	var card *Card
	count := 0
Loop:
	for i := range l.unexploredCombinations {
		if l.unexploredCombinations[i].Kind() != CombinationSingle {
			continue
		}
		for j := range l.unexploredCombinations {
			if l.unexploredCombinations[j].Kind() == CombinationSingle {
				continue
			}
			if containsCard(l.unexploredCombinations[j].Cards(), l.unexploredCombinations[i].(*SingleCard).card) {
				continue Loop
			}
		}
		count++
		if count == 2 {
			card = l.unexploredCombinations[i].(*SingleCard).card
			break
		}
	}
	return card
}
