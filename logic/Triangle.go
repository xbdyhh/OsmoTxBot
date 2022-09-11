package logic

import (
	"fmt"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/xbdyhh/OsmoTxBot/logic/module"
	"github.com/xbdyhh/OsmoTxBot/tool"
	tm "github.com/xbdyhh/OsmoTxBot/tool/module"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
	"regexp"
	"strconv"
	"time"
)

const (
	MNEMONIC   = ""
	OSMO_DENOM = "uosmo"
)

// Pool[from][to]
type PoolMap map[string]map[string]module.Path

func NewPoolMap(ctx *tool.MyContext) PoolMap {
	return make(map[string]map[string]module.Path)
}
func (p PoolMap) FreshMap(ctx *tool.MyContext, pools []module.Pool) []module.Pool {
	TransactionRouters = TransactionRouters[:0]
	newPools := make([]module.Pool, 0, 0)
	for _, v := range pools {
		for _, from := range v.PoolAssets {
			for _, to := range v.PoolAssets {
				if from.TokenDenom == to.TokenDenom {
					continue
				}
				if _, ok := p[from.TokenDenom]; !ok {
					p[from.TokenDenom] = make(map[string]module.Path)
				}
				path := module.Path{
					ID:         v.ID,
					Ratio:      0,
					WeightFrom: from.Weight,
					WeightTo:   to.Weight,
					AmountFrom: from.Amount,
					AmountTo:   to.Amount,
					Fees:       v.SwapFees,
				}
				path.GetRatio()
				oldpath, ok := p[from.TokenDenom][to.TokenDenom]
				if !ok {
					p[from.TokenDenom][to.TokenDenom] = path
				} else {
					if to.TokenDenom == OSMO_DENOM && to.Amount < 10000000 {
						continue
					}
					if from.TokenDenom == OSMO_DENOM && from.Amount < 10000000 {
						continue
					}
					if oldpath.Ratio < path.Ratio {
						p[from.TokenDenom][to.TokenDenom] = path
					}
				}
			}
		}
		if v.PoolAssets[0].TokenDenom != OSMO_DENOM && v.PoolAssets[1].TokenDenom != OSMO_DENOM {
			newPools = append(newPools, v)
		}
	}
	ctx.Logger.Debug("fresh pools map done.")
	time.Sleep(3 * time.Second)
	return newPools
}

// 寻找兑换比率大于1的routers
func (p PoolMap) FindProfitMargins(ctx *tool.MyContext, pools []module.Pool, balance uint64) ([]module.Router, error) {
	routers := make([]module.Router, 0, 0)
	for _, v := range pools {
	nextpool:
		for _, from := range v.PoolAssets {
			for _, to := range v.PoolAssets {
				if to.TokenDenom == from.TokenDenom {
					continue
				}
				path := module.Path{
					ID:         v.ID,
					Ratio:      0,
					WeightFrom: from.Weight,
					WeightTo:   to.Weight,
					AmountFrom: from.Amount,
					AmountTo:   to.Amount,
					Fees:       v.SwapFees,
				}
				path.GetRatio()
				if ratio := path.Ratio * p[OSMO_DENOM][from.TokenDenom].Ratio * p[to.TokenDenom][OSMO_DENOM].Ratio; ratio > 1 {
					fmt.Println(ratio, " ", path.GetDepth())
					ids := []uint64{p[OSMO_DENOM][from.TokenDenom].ID, path.ID, p[to.TokenDenom][OSMO_DENOM].ID}
					fmt.Println(ids)
					fmt.Println()
					out := []string{from.TokenDenom, to.TokenDenom, OSMO_DENOM}
					depth1 := p[OSMO_DENOM][from.TokenDenom].GetDepth()
					depth2 := uint64(float64(path.GetDepth()) / p[OSMO_DENOM][from.TokenDenom].Ratio)
					depth3 := uint64(float64(p[to.TokenDenom][OSMO_DENOM].GetDepth()) / p[OSMO_DENOM][from.TokenDenom].Ratio / path.Ratio)
					routers = append(routers, module.Router{
						PoolIds:       ids,
						TokenOutDenom: out,
						Depth:         MinDepth(depth1, depth2, depth3, balance),
						Ratio:         ratio,
					})
					break nextpool
				}
			}
		}
	}
	for _, v := range pools {
	next4pool:
		for _, from := range v.PoolAssets {
			for _, to := range v.PoolAssets {
				if to.TokenDenom == from.TokenDenom {
					continue
				}
				path := module.Path{
					ID:         v.ID,
					Ratio:      0,
					WeightFrom: from.Weight,
					WeightTo:   to.Weight,
					AmountFrom: from.Amount,
					AmountTo:   to.Amount,
					Fees:       v.SwapFees,
				}
				path.GetRatio()
				for tokendemnom3, path3 := range p[to.TokenDenom] {
					if ratio := path.Ratio * p[OSMO_DENOM][from.TokenDenom].Ratio * path3.Ratio * p[tokendemnom3][OSMO_DENOM].Ratio; ratio > 1 {
						fmt.Println(ratio, " ", path.GetDepth())
						ids := []uint64{p[OSMO_DENOM][from.TokenDenom].ID, path.ID, path3.ID, p[tokendemnom3][OSMO_DENOM].ID}
						fmt.Println(ids)
						fmt.Println()
						out := []string{from.TokenDenom, to.TokenDenom, tokendemnom3, OSMO_DENOM}
						depth1 := p[OSMO_DENOM][from.TokenDenom].GetDepth()
						depth2 := uint64(float64(path.GetDepth()) / p[OSMO_DENOM][from.TokenDenom].Ratio)
						depth3 := uint64(float64(p[to.TokenDenom][tokendemnom3].GetDepth()) / p[OSMO_DENOM][from.TokenDenom].Ratio / path.Ratio)
						depth4 := uint64(float64(p[tokendemnom3][OSMO_DENOM].GetDepth()) / p[OSMO_DENOM][from.TokenDenom].Ratio / path.Ratio / path3.Ratio)
						routers = append(routers, module.Router{
							PoolIds:       ids,
							TokenOutDenom: out,
							Depth:         MinDepth(depth1, depth2, depth3, depth4, balance),
							Ratio:         ratio,
						})
						break next4pool
					}
				}

			}
		}
	}

	routers = CombineRouters(ctx, routers)
	return routers, nil
}

// 组合routers以提高成功率
func CombineRouters(ctx *tool.MyContext, routers []module.Router) []module.Router {
	ctx.Logger.Debug("meta router is:", routers)
	newrouters := make([]module.Router, 0, 0)
	userouter := make(map[int]bool)
	for i, v := range routers {
		if userouter[i] {
			continue
		}
		router := v
		for i2, to := range routers[i:] {
			if router.TokenOutDenom[len(router.TokenOutDenom)-2] == to.TokenOutDenom[1] &&
				float64(router.Depth)/float64(to.Depth) > 0.9 && float64(router.Depth)/float64(to.Depth) < 1.1 && !userouter[i2] {
				router.PoolIds = append(router.PoolIds[0:len(router.TokenOutDenom)-2], to.PoolIds[1:]...)
				router.TokenOutDenom = append(router.TokenOutDenom[0:len(router.TokenOutDenom)-2], to.TokenOutDenom[1:]...)
				router.Depth = MinDepth(router.Depth, to.Depth)
				router.Ratio = router.Ratio * to.Ratio
				userouter[i] = true
			}
		}
		newrouters = append(newrouters, router)
	}
	return newrouters
}

var TransactionRouters []module.Router

func FreshPoolMap(ctx *tool.MyContext) {
	pMap := NewPoolMap(ctx)
	for {
		//拉取数据
		res, err := osmo.QueryOsmoPoolInfo(ctx)
		if err != nil {
			ctx.Logger.Errorf("pull pools err:%v", err)
			continue
		}
		//删除流动性小于1000的pool
		pools, err := DeleteLittlePools(ctx, res)
		//生成最低直接路径的pool map
		newpools := pMap.FreshMap(ctx, pools)
		//遍历map去处与osmo直接相关的pool后，计算三角赔率（x*y*z*0.97*0.98*0.98），将大于1的添加入待执行名单
		//得到切片1
		priv, err := tool.NewPrivateKeyByMnemonic(MNEMONIC)
		if err != nil {

			panic(err)
		}
		address, err := osmo.PrivateToOsmoAddress(priv)
		if err != nil {
			panic(err)
		}

		balance, err := osmo.QueryOsmoBalanceInfo(ctx, address)
		var balamount uint64
		for _, v := range balance.Balances {
			if v.Denom == "uosmo" {
				balamount, _ = strconv.ParseUint(v.Amount, 10, 64)
			}
		}
		routers, err := pMap.FindProfitMargins(ctx, newpools, balamount)
		//组合过滤
		TransactionRouters = SortRouters(ctx, routers)
		ctx.Logger.Debugf("Finally transaction routers is:%v\n", routers)
		SendOsmoTriTx(ctx)
	}
}

func SortRouters(ctx *tool.MyContext, routers []module.Router) []module.Router {
	newRouters := make([]module.Router, 0, 0)
	for _, v := range routers {
		if float64(v.Depth)*(v.Ratio-1) > osmo.GAS_FEE {
			newRouters = append(newRouters, v)
		}
	}
	//sort.SliceStable(routers, func(i, j int) bool {
	//	return float64(routers[i].Depth)*routers[i].Ratio > float64(routers[j].Depth)*routers[j].Ratio
	//})
	return newRouters
}

func DeleteLittlePools(ctx *tool.MyContext, pools *tm.Pools) ([]module.Pool, error) {
	if len(pools.Pools) == 0 {
		ctx.Logger.Errorf("Input pool data is null!")
		return nil, errors.Errorf("The pools is no data!!!")
	}
	ans := make([]module.Pool, 0, 0)
	for _, v := range pools.Pools {
		ok := false
		for _, val := range v.PoolAssets {
			ok2, _ := regexp.Match("gamm.*", []byte(val.Token.Denom))
			ok = ok || ok2
		}
		checkliquid, _ := strconv.ParseUint(v.PoolAssets[0].Token.Amount, 10, 64)
		checkliquid2, _ := strconv.ParseUint(v.PoolAssets[1].Token.Amount, 10, 64)
		ok = ok || checkliquid < 100000000 || checkliquid2 < 100000000
		if !ok {
			pool := module.Pool{}
			pool.ID, _ = strconv.ParseUint(v.ID, 10, 64)
			pool.PoolAssets = module.PoolAssets{}
			for _, asset := range v.PoolAssets {
				var poolAsset = module.PoolAsset{}
				var err error
				poolAsset.TokenDenom = asset.Token.Denom
				poolAsset.Amount, err = strconv.ParseFloat(asset.Token.Amount, 64)
				if err != nil {
					return nil, err
				}
				poolAsset.Weight, _ = strconv.ParseUint(asset.Weight, 10, 64)
				pool.PoolAssets = append(pool.PoolAssets, poolAsset)
			}
			pool.SwapFees, _ = strconv.ParseFloat(v.PoolParams.SwapFee, 64)
			ans = append(ans, pool)
		}
	}
	return ans, nil
}

func MinDepth(uis ...uint64) uint64 {
	var ans uint64 = 0
	for i, v := range uis {
		if v < ans || i == 0 {
			ans = v
		}
	}
	return ans
}

func SendOsmoTriTx(ctx *tool.MyContext) {
	priv, err := tool.NewPrivateKeyByMnemonic(MNEMONIC)
	if err != nil {

		panic(err)
	}
	address, err := osmo.PrivateToOsmoAddress(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("send osmo triangle start!!!")
	acc, err := osmo.QueryOsmoAccountInfo(ctx, address)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
		return
	}
	seq, err := strconv.ParseUint(acc.Account.Sequence, 10, 64)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
	}

	balance, err := osmo.QueryOsmoBalanceInfo(ctx, address)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
	}
	//查询钱包余额
	var balAmount uint64
	for _, v := range balance.Balances {
		if v.Denom == OSMO_DENOM {
			balAmount, err = strconv.ParseUint(v.Amount, 10, 64)
			if err != nil {
				ctx.Logger.Errorf("%v", err)
			}
			break
		}
	}
	accnum, err := strconv.ParseUint(acc.Account.AccountNumber, 10, 64)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
	}
	txs := make([]string, 0, 0)
	for i, v := range TransactionRouters {
		amountin := Min(v.Depth, balAmount)
		if amountin == balAmount {
			amountin -= osmo.GAS_FEE
		}
		tokenMinOut := strconv.FormatUint(amountin+osmo.GAS_FEE, 10)
		//判断利润是否达标
		if float64(amountin)*(v.Ratio-1) > float64(osmo.GAS_FEE) {
			fmt.Printf("hope profit is: %v:amount is %d:ratio is %v:bal is %v:depth is %v \n",
				float64(amountin)*(v.Ratio-1), amountin, v.Ratio, balAmount, v.Depth)
			ctx.Logger.Debugf("hope profit is: %v:amount is %d:ratio is %v\n", float64(amountin)*(v.Ratio-1), amountin-osmo.GAS_FEE, v.Ratio)
			resp, err := osmo.SendOsmoTx(ctx, MNEMONIC, OSMO_DENOM, tokenMinOut, amountin, seq, accnum, v.PoolIds, v.TokenOutDenom)
			if err != nil {
				ctx.Logger.Errorf("%d tx err:%v", i, err)
				continue
			}
			if resp == nil {
				continue
			} else if resp.Code == 0 {
				seq++
				balAmount -= amountin + osmo.GAS_FEE
				txs = append(txs, resp.TxHash)

			} else if resp.Code == 32 {
				acc, err := osmo.QueryOsmoAccountInfo(ctx, address)
				if err != nil {
					ctx.Logger.Errorf("%v", err)
					panic(err)
				}
				seq, err = strconv.ParseUint(acc.Account.Sequence, 10, 64)
			}

		}
		if balAmount <= osmo.GAS_FEE*3 {
			break
		}
	}
	ctx.Logger.Infof("finish send msg!!!"+"tx is:", txs)
	//等待全部交易完成
	for {
		ok, err := osmo.IsSendSuccess(ctx, txs...)
		if err != nil {
			ctx.Logger.Errorf("query tx err happend!:%v\n", err)
		}
		fmt.Println("query tx response is:", ok)
		if ok {
			break
		}
		time.Sleep(7 * time.Second)
	}

}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
