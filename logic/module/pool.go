package module

import sdk "github.com/cosmos/cosmos-sdk/types"

type Path struct {
	ID         int64
	Ratio      float64
	WeightFrom uint64
	WeightTo   uint64
	AmountFrom uint64
	AmountTo   uint64
	fees       float64
}

type Pool struct {
	ID         int64
	PoolAssets PoolAssets
	SwapFees   float64
}
type PoolAssets []PoolAsset
type PoolAsset struct {
	Token  sdk.Coin
	Weight uint64
}

func (p *Path) GetRatio() {
	p.Ratio = float64(p.WeightFrom) / float64(p.WeightTo) * float64(p.AmountTo) / float64(p.AmountFrom) * (1 - p.fees)
}

type Router struct {
	PoolIds       []uint64
	TokenOutDenom []string
	Depth         uint64
	Ratio         float64
}
