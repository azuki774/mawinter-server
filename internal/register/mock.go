package register

import (
	"context"
	"mawinter-server/internal/model"
	"time"
)

type mockRepo struct {
	err error
}

type mockBillFetcher struct {
	err error
}

func (m *mockRepo) InsertUniqueCatIDRecord(req model.Recordstruct) (res model.Recordstruct, err error) {
	if m.err != nil {
		return model.Recordstruct{}, m.err
	}

	return model.Recordstruct{
		ID:         1,
		CategoryID: 210,
		Datetime:   time.Now(),
		From:       "bill-manager-api",
		Type:       "",
		Price:      1234,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *mockBillFetcher) FetchBills(ctx context.Context, yyyymm string) (ress []model.BillAPIResponse, err error) {
	if m.err != nil {
		return []model.BillAPIResponse{}, m.err
	}

	ress = []model.BillAPIResponse{
		{
			BillName: "water",
			Price:    12345,
		},
		{
			BillName: "gas",
			Price:    1234,
		},
		{
			BillName: "elect",
			Price:    123,
		},
	}
	return ress, nil
}
