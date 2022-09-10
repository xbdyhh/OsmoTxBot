package module

type SimulateBody struct {
	Txmsg   TxMsg  `json:"tx"`
	TxBytes string `json:"tx_bytes"`
}
type SimulateResponse struct {
	GasInfo struct {
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
	} `json:"gas_info"`
	Result struct {
		Data   string `json:"data"`
		Log    string `json:"log"`
		Events []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"result"`
}
