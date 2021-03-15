package tienlen_bot

import (
	"fmt"
	"testing"
)

func TestGetSequence(t *testing.T) {
	t.Log(fmt.Sprintf("%+v", GetSequence(sortCard(parseCards("3♣, 3♥,4♥, 5♥, 6♣")))))
}
