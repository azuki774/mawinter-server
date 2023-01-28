package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mawinter-server/internal/model"
	"net/http"
	"time"
)

const httpTimeout = 10 * time.Second

type BillFetcher struct {
	BillEndpoint string // ex: http://localhost:8080/bill/
}

func (b *BillFetcher) FetchBills(ctx context.Context, yyyymm string) (ress []model.BillAPIResponse, err error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	req, err := http.NewRequest("GET", b.BillEndpoint+yyyymm, nil)
	if err != nil {
		return []model.BillAPIResponse{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return []model.BillAPIResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []model.BillAPIResponse{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []model.BillAPIResponse{}, err
	}

	err = json.Unmarshal(body, &ress)
	if err != nil {
		return []model.BillAPIResponse{}, err
	}

	return ress, nil
}
