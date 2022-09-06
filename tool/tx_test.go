package tool

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	InitCcontext()
	ctx := InitMyContext()
	acc, err := QueryOsmoAccountInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	fmt.Println(acc.Account.Sequence)
}

func TestQueryBalanceInfo(t *testing.T) {
	InitCcontext()
	ctx := InitMyContext()
	bal, err := QueryOsmoBalanceInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	amount, err := bal.Get("uosmo")
	fmt.Println(amount)
}

func TestQueryPoolInfo(t *testing.T) {
	InitCcontext()
	ctx := InitMyContext()
	pool, err := QueryOsmoPoolInfo(ctx)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(len(pool.Pools))
}

func TestSendTx(t *testing.T) {
	InitCcontext()
	ctx := InitMyContext()
	SendOsmoTx(ctx, "",
		"uosmo", "1000", 10000, 2, 656400, []string{"1"},
		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"})
}
