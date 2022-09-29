package module

import "time"

type BlockInfo struct {
	BlockID BlockID `json:"block_id"`
	Block   Block   `json:"block"`
}
type Parts struct {
	Total int    `json:"total"`
	Hash  string `json:"hash"`
}
type BlockID struct {
	Hash  string `json:"hash"`
	Parts Parts  `json:"parts"`
}
type Version struct {
	Block string `json:"block"`
	App   string `json:"app"`
}
type LastBlockID struct {
	Hash  string `json:"hash"`
	Parts Parts  `json:"parts"`
}
type Header struct {
	Version            Version     `json:"version"`
	ChainID            string      `json:"chain_id"`
	Height             string      `json:"height"`
	Time               time.Time   `json:"time"`
	LastBlockID        LastBlockID `json:"last_block_id"`
	LastCommitHash     string      `json:"last_commit_hash"`
	DataHash           string      `json:"data_hash"`
	ValidatorsHash     string      `json:"validators_hash"`
	NextValidatorsHash string      `json:"next_validators_hash"`
	ConsensusHash      string      `json:"consensus_hash"`
	AppHash            string      `json:"app_hash"`
	LastResultsHash    string      `json:"last_results_hash"`
	EvidenceHash       string      `json:"evidence_hash"`
	ProposerAddress    string      `json:"proposer_address"`
}
type Data struct {
	Txs []string `json:"txs"`
}
type Evidence struct {
	Evidence []interface{} `json:"evidence"`
}
type Signatures struct {
	BlockIDFlag      int       `json:"block_id_flag"`
	ValidatorAddress string    `json:"validator_address"`
	Timestamp        time.Time `json:"timestamp"`
	Signature        string    `json:"signature"`
}
type LastCommit struct {
	Height     string       `json:"height"`
	Round      int          `json:"round"`
	BlockID    BlockID      `json:"block_id"`
	Signatures []Signatures `json:"signatures"`
}
type Block struct {
	Header     Header     `json:"header"`
	Data       Data       `json:"data"`
	Evidence   Evidence   `json:"evidence"`
	LastCommit LastCommit `json:"last_commit"`
}
