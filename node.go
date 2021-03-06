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
	String() string
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
			} else if game.GetCurrentPlayerIndex() != game.GetPreviousPlayerIndex() {
				node.unexploredCombinations = append(node.unexploredCombinations, NewPass())
			}
			if isNil(parent) {
				if !game.HasNoLastDealtCombination() {
					node.keepConsecutivePairsForDefeating2(game)
				}
				if node.allHasOneCardLeft(game) {
					if game.GetCurrentPlayerIndex() == game.GetPreviousPlayerIndex() {
						node.removeSingleCardIfTheyAllHaveOneCardLeft(game)
					}
				} else {
					if game.HasNoLastDealtCombination() {
						// ch??? ??u ti??n ????nh 2 h??n th??ng v?? t??? qu?? n???u kh??ng ph???i turn ch???n
						node.removeStrongCombinationsThan2IfHave2()
						// x??a 2 n???u l?? l?????t ????nh kh??ng ch???t
						// v?? ch??? c??n 1 con 2 v???i c?? ??t nh???t 1 con l??? nh??? h??n 2
						// v?? kh??ng ai c??n 1 l?? tr??n b??n
						node.remove2IfIsFirstTurn(game)
						// x??a 3, 4 ????i th??ng ho???c t??? qu?? n???u ng?????i ch??i kia c??n 2 con 1 con 2 v?? 1 con l???
						// con l??? nh??? h??n ??t nh???t 1 con l??? c???a m??nh
						node.removeCombinationStrongerThan2IfTheyHave2AndOneSmallSingleCardLeft(game)
					}
				}
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
		s += fmt.Sprintf("%-40s", "Node:   " + node.GetCombination().String())
		s += "|"
		s += fmt.Sprintf("%-20s", fmt.Sprintf("Visit:  %d", node.GetVisit()))
		s += "|"
		s += fmt.Sprintf("%-30s", fmt.Sprintf("Reward: %+v", node.GetReward()))
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

//  x??a 3 ????i th??ng, 4 ????i th??ng v?? 2 ??i
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

	// n???u v?? t??nh x??a h???t con m??? n?? n?????c ??i th?? th??i coi nh?? x?? x??a
	if len(l.unexploredCombinations) == 0 {
		l.unexploredCombinations = temp
	}
}

// ??u ti??n ????nh 2 tr?????c khi ra t??? qu??, 3 ????i th??ng ho???c 4 ????i th??ng
func (l *LocalNode) removeStrongCombinationsThan2IfHave2() {
	removedList1 := []Combination{}
	//removedList2 := []Combination{}
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

		//if containsRank(o.Cards(), Two) && o.Kind() != CombinationSingle {
		//	removedList2 = append(removedList2, o)
		//}
	}

	//if len(removedList2) < len(l.unexploredCombinations)-1 {
	//	for i := range removedList2 {
	//		l.removeUnexploredCombination(removedList2[i])
	//	}
	//}

	if !contains2sCard {
		return
	}

	for i := range removedList1 {
		l.removeUnexploredCombination(removedList1[i])
	}
}

// lu??n d??ng t??? qu??, 3 ????i th??ng ho???c 4 ????i th??ng n???u ng?????i tr?????c ????nh 2
func (l *LocalNode) keepConsecutivePairsForDefeating2(game Game) {
	// n???u con ????nh ko ph???i 2 ho???c ????i 2 ho???c tam 2 th?? th??i
	if !containsRank(game.GetLastDealtCombination().Cards(), Two) {
		return
	}
	// n???u bot ????nh ????i 2, ho???c tam 2 m?? ch???n ???????c th?? ch???n lu??n
	if len(game.GetLastDealtCombination().Cards()) >= 2 &&
		len(l.unexploredCombinations) >= 2 {
		l.removePass()
		return
	}

	//n???u bot ????nh 1 con 2 l???
	if game.GetLastDealtCombination().Kind() == CombinationSingle {
		if !l.hasStrongCombination(l.unexploredCombinations) {
			return
		} else if game.GetMaxPlayerNumber() == 1 || rand.Intn(100) < 70 {
			// 70% remove 2 if not ch???t turn
			l.removeAllSingleCard()
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

		// n???u ng?????i ch??i tr?????c kh??ng c??n 2 th?? ch???n lu??n
		if !contains2 {
			l.removePass()
			return
		}

		// n???u ng?????i ch??i tr?????c c?? 2 th?? 100% ch???n n???u con v???a ????nh l?? 2 ????? v?? 90% ch???n n???u l?? 2 ??en
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

// lo???i con 2 ra n???u turn n??y m??nh kh??ng ph???i ch???n ai
func (l *LocalNode) remove2IfIsFirstTurn(game Game) {
	// check l???i n???u c?? 1 ng?????i c??n 1 con th?? kh??ng ???????c lo???i 2
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		if len(game.GetPlayerAt(i).AllAvailableCombinations()) == 1 {
			return
		}
	}
	// l???y qu??n b??i l??? g???n nh??? nh???t (nh??? th??? 2) m?? kh??ng n???m trong b??? n??o
	card := l.getAlmostSmallestSingleCardWhichNotInAnyCombination()
	// n???u kh??ng c?? b??i n??o ho???c l?? ???? c??ng l?? 2
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

// b???t bu???c ch???n con l??? n???u c?? th???
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

// x??a pass t???c l?? b???t bu???c ????nh
func (l *LocalNode) removePass() {
	for i := 0; i < len(l.unexploredCombinations); i++ {
		if l.unexploredCombinations[i].Kind() == CombinationPass {
			l.removeUnexploredCombination(l.unexploredCombinations[i])
			break
		}
	}
}

// n???u m???i ng?????i to??n c??n 1 qu??n th?? ????nh b??? tr?????c, b??? h???t qu??n l??? ra ngo??i ????nh b??? h???t tr?????c
func (l *LocalNode) removeSingleCardIfTheyAllHaveOneCardLeft(game Game) {
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		if game.GetPlayerAt(i).GetCardsLength() == 1 {
			continue
		}
		return
	}

	backup := make([]Combination, len(l.unexploredCombinations))
	copy(backup, l.unexploredCombinations)

	removedList := []int{}
	keep2List := []Combination{}
	for i := range l.unexploredCombinations {
		if l.unexploredCombinations[i].Kind() == CombinationSingle {
			removedList = append(removedList, i)
			if l.unexploredCombinations[i].Cards()[0].rank == Two {
				keep2List = append(keep2List, l.unexploredCombinations[i])
			}
		}
	}

	for i := len(removedList) - 1; i >=0; i-- {
		l.removeUnexploredCombinationAt(removedList[i])
	}

	if len(l.unexploredCombinations) == 0 {
		l.unexploredCombinations = keep2List
	}
	if len(l.unexploredCombinations) == 0 {
		l.unexploredCombinations = backup
	}
}

// n???u ?????i ph????ng c??n 1 con 2 v?? 1 con l??? (kh??ng ph???i 2)
// b??? c??c b??? m???nh h??n con 2 kia ??i n???u b??? ??i m?? v???n c?? qu??n l??? l???n h??n con l??? c??n l???i c???a ng?????i kia
// ch??? t??nh tr?????ng h???p 2 ng?????i ch??i
// trong turn kh??ng ph???i turn ch???t
func (l *LocalNode) removeCombinationStrongerThan2IfTheyHave2AndOneSmallSingleCardLeft(game Game) {
	if game.GetMaxPlayerNumber() != 2 {
		return
	}
	if game.GetPlayerAt(1 - game.GetCurrentPlayerIndex()).GetCardsLength() != 2 {
		return
	}
	com := game.GetPlayerAt(1 - game.GetCurrentPlayerIndex()).AllAvailableCombinations()
	if len(com) != 2 {
		return
	}
	if com[0].Kind() != CombinationSingle || com[1].Kind() != CombinationSingle {
		return
	}
	if com[0].(*SingleCard).card.rank != Two && com[1].(*SingleCard).card.rank != Two {
		return
	}
	if com[0].(*SingleCard).card.rank == Two && com[1].(*SingleCard).card.rank == Two {
		return
	}
	var card *SingleCard
	if com[0].(*SingleCard).card.rank == Two {
		card = com[1].(*SingleCard)
	} else {
		card = com[0].(*SingleCard)
	}

	rmList := []Combination{}
	for i := range l.unexploredCombinations {
		if l.unexploredCombinations[i].Kind() == CombinationThreeConsecutivePairs ||
			l.unexploredCombinations[i].Kind() == CombinationFourConsecutivePairs ||
			l.unexploredCombinations[i].Kind() == CombinationQuads {
			rmList = append(rmList, l.unexploredCombinations[i])
		}
	}
	rmListLen := len(rmList)
	for i := 0; i < rmListLen; i++ {
		rmList = append(rmList, game.GetCurrentPlayer().GetAllCombinationsHasSameAtLeastOneCardWith(rmList[i])...)
	}

	backup := make([]Combination, len(l.unexploredCombinations))
	copy(backup, l.unexploredCombinations)

	for i := range rmList {
		l.removeUnexploredCombination(rmList[i])
	}

	if len(l.unexploredCombinations) > 0 {
		for i := range l.unexploredCombinations {
			if l.unexploredCombinations[i].Defeats(card) {
				if card.card.rank != Ace {
					// n???u qu??n l??? kia kp qu??n ??t th?? b??? c??c b??? ????i A, tam A ra, ????nh c??c A c??u 2
					rmList = []Combination{}
					for j := range l.unexploredCombinations {
						if (l.unexploredCombinations[j].Kind() == CombinationDubs ||
							l.unexploredCombinations[j].Kind() == CombinationTrips) &&
							containsRank(l.unexploredCombinations[j].Cards(), Ace) {
							rmList = append(rmList, l.unexploredCombinations[j])
						}
					}
					for _, c := range rmList {
						l.removeUnexploredCombination(c)
					}
				}
				return
			}
		}
	}
	l.unexploredCombinations = backup
}

// x??a b??? kh???i list
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

// c?? t??? qu??, 3 ????i th??ng, 4 ????i th??ng
func (l *LocalNode) hasStrongCombination(combinations []Combination) bool {
	for i := range combinations {
		if combinations[i].Kind() == CombinationQuads ||
			combinations[i].Kind() == CombinationThreeConsecutivePairs ||
			combinations[i].Kind() == CombinationFourConsecutivePairs {
			return true
		}
	}
	return false
}

// l???y l?? l??? c?? gi?? tr??? nh??? th??? 2 m?? ko c?? trong b???t k?? b??? n??o
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

// true n???u t???t c??? m???i ng?????i ch??i kh??c c??n 1 l??
func (l *LocalNode) allHasOneCardLeft(game Game) bool {
	for i := 0; i < game.GetMaxPlayerNumber(); i++ {
		if i == game.GetCurrentPlayerIndex() {
			continue
		}
		if game.GetPlayerAt(i).GetCardsLength() == 1 {
			continue
		}
		return false
	}
	return true
}

func (l *LocalNode) removeAllSingleCard() {
	i := 0
	for _, c := range l.unexploredCombinations {
		if c.Kind() != CombinationSingle {
			l.unexploredCombinations[i] = c
			i++
		}
	}
	l.unexploredCombinations = l.unexploredCombinations[:i]
}

func (l *LocalNode) String() string {
	info := ""
	for _, node := range l.children {
		info += fmt.Sprintf("%-40s", "Node " + node.GetCombination().String())
		info += "|"
		info += fmt.Sprintf("%-20s", "Visit " + fmt.Sprintf("%d", node.GetVisit()))
		info += "|"
		info += fmt.Sprintf("%-30s", "Reward " + fmt.Sprintf("%+v", node.GetReward()))
		info += "\n"
	}
	return info
}