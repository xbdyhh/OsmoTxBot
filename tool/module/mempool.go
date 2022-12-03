package module

type MemTxs struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}
type Result struct {
	NTxs       string   `json:"n_txs"`
	Total      string   `json:"total"`
	TotalBytes string   `json:"total_bytes"`
	Txs        []string `json:"txs"`
}

type MemTx struct {
	Body struct {
		Messages []struct {
			Type   string `json:"@type"`
			Sender string `json:"sender,omitempty"`
			Routes []struct {
				PoolID        string `json:"pool_id,omitempty"`
				TokenOutDenom string `json:"token_out_denom,omitempty"`
				TokenInDenom  string `json:"token_in_denom,omitempty"`
			} `json:"routes,omitempty"`
			TokenIn struct {
				Denom  string `json:"denom,omitempty"`
				Amount string `json:"amount,omitempty"`
			} `json:"token_in,omitempty"`
			TokenOut struct {
				Denom  string `json:"denom,omitempty"`
				Amount string `json:"amount,omitempty"`
			} `json:"token_out,omitempty"`
		} `json:"messages"`
	} `json:"body"`
}
