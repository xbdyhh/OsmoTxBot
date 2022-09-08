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
					if to.TokenDenom == "uosmo" && to.Amount < 10000000 {
						continue
					}
					if from.TokenDenom == "uosmo" && from.Amount < 10000000 {
						continue
					}

					p[from.TokenDenom][to.TokenDenom] = path
				} else {
					if oldpath.Ratio < path.Ratio &&
						float64(path.GetDepth(1.01))/p["uosmo"][from.TokenDenom].Ratio > 1000000 {
						p[from.TokenDenom][to.TokenDenom] = path
					}
				}
			}
		}
	}
	ctx.Logger.Info("fresh pools map done.")
}

func (p PoolMap) FindProfitMargins(ctx *tool.MyContext) ([]module.Router, error) {
	routers := make([]module.Router, 0, 0)
	for from, frommap := range p {
		for to, v := range frommap {
			if from == "uosmo" || to == "uosmo" {
				continue
			}
			ratio := v.Ratio * p["uosmo"][from].Ratio * p[to]["uosmo"].Ratio
			if ratio > 1 {
				ids := []uint64{p["uosmo"][from].ID, p[from][to].ID, p[to]["uosmo"].ID}
				tokenOutDenoms := []string{from, to, "uomos"}
				depthFrom := p["uosmo"][from].GetDepth(ratio)
				depth := uint64(float64(v.GetDepth(ratio)) / p["uosmo"][from].Ratio)
				depthTo := uint64(float64(p[to]["uosmo"].GetDepth(ratio)) / p["uosmo"][from].Ratio / p[from][to].Ratio)
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
		fmt.Println(routers)
		break
	}
}

func SortRouters(ctx *tool.MyContext, routers []module.Router) []module.Router {
	newRouters := make([]module.Router, 0, 0)
	for _, v := range routers {
		if float64(v.Depth)*(v.Ratio-1) > 2000 {
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
		if v < ans || i == 1 {
			ans = v
		}
	}
	return ans
}
