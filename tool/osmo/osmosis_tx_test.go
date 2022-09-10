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

//func TestQueryOSMOPoolInfo(t *testing.T) {
//	InitCcontext()
//	ctx := tool.InitMyContext()
//	pool, err := QueryOsmoPoolInfo(ctx)
//	if err != nil {
//		t.Errorf("%v", err)
//	}
//	fmt.Println(len(pool.Pools))
//}

//func TestSendTx(t *testing.T) {
//	InitCcontext()
//	ctx := tool.InitMyContext()
//	SendOsmoTx(ctx, "",
//		"uosmo", "1000", 10000, 2, 656400, []uint64{1},
//		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"})
//}

//func TestSendTx2(t *testing.T) {
//	InitCcontext()
//	ctx := tool.InitMyContext()
//	SendOsmoTx2(ctx, &Ccontext, "enroll saddle syrup movie steak tunnel invest old trophy brown angry multiply", "osmo16kydz6vznpgtpgws733panrs6atdsefcfxa97j", 44, 656400, "uosmo", "100", 1000, []uint64{1, 1},
//		[]string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "uosmo"})
//}

func TestIsOsmoSuccess(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	ok, err := IsSendSuccess(ctx, "BBFBEFC889E1B506612405F3F88982AA333D928A3D0430D586012DB933B651B8", "7228BA7DCD97099E9E2F72EAE976D804541D54E0914AA38FB24ABC4816BAFF05")
	if err != nil || !ok {
		t.Errorf("an err happend%v", err)
	}
}

func TestQuerySimulate(t *testing.T) {
	InitCcontext()
	ctx := tool.InitMyContext()
	jsonmsg := "{\"body\":{\"messages\":[{\"@type\":\"/osmosis.gamm.v1beta1.MsgSwapExactAmountIn\",\"sender\":\"osmo16kydz6vznpgtpgws733panrs6atdsefcfxa97j\",\"routes\":[{\"poolId\":\"1\",\"tokenOutDenom\":\"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\"},{\"poolId\":\"22\",\"tokenOutDenom\":\"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076\"},{\"poolId\":\"42\",\"tokenOutDenom\":\"uosmo\"}],\"tokenIn\":{\"denom\":\"uosmo\",\"amount\":\"10535414\"},\"tokenOutMinAmount\":\"1\"}],\"memo\":\"\",\"timeout_height\":\"0\",\"extension_options\":[],\"non_critical_extension_options\":[]},\"auth_info\":{\"signer_infos\":[{\"public_key\":{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"ArBXS3ckpnjWCSVsPsGddT5veYQiQjcBeuDGH6zkPvAC\"},\"mode_info\":{\"single\":{\"mode\":\"SIGN_MODE_DIRECT\"}},\"sequence\":\"10415\"}],\"fee\":{\"amount\":[],\"gas_limit\":\"1000000\",\"payer\":\"\",\"granter\":\"\"}},\"signatures\":[\"K3cgv1iZmQba0tq37KU4UoLXSDfHI4ZWmpTa2KTOttM3jWENL3/8juECrh8A5Y0KbHD6jr7l51jVlRYldzUvFw==\"]}"
	ok, err := QuerySimulate(ctx, []byte(jsonmsg))
	if err != nil || !ok {
		t.Errorf("an err happend%v", err)
	}

}
