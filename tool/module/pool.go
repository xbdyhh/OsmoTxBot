package module

type Pools struct {
	Pools []struct {
		Type       string `json:"@type,omitempty"`
		Address    string `json:"address,omitempty"`
		ID         string `json:"id,omitempty"`
		PoolParams struct {
			SwapFee                  string      `json:"swap_fee,omitempty"`
			ExitFee                  string      `json:"exit_fee,omitempty"`
			SmoothWeightChangeParams interface{} `json:"smooth_weight_change_params,omitempty"`
		} `json:"pool_params,omitempty"`
		FuturePoolGovernor string `json:"future_pool_governor,omitempty"`
		TotalShares        struct {
			Denom  string `json:"denom,omitempty"`
			Amount string `json:"amount,omitempty"`
		} `json:"total_shares,omitempty"`
		PoolAssets []struct {
			Token struct {
				Denom  string `json:"denom,omitempty"`
				Amount string `json:"amount,omitempty"`
			} `json:"token,omitempty"`
			Weight string `json:"weight,omitempty"`
		} `json:"pool_assets,omitempty"`
		TotalWeight string `json:"total_weight,omitempty"`
	} `json:"pools,omitempty"`
	Pagination struct {
		NextKey string `json:"next_key,omitempty"`
		Total   string `json:"total,omitempty"`
	} `json:"pagination,omitempty"`
}
