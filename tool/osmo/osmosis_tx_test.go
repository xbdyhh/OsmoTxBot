package osmo

import (
	"fmt"
	"github.com/xbdyhh/OsmoTxBot/tool"
	"testing"
)

func TestQueryOsmoAccountInfo(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	acc, err := QueryOsmoAccountInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	fmt.Println(acc.Account.Sequence)
}

func TestQueryOsmoBalanceInfo(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	bal, err := QueryOsmoBalanceInfo(ctx, "osmo1w3t6kvkvhudyrcvveu9yzyh3sv7ykpst24rc4p")
	if err != nil {
		t.Errorf("an err happend%v", err)
	}
	amount, err := bal.Get("uosmo")
	fmt.Println(amount)
}

func TestQueryOSMOPoolInfo(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	pool, err := QueryOsmoPoolInfo(ctx)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(len(pool.Pools))
}

//func TestSendTx(t *testing.T) {
//	InitCcontext()
//	ctx := tool.InitMyContext()
//	SendOsmoTx(ctx, "",
//		"uosmo", "1000", 10000, 2, 656400, []uint64{1},
//		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"})
//}

func TestSendTx2(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	SendOsmoTx2(ctx, &Ccontext, "enroll saddle syrup movie steak tunnel invest old trophy brown angry multiply", "osmo16kydz6vznpgtpgws733panrs6atdsefcfxa97j", 44, 656400, "uosmo", "100", 1000, []uint64{1, 1},
		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "uosmo"})
}
