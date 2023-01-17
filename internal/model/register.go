package model

import (
	"fmt"
	"time"
)

type BillAPIResponse struct {
	BillName string `json:"bill_name"`
	Price    int    `json:"price"`
}

func (b *BillAPIResponse) NewRecordstruct() (req Recordstruct, err error) {
	req = Recordstruct{
		Datetime: time.Now().Local(),
		From:     "bill-manager-api",
		Price:    b.Price,
	}

	switch b.BillName {
	case "elect":
		req.CategoryID = 220
	case "gas":
		req.CategoryID = 221
	case "water":
		req.CategoryID = 222
	default:
		return Recordstruct{}, fmt.Errorf("unknown billname")
	}

	return req, nil
}
