package tienlen_bot

import (
	"fmt"
	"testing"
)

func TestGetSequence(t *testing.T) {
	t.Log(fmt.Sprintf("%+v", GetSequence(SortCard(ParseCards("3♣, 3♥,4♥, 5♥, 6♣")))))
}
