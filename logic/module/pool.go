package module

type Path struct {
	ID         uint64
	Ratio      float64
	WeightFrom uint64
	WeightTo   uint64
	AmountFrom float64
	AmountTo   float64
	Fees       float64
}

type Pool struct {
	ID         uint64
	PoolAssets PoolAssets
	SwapFees   float64
}
type PoolAssets []PoolAsset
type PoolAsset struct {
	TokenDenom string
	Amount     float64
	Weight     uint64
}

func (p *Path) GetRatio() {
	p.Ratio = (float64(p.WeightFrom) / float64(p.WeightTo)) *
		(p.AmountTo / p.AmountFrom) *
		(1 - p.Fees)
}
func (p Path) GetDepth() uint64 {
	return uint64(0.001 * float64(p.AmountFrom))
}

type Router struct {
	PoolIds       []uint64
	TokenOutDenom []string
	Depth         uint64
	Ratio         float64
}
