package model

import (
	"fmt"
	"mawinter-server/internal/openapi"
	"strconv"
	"time"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

type BillAPIResponse struct {
	BillName string `json:"bill_name"`
	Price    int    `json:"price"`
}

func (b *BillAPIResponse) NewRecordstruct() (req openapi.Record, err error) {
	req = openapi.Record{
		Datetime: time.Now().Local(),
		From:     "bill-manager-api",
		Price:    b.Price,
	}

	switch b.BillName {
	case "elect":
		req.CategoryId = 220
	case "gas":
		req.CategoryId = 221
	case "water":
		req.CategoryId = 222
	default:
		return openapi.Record{}, fmt.Errorf("unknown billname")
	}

	return req, nil
}

func NewMailMonthlyFixBilling(recs []openapi.Record) (text string) {
	for _, rec := range recs {
		text += fmt.Sprintf("%v,%v,%v,%v\n", rec.CategoryId, rec.Price, rec.Type, rec.Memo)
	}
	return text
}

func NewMailMonthlyRegistBill(ress []BillAPIResponse) (text string) {
	for _, res := range ress {
		text += fmt.Sprintf("%s,%d\n", res.BillName, res.Price)
	}
	return text
}

type MonthlyFixBilling struct {
	CategoryID int
	Day        int
	Price      int
	Type       string
	Memo       string
}

func (m *MonthlyFixBilling) ConvAddDBModel(yyyymm string) (Record, error) {
	yyyynum, err := strconv.Atoi(yyyymm[0:4])
	if err != nil {
		return Record{}, err
	}

	mmnum, err := strconv.Atoi(yyyymm[5:6])
	if err != nil {
		return Record{}, err
	}

	return Record{
		CategoryID: int64(m.CategoryID),
		Datetime:   time.Date(yyyynum, time.Month(mmnum), m.Day, 0, 0, 0, 0, jst),
		From:       "fixmonth", // 固定値
		Price:      int64(m.Price),
		Type:       m.Type,
		Memo:       m.Memo,
	}, nil
}
