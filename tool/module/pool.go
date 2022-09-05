package module

type Pools struct {
	Pools []struct {
		Type       string `json:"@type"`
		Address    string `json:"address"`
		ID         string `json:"id"`
		PoolParams struct {
			SwapFee                  string      `json:"swapFee"`
			ExitFee                  string      `json:"exitFee"`
			SmoothWeightChangeParams interface{} `json:"smoothWeightChangeParams"`
		} `json:"poolParams"`
		FuturePoolGovernor string `json:"future_pool_governor"`
		TotalShares        struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"totalShares"`
		PoolAssets []struct {
			Token struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"token"`
			Weight string `json:"weight"`
		} `json:"poolAssets"`
		TotalWeight string `json:"totalWeight"`
	} `json:"pools"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
