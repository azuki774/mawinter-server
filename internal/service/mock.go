package service

import (
	"context"
	"mawinter-server/internal/openapi"
	"time"
)

type mockAp struct{}
type mockMailClient struct{}

func (m *mockAp) GetYYYYMMRecords(ctx context.Context, yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error) {
	return []openapi.Record{
		{
			CategoryId: 100,
			Datetime:   time.Date(2000, 1, 23, 12, 0, 0, 0, jst),
			From:       "from1",
			Price:      1234,
		},
		{
			CategoryId: 100,
			Datetime:   time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
			From:       "from1",
			Price:      1234,
		},
		{
			CategoryId: 200,
			Datetime:   time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
			From:       "from1",
			Price:      1234,
		},
	}, nil
}

func (m *mockMailClient) Send(ctx context.Context, to string, title string, body string) (err error) {
	return nil
}
