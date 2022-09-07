package logic

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/xbdyhh/OsmoTxBot/logic/module"
	"github.com/xbdyhh/OsmoTxBot/tool"
	tm "github.com/xbdyhh/OsmoTxBot/tool/module"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
	"strconv"
	"time"
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
				if from.Token.Denom == to.Token.Denom {
					break
				}
				if _, ok := p[from.Token.Denom]; !ok {
					p[from.Token.Denom] = make(map[string]module.Path)
				}
				path := module.Path{
					ID:         v.ID,
					Ratio:      0,
					WeightFrom: from.Weight,
					WeightTo:   to.Weight,
					AmountFrom: from.Token.Amount.Uint64(),
					AmountTo:   to.Token.Amount.Uint64(),
				}
				path.GetRatio()
				oldpath, ok := p[from.Token.Denom][to.Token.Denom]
				if !ok {
					p[from.Token.Denom][to.Token.Denom] = path
				} else {
					if oldpath.Ratio < path.Ratio {
						p[from.Token.Denom][to.Token.Denom] = path
					}
				}
			}
		}
	}
	ctx.Logger.Info("fresh pools map done.")
}

func FreshPoolMap(ctx *tool.MyContext) {
	pMap := NewPoolMap(ctx)
	for {
		//拉取数据
		res, err := osmo.QueryOsmoPoolInfo(ctx)
		st := time.Now()
		fmt.Println("finish query osmo pool info!")
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
		//计算切片深度（depth1 = ）
		//组合过滤
		//发起交易
		//储存数据
		ans := time.Since(st)
		fmt.Println(ans)
		break
	}
}

func DeleteLittlePools(ctx *tool.MyContext, pools *tm.Pools) ([]module.Pool, error) {
	if len(pools.Pools) == 0 {
		ctx.Logger.Errorf("Input pool data is null!")
		return nil, errors.Errorf("The pools is no data!!!")
	}
	ans := make([]module.Pool, 0, 0)
	for _, v := range pools.Pools {
		if totalWeight, _ := strconv.ParseUint(v.TotalWeight, 10, 64); totalWeight > 1000 {
			pool := module.Pool{}
			pool.ID, _ = strconv.ParseInt(v.ID, 10, 64)
			pool.PoolAssets = module.PoolAssets{}
			for _, asset := range v.PoolAssets {
				var poolAsset = module.PoolAsset{}
				amount, _ := strconv.ParseInt(asset.Token.Amount, 10, 64)
				poolAsset.Token = sdk.NewInt64Coin(asset.Token.Denom, amount)
				poolAsset.Weight, _ = strconv.ParseUint(asset.Weight, 10, 64)
				pool.PoolAssets = append(pool.PoolAssets, poolAsset)
			}
			pool.SwapFees, _ = strconv.ParseFloat(v.PoolParams.SwapFee, 64)
			ans = append(ans, pool)
		}
	}
	return ans, nil
}
