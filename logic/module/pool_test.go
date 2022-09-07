package module

import (
	"fmt"
	"testing"
)

func TestPool_GetRatio(t *testing.T) {
	p := &Pool{
		id:         1,
		Ratio:      0,
		WeightFrom: 8,
		WeightTo:   2,
		AmountFrom: 22308458,
		AmountTo:   24635,
		fees:       0.002,
	}
	fmt.Println(p.GetRatio())
}
