package logic

import (
	"fmt"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/xbdyhh/OsmoTxBot/logic/module"
	"github.com/xbdyhh/OsmoTxBot/tool"
	tm "github.com/xbdyhh/OsmoTxBot/tool/module"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const (
	MNEMONIC   = ""
	OSMO_DENOM = "uosmo"
)

// Pool[from][to]
type PoolMap map[string]map[string][]module.Path

func NewPoolMap(ctx *tool.MyContext) PoolMap {
	return make(map[string]map[string][]module.Path)
}
func (p PoolMap) FreshMap(ctx *tool.MyContext, pools []module.Pool, bal uint64) []module.Pool {
	TransactionRouters = TransactionRouters[:0]
	newPools := make([]module.Pool, 0, 0)
	for _, v := range pools {
		for _, from := range v.PoolAssets {
			for _, to := range v.PoolAssets {
				if from.TokenDenom == to.TokenDenom {
					continue
				}
				if _, ok := p[from.TokenDenom]; !ok {
					p[from.TokenDenom] = make(map[string][]module.Path)
				}
				if _, ok := p[from.TokenDenom][to.TokenDenom]; !ok {
					p[from.TokenDenom][to.TokenDenom] = make([]module.Path, 0, 0)
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
				path.GetRatio(bal)
				p[from.TokenDenom][to.TokenDenom] = append(p[from.TokenDenom][to.TokenDenom], path)
			}
		}
		if v.PoolAssets[0].TokenDenom != OSMO_DENOM && v.PoolAssets[1].TokenDenom != OSMO_DENOM {
			newPools = append(newPools, v)
		}
	}
	ctx.Logger.Debug("fresh pools map done.")
	return newPools
}

// 寻找兑换比率大于1的routers
func (p PoolMap) FindProfitMargins(ctx *tool.MyContext, pools []module.Pool, balance uint64) ([]module.Router, error) {
	routers := make([]module.Router, 0, 0)
	for fromkey, fromArr := range p[OSMO_DENOM] {
		for tokey, toMap := range p {
			if fromkey == tokey {
				for _, from := range fromArr {
					for _, to := range toMap[OSMO_DENOM] {
						if ratio := to.Ratio * from.Ratio; ratio > 1 && to.ID != from.ID {
							ids := []uint64{from.ID, to.ID}
							fmt.Println(ids)
							fmt.Println()
							out := []string{tokey, OSMO_DENOM}
							depth1 := from.GetDepth(balance)
							depth2 := uint64(float64(to.GetDepth(balance)) / from.Ratio)
							routers = append(routers, module.Router{
								PoolIds:       ids,
								TokenOutDenom: out,
								Depth:         MinDepth(depth1, depth2, balance),
								Ratio:         ratio,
							})
						}
					}
				}
			}
		}
	}

	for _, v := range pools {
	nextpool3:
		for _, fromAss := range v.PoolAssets {
			for _, toAss := range v.PoolAssets {
				if toAss.TokenDenom == fromAss.TokenDenom {
					continue
				}
				path := module.Path{
					ID:         v.ID,
					Ratio:      0,
					WeightFrom: fromAss.Weight,
					WeightTo:   toAss.Weight,
					AmountFrom: fromAss.Amount,
					AmountTo:   toAss.Amount,
					Fees:       v.SwapFees,
				}
				path.GetRatio(balance)
				var router module.Router
				for _, from := range p[OSMO_DENOM][fromAss.TokenDenom] {
					for _, to := range p[toAss.TokenDenom][OSMO_DENOM] {
						fmt.Println(path.Ratio, from.Ratio, to.Ratio)
						if ratio := path.Ratio * from.Ratio * to.Ratio; ratio > 1 {
							ids := []uint64{from.ID, path.ID, to.ID}
							out := []string{fromAss.TokenDenom, toAss.TokenDenom, OSMO_DENOM}
							depth1 := from.GetDepth(balance)
							depth2 := uint64(float64(path.GetDepth(balance)) / from.Ratio)
							depth3 := uint64(float64(to.GetDepth(balance)) / from.Ratio / path.Ratio)
							depth := MinDepth(depth1, depth2, depth3, balance)
							if float64(depth)*ratio > float64(router.Depth)*router.Ratio {
								router = module.Router{
									PoolIds:       ids,
									TokenOutDenom: out,
									Depth:         depth,
									Ratio:         ratio,
								}
							}
						}
					}
				}
				if router.PoolIds != nil {
					routers = append(routers, router)
				}
				break nextpool3
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
				path.GetRatio(balance)
				var router module.Router
				for path3key, path3Arr := range p[to.TokenDenom] {
					for _, path3 := range path3Arr {
						for _, osmoFrom := range p[OSMO_DENOM][from.TokenDenom] {
							for _, osmoTo := range p[to.TokenDenom][OSMO_DENOM] {
								if ratio := path.Ratio * osmoFrom.Ratio * path3.Ratio * osmoTo.Ratio; ratio > 1 {
									ids := []uint64{osmoFrom.ID, path.ID, path3.ID, osmoTo.ID}
									out := []string{from.TokenDenom, to.TokenDenom, path3key, OSMO_DENOM}
									depth1 := osmoFrom.GetDepth(balance)
									depth2 := uint64(float64(path.GetDepth(balance)) / osmoFrom.Ratio)
									depth3 := uint64(float64(path3.GetDepth(balance)) / osmoFrom.Ratio / path.Ratio)
									depth4 := uint64(float64(osmoTo.GetDepth(balance)) / osmoFrom.Ratio / path.Ratio / path3.Ratio)
									depth := MinDepth(depth1, depth2, depth3, depth4, balance)
									if float64(depth)*ratio > float64(router.Depth)*router.Ratio {
										router = module.Router{
											PoolIds:       ids,
											TokenOutDenom: out,
											Depth:         depth,
											Ratio:         ratio,
										}
									}
								}
							}
						}
					}
				}
				if router.PoolIds != nil {
					routers = append(routers, router)
				}
				break next4pool
			}
		}
	}

	return routers, nil
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

		//删除流动性小于1000的pool
		pools, err := DeleteLittlePools(ctx, res)
		//生成最低直接路径的pool map
		newpools := pMap.FreshMap(ctx, pools, balamount)
		//遍历map去处与osmo直接相关的pool后，计算三角赔率（x*y*z*0.97*0.98*0.98），将大于1的添加入待执行名单
		//得到切片1
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
	sort.SliceStable(routers, func(i, j int) bool {
		return float64(routers[i].Depth)*routers[i].Ratio > float64(routers[j].Depth)*routers[j].Ratio
	})
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
