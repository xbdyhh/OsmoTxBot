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
	"sync"
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
func (p PoolMap) FreshMap(ctx *tool.MyContext, pools []module.Pool) {
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
					if to.TokenDenom == OSMO_DENOM && to.Amount < 10000000 {
						continue
					}
					if from.TokenDenom == OSMO_DENOM && from.Amount < 10000000 {
						continue
					}

					p[from.TokenDenom][to.TokenDenom] = path
				} else {
					if oldpath.Ratio < path.Ratio &&
						float64(path.GetDepth(1.01))/p[OSMO_DENOM][from.TokenDenom].Ratio > 1000000 {
						p[from.TokenDenom][to.TokenDenom] = path
					}
				}
			}
		}
	}
	ctx.Logger.Info("fresh pools map done.")
	time.Sleep(3 * time.Second)
}

func (p PoolMap) FindProfitMargins(ctx *tool.MyContext) ([]module.Router, error) {
	routers := make([]module.Router, 0, 0)
	for from, frommap := range p {
		for to, v := range frommap {
			if from == OSMO_DENOM || to == OSMO_DENOM {
				continue
			}
			ratio := v.Ratio * p[OSMO_DENOM][from].Ratio * p[to][OSMO_DENOM].Ratio
			if ratio > 1 {
				ids := []uint64{p[OSMO_DENOM][from].ID, p[from][to].ID, p[to][OSMO_DENOM].ID}
				tokenOutDenoms := []string{from, to, OSMO_DENOM}
				depthFrom := p[OSMO_DENOM][from].GetDepth(ratio)
				depth := uint64(float64(v.GetDepth(ratio)) / p[OSMO_DENOM][from].Ratio)
				depthTo := uint64(float64(p[to][OSMO_DENOM].GetDepth(ratio)) / p[OSMO_DENOM][from].Ratio / p[from][to].Ratio)
				router := module.Router{
					PoolIds:       ids,
					TokenOutDenom: tokenOutDenoms,
					Depth:         MinDepth(depthFrom, depth, depthTo),
					Ratio:         ratio,
				}
				routers = append(routers, router)
			}
		}

	}
	return routers, nil
}

var TransactionRouters []module.Router

var RoterLock sync.Mutex

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
		pMap.FreshMap(ctx, pools)
		//遍历map去处与osmo直接相关的pool后，计算三角赔率（x*y*z*0.97*0.98*0.98），将大于1的添加入待执行名单
		//得到切片1
		routers, err := pMap.FindProfitMargins(ctx)
		//组合过滤
		RoterLock.Lock()
		TransactionRouters = SortRouters(ctx, routers)
		ctx.Logger.Infof("Finally transaction routers is:%v", routers)
		RoterLock.Unlock()
	}
	ctx.Wg.Done()
}

func SortRouters(ctx *tool.MyContext, routers []module.Router) []module.Router {
	newRouters := make([]module.Router, 0, 0)
	for _, v := range routers {
		if float64(v.Depth)*(v.Ratio-1) > osmo.GAS_FEE*20 {
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
		ok = ok || checkliquid < 100000000
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
	fmt.Println("address is:" + address)
	time.Sleep(3 * time.Second)
	fmt.Println("sleeping end!")
	acc, err := osmo.QueryOsmoAccountInfo(ctx, address)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
		panic(err)
	}
	seq, err := strconv.ParseUint(acc.Account.Sequence, 10, 64)
	if err != nil {
		ctx.Logger.Errorf("%v", err)
	}
	for {
		balance, err := osmo.QueryOsmoBalanceInfo(ctx, address)
		if err != nil {
			ctx.Logger.Errorf("%v", err)
			continue
		}
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
			continue
		}
		RoterLock.Lock()
		for i, v := range TransactionRouters {
			amountin := Min(v.Depth, balAmount-osmo.GAS_FEE)
			tokenMinOUt := strconv.FormatUint(amountin, 10)
			fmt.Println("send msg")
			resp, err := osmo.SendOsmoTx(ctx, MNEMONIC, OSMO_DENOM, tokenMinOUt, amountin, seq, accnum, v.PoolIds, v.TokenOutDenom)
			fmt.Println("finish at ")
			if err != nil {
				ctx.Logger.Errorf("%d tx err:%v", i, err)
				continue
			}
			if resp != nil && resp.Code == 0 {
				seq++
			}
			balAmount -= (amountin + osmo.GAS_FEE)
			if balAmount <= 0 {
				break
			}
		}
		TransactionRouters = TransactionRouters[:0]
		RoterLock.Unlock()
		time.Sleep(2 * time.Second)
	}
	ctx.Wg.Done()
}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
