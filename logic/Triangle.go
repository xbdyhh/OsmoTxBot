package logic

import (
	"fmt"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/xbdyhh/OsmoTxBot/logic/module"
	"github.com/xbdyhh/OsmoTxBot/tool"
	tm "github.com/xbdyhh/OsmoTxBot/tool/module"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
	"sort"
	"strconv"
	"time"
)

const (
	MNEMONIC   = ""
	OSMO_DENOM = "uosmo"
)

var totalPath = 0

// Pool[from][to]
type PoolMap map[string]map[string][]module.Path

func NewPoolMap(ctx *tool.MyContext) PoolMap {
	return make(map[string]map[string][]module.Path)
}
func (p PoolMap) FreshMap(ctx *tool.MyContext, pools []module.Pool) {
	TransactionRouters = TransactionRouters[:0]
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
				path.GetRatio()
				if !IsPoolIn(p[from.TokenDenom][to.TokenDenom], path.ID) {
					p[from.TokenDenom][to.TokenDenom] = append(p[from.TokenDenom][to.TokenDenom], path)
				}
			}
		}
	}
	ctx.Logger.Debug("fresh pools map done.")
}

func IsPoolIn(ids []module.Path, id uint64) bool {
	for _, v := range ids {
		if v.ID == id {
			return true
		}
	}
	return false
}

// 寻找兑换比率大于1的routers
func (p PoolMap) FindProfitMargins(ctx *tool.MyContext, balance uint64) ([]module.Router, error) {
	routers := make([]module.Router, 0, 0)
	ids := make([]uint64, 0, 0)
	denoms := make([]string, 0, 0)
	routers = p.FindPath(ctx, ids, 1000000000, 1, denoms, OSMO_DENOM)
	fmt.Println("total paths is:", totalPath)
	return routers, nil
}

func (p PoolMap) FindPath(ctx *tool.MyContext, oldids []uint64, depth uint64, ratio float64, olddenoms []string, denom string) []module.Router {
	ids := make([]uint64, len(oldids), len(oldids))
	copy(ids, oldids)
	routers := make([]module.Router, 0, 0)
	denoms := make([]string, len(oldids), len(oldids))
	copy(denoms, olddenoms)
	if len(denoms) != 0 && denoms[len(denoms)-1] == OSMO_DENOM {
		totalPath++
		if ratio > 1 && depth > 1000000 {
			routers = append(routers, module.Router{
				PoolIds:       ids,
				TokenOutDenom: denoms,
				Depth:         depth,
				Ratio:         ratio,
			})
		}
		return routers
	}
	if len(ids) > 2 {
		if patharr, ok := p[denom][OSMO_DENOM]; ok {
			for _, path := range patharr {
				depth2 := uint64(float64(path.GetDepth()) / ratio)
				totalPath++
				if ratio > 1 && depth > 1000000 {
					routers = append(routers, module.Router{
						PoolIds:       append(ids, path.ID),
						TokenOutDenom: append(denoms, OSMO_DENOM),
						Depth:         MinDepth(depth2, depth),
						Ratio:         ratio * path.Ratio,
					})
				}
			}
		}
		return routers
	}
	for key, patharr := range p[denom] {
		for _, path := range patharr {
			depth2 := uint64(float64(path.GetDepth()) / ratio)
			newrouters := make([]module.Router, 0, 0)
			if !IsIdIn(ids, path.ID) {
				newrouters = p.FindPath(ctx, append(ids, path.ID), MinDepth(depth2, depth), ratio*path.Ratio, append(denoms, key), key)
			}
			routers = append(routers, newrouters...)
		}
	}
	return routers

}

func IsIdIn(ids []uint64, id uint64) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

var TransactionRouters []module.Router

func FreshPoolMap(ctx *tool.MyContext) {
	priv, err := tool.NewPrivateKeyByMnemonic(MNEMONIC)
	if err != nil {

		panic(err)
	}
	address, err := osmo.PrivateToOsmoAddress(priv)
	if err != nil {
		panic(err)
	}
	for {
		totalPath = 0
		pMap := NewPoolMap(ctx)
		//拉取数据
		res, err := osmo.QueryOsmoPoolInfo(ctx)
		if err != nil {
			ctx.Logger.Errorf("pull pools err:%v", err)
			continue
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
		pMap.FreshMap(ctx, pools)
		//遍历map去处与osmo直接相关的pool后，计算三角赔率（x*y*z*0.97*0.98*0.98），将大于1的添加入待执行名单
		//得到切片1
		routers, err := pMap.FindProfitMargins(ctx, balamount)
		//组合过滤
		TransactionRouters = SortRouters(ctx, routers)
		ctx.Logger.Debugf("Finally transaction routers is:%v\n", routers)
		SendOsmoTriTx(ctx)
	}
}

func SortRouters(ctx *tool.MyContext, routers []module.Router) []module.Router {
	newRouters := make([]module.Router, 0, 0)
	newRouters = routers
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
		//for _, val := range v.PoolAssets {
		//	ok2, _ := regexp.Match("gamm.*", []byte(val.Token.Denom))
		//	ok = ok || ok2
		//}
		num, _ := strconv.ParseUint(v.PoolAssets[0].Token.Amount, 10, 64)
		if num > 100000 && v.ID != "551" {
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
			amountin -= osmo.GAS_FEE*uint64(len(v.PoolIds)) + 2500
		}
		tokenMinOut := strconv.FormatUint(amountin+uint64(len(v.PoolIds))*osmo.GAS_FEE+2500, 10)
		//判断利润是否达标
		if float64(amountin)*(v.Ratio-1) > 10000 {
			fmt.Printf("hope profit is: %v:amount is %d:ratio is %v:bal is %v:depth is %v,path is %v \n",
				float64(amountin)*(v.Ratio-1), amountin, v.Ratio, balAmount, v.Depth, v.PoolIds)
			ctx.Logger.Debugf("hope profit is: %v:amount is %d:ratio is %v\n", float64(amountin)*(v.Ratio-1), amountin-osmo.GAS_FEE*uint64(len(v.PoolIds))+2500, v.Ratio)
			resp, err := osmo.SendOsmoTx(ctx, MNEMONIC, OSMO_DENOM, tokenMinOut, amountin, seq, accnum, v.PoolIds, v.TokenOutDenom, int64(osmo.GAS_FEE*uint64(len(v.PoolIds))+2500))
			if err != nil {
				ctx.Logger.Errorf("%d tx err:%v", i, err)
				continue
			}
			if resp == nil {
				continue
			} else if resp.Code == 0 {
				seq++
				balAmount -= amountin + osmo.GAS_FEE*uint64(len(v.PoolIds)) + 2500
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
		if balAmount <= osmo.GAS_FEE*2 {
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
		time.Sleep(6 * time.Second)
	}
}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
