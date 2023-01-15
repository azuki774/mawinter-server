package client

import (
	"context"
	"mawinter-server/internal/model"
)

type BillFetcher struct {
	Host string
	Port string
}

func (b *BillFetcher) FetchBills(ctx context.Context, yyyymm string) (ress []model.BillAPIResponse, err error) {
	// TODO
	return ress, nil
}
