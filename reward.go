package tienlen_bot

const (
	FactorPly               float64 = 0.005
	NumberOfCardsAtLateGame int     = 6
	Multiplier              float64 = 1
	FactorRed2sCard         float64 = 0.1 * Multiplier
	FactorBlack2sCard       float64 = 0.05 * Multiplier
	FactorNormalSingleCard  float64 = 0.01 * Multiplier
	FactorThreePairs        float64 = 0.1 * Multiplier
	FactorFourPairs         float64 = 0.3 * Multiplier
	FactorQuads             float64 = 0.2 * Multiplier
)

type Reward interface {
	AddReward(reward Reward)
	SetScore(playerIndex int, score float64)
	GetScoreOfPlayer(playerIndex int) float64
	Copy() Reward
}

type LocalReward struct {
	scores             []float64
	maxNumberOfPlayers int
}

func NewReward(maxNumberOfPlayers int) Reward {
	return &LocalReward{
		scores:             make([]float64, maxNumberOfPlayers),
		maxNumberOfPlayers: maxNumberOfPlayers,
	}
}

func (l *LocalReward) AddReward(reward Reward) {
	for i := 0; i < l.maxNumberOfPlayers; i++ {
		l.scores[i] += reward.GetScoreOfPlayer(i)
	}
}

func (l *LocalReward) SetScore(playerIndex int, score float64) {
	l.scores[playerIndex] = score
}

func (l *LocalReward) GetScoreOfPlayer(playerIndex int) float64 {
	return l.scores[playerIndex]
}

func (l *LocalReward) Copy() Reward {
	r := &LocalReward{
		scores:             make([]float64, l.maxNumberOfPlayers),
		maxNumberOfPlayers: l.maxNumberOfPlayers,
	}
	copy(r.scores, l.scores)
	return r
}
