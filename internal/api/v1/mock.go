package api

import (
	"mawinter-server/internal/model"
)

type mockRepo struct {
	errCreateRecordTable       error
	errGetCategoryInfo         error
	errSumPriceForEachCatID    error
	RecordYYYYMMNum            int
	sumPriceForEachCatIDOffset int
}

func (m *mockRepo) CreateRecordTable(yyyymm string) (err error) {
	if m.errCreateRecordTable != nil {
		return m.errCreateRecordTable
	}
	return nil
}

func (m *mockRepo) InsertRecord(req model.Recordstruct) (res model.Recordstruct, err error) {
	return model.Recordstruct{}, nil
}

func (m *mockRepo) GetCategoryInfo() (info []model.Category, err error) {
	if m.errGetCategoryInfo != nil {
		return []model.Category{}, m.errGetCategoryInfo
	}

	info = []model.Category{
		{
			ID:         1,
			CategoryID: 100,
			Name:       "カテゴリ1",
		},
		{
			ID:         2,
			CategoryID: 200,
			Name:       "カテゴリ2",
		},
		{
			ID:         3,
			CategoryID: 300,
			Name:       "カテゴリ3",
		},
	}
	return info, nil
}

func (m *mockRepo) SumPriceForEachCatID(yyyymm string) (sum []model.SumPriceCategoryID, err error) {
	if m.errSumPriceForEachCatID != nil {
		return []model.SumPriceCategoryID{}, m.errSumPriceForEachCatID
	}
	if m.RecordYYYYMMNum <= m.sumPriceForEachCatIDOffset {
		// 0 records
		return []model.SumPriceCategoryID{}, nil
	}

	sum = []model.SumPriceCategoryID{
		{
			CategoryID: 100,
			Count:      5,
			Sum:        10,
		},
		{
			CategoryID: 200,
			Count:      5,
			Sum:        100,
		},
		// {
		// 	CategoryID: 300,
		// 	Count:      5,
		// 	Sum:        ,
		// },
	}
	m.sumPriceForEachCatIDOffset = m.sumPriceForEachCatIDOffset + 1
	return sum, nil
}
