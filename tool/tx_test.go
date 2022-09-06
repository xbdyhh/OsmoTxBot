package tool

import (
	"context"
	"testing"
)

//func TestQuery(t *testing.T) {
//	InitContext()
//	ctx := context.Background()
//	acc, err := QueryAccountInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
//	if err != nil {
//		t.Errorf("an err happend%v", err)
//	}
//	fmt.Println(acc.Account.Sequence)
//}
//
//func TestQueryBalanceInfo(t *testing.T) {
//	InitContext()
//	ctx := context.Background()
//	bal, err := QueryBalanceInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
//	if err != nil {
//		t.Errorf("an err happend%v", err)
//	}
//	amount, err := bal.Get("uosmo")
//	fmt.Println(amount)
//}
//
//func TestQueryPoolInfo(t *testing.T) {
//	InitContext()
//	ctx := context.Background()
//	pool, err := QueryPoolInfo(ctx)
//	if err != nil {
//		t.Errorf("%v", err)
//	}
//	fmt.Println(len(pool.Pools))
//}

func TestSendTx(t *testing.T) {
	InitContext()
	ctx := context.Background()
	SendOsmoTx(ctx, "",
		"uosmo", "1000", 10000, 2, 656400, []string{"1"},
		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"})
}
