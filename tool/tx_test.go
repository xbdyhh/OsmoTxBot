package tool

import (
	"context"
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	InitContext()
	ctx := context.Background()
	acc, err := QueryAccountInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	fmt.Println(acc.Account.Sequence)
}

func TestQueryBalanceInfo(t *testing.T) {
	InitContext()
	ctx := context.Background()
	bal, err := QueryBalanceInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	amount, err := bal.Get("uosmo")
	fmt.Println(amount)
}

func TestQueryPoolInfo(t *testing.T) {
	InitContext()
	ctx := context.Background()
	pool, err := QueryPoolInfo(ctx)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(len(pool.Pools))
}
