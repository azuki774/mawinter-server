package register

import (
	"context"
	"mawinter-server/internal/model"
	"time"
)

type mockRepo struct {
	err            error
	monthlyFixDone bool
	errGetMonthly  error
}

type mockMailClient struct {
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

func (m *mockRepo) GetMonthlyFixDone(yyyymm string) (flag bool, err error) {
	if m.errGetMonthly != nil {
		return false, m.errGetMonthly
	}
	return m.monthlyFixDone, nil
}

func (m *mockRepo) GetMonthlyFixBilling() (fixBills []model.MonthlyFixBilling, err error) {
	if m.errGetMonthly != nil {
		return []model.MonthlyFixBilling{}, m.errGetMonthly
	}
	return []model.MonthlyFixBilling{
		{
			CategoryID: 100,
			Day:        2,
			Type:       "type1",
			Memo:       "memo1",
		},
		{
			CategoryID: 101,
			Day:        4,
			Type:       "type2",
			Memo:       "memo2",
		},
	}, nil
}
func (m *mockRepo) InsertMonthlyFixBilling(yyyymm string, fixBills []model.MonthlyFixBilling) (err error) {
	return m.errGetMonthly
}

func (m *mockMailClient) Send(ctx context.Context, to string, title string, body string) (err error) {
	if m.err != nil {
		return m.err
	}
	return nil
}
