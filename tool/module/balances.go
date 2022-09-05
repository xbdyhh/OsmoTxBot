package module

import (
	"github.com/graph-gophers/graphql-go/errors"
	"strconv"
)

type Balances struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

func (pb *Balances) Get(denom string) (int, error) {
	for _, v := range pb.Balances {
		if v.Denom == denom {
			val, err := strconv.Atoi(v.Amount)
			if err != nil {
				return 0, err
			}
			return val, nil
		}
	}
	return 0, errors.Errorf("key not found!")
}
