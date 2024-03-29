package module

import "time"

type TxResp struct {
	Tx         TxMsg `json:"tx"`
	TxResponse struct {
		Height    string        `json:"height"`
		Txhash    string        `json:"txhash"`
		Codespace string        `json:"codespace"`
		Code      int           `json:"code"`
		Data      string        `json:"data"`
		RawLog    string        `json:"raw_log"`
		Logs      []interface{} `json:"logs"`
		Info      string        `json:"info"`
		GasWanted string        `json:"gas_wanted"`
		GasUsed   string        `json:"gas_used"`
		Tx        struct {
			Type string `json:"@type"`
			Body struct {
				Messages []struct {
					Type   string `json:"@type"`
					Sender string `json:"sender"`
					Routes []struct {
						PoolID        string `json:"poolId"`
						TokenOutDenom string `json:"tokenOutDenom"`
					} `json:"routes"`
					TokenIn struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"tokenIn"`
					TokenOutMinAmount string `json:"tokenOutMinAmount"`
				} `json:"messages"`
				Memo                        string        `json:"memo"`
				TimeoutHeight               string        `json:"timeout_height"`
				ExtensionOptions            []interface{} `json:"extension_options"`
				NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						Type string `json:"@type"`
						Key  string `json:"key"`
					} `json:"public_key"`
					ModeInfo struct {
						Single struct {
							Mode string `json:"mode"`
						} `json:"single"`
					} `json:"mode_info"`
					Sequence string `json:"sequence"`
				} `json:"signer_infos"`
				Fee struct {
					Amount   []interface{} `json:"amount"`
					GasLimit string        `json:"gas_limit"`
					Payer    string        `json:"payer"`
					Granter  string        `json:"granter"`
				} `json:"fee"`
			} `json:"auth_info"`
			Signatures []string `json:"signatures"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string      `json:"key"`
				Value interface{} `json:"value"`
				Index bool        `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"tx_response"`
}

type TxMsg struct {
	Body struct {
		Messages []struct {
			Type   string `json:"@type"`
			Sender string `json:"sender"`
			Routes []struct {
				PoolID        string `json:"poolId"`
				TokenOutDenom string `json:"tokenOutDenom"`
			} `json:"routes"`
			TokenIn struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"tokenIn"`
			TokenOutMinAmount string `json:"tokenOutMinAmount"`
		} `json:"messages"`
		Memo                        string        `json:"memo"`
		TimeoutHeight               string        `json:"timeout_height"`
		ExtensionOptions            []interface{} `json:"extension_options"`
		NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
	} `json:"body"`
	AuthInfo struct {
		SignerInfos []struct {
			PublicKey struct {
				Type string `json:"@type"`
				Key  string `json:"key"`
			} `json:"public_key"`
			ModeInfo struct {
				Single struct {
					Mode string `json:"mode"`
				} `json:"single"`
			} `json:"mode_info"`
			Sequence string `json:"sequence"`
		} `json:"signer_infos"`
		Fee struct {
			Amount   []interface{} `json:"amount"`
			GasLimit string        `json:"gas_limit"`
			Payer    string        `json:"payer"`
			Granter  string        `json:"granter"`
		} `json:"fee"`
	} `json:"auth_info"`
	Signatures []string `json:"signatures"`
}
