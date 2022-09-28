package module

type BlockInfo struct {
	Height string `json:"height"`
	Result Result `json:"result"`
}
type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Validators struct {
	Address          string `json:"address"`
	PubKey           PubKey `json:"pub_key"`
	ProposerPriority string `json:"proposer_priority"`
	VotingPower      string `json:"voting_power"`
}
type Result struct {
	BlockHeight string       `json:"block_height"`
	Validators  []Validators `json:"validators"`
	Total       string       `json:"total"`
}
