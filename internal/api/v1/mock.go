package v1

import (
	"mawinter-server/internal/model"
)

type mockRepo struct {
	errGetCategoryInfo              error
	errGetCategoryMonthSummary      error
	CategoryMonthSummaryMonthTotal  []int
	CategoryMonthSummaryMonthCount  []int
	categoryMonthSummaryMonthOffset int
}

func (m *mockRepo) InsertRecord(req *model.RecordRequest) (res *model.Record_YYYYMM, err error) {
	return nil, nil
}

func (m *mockRepo) GetCategoryInfo(categoryID int) (info *model.Category, err error) {
	if m.errGetCategoryInfo != nil {
		return nil, m.errGetCategoryInfo
	}

	info = &model.Category{
		ID:         1,
		CategoryID: int64(categoryID),
		Name:       "カテゴリ1",
	}
	return info, nil
}

func (m *mockRepo) GetCategoryMonthSummary(categoryID int, yyyymm string) (sum *model.CategoryMonthSummary, err error) {
	if m.errGetCategoryMonthSummary != nil {
		return nil, m.errGetCategoryMonthSummary
	}

	l := len(m.CategoryMonthSummaryMonthTotal)
	if l <= m.categoryMonthSummaryMonthOffset {
		// Table not found
		return nil, model.ErrTableNotFound
	}
	sum = &model.CategoryMonthSummary{
		CategoryID:   categoryID,
		CategoryName: "カテゴリ1",
		Count:        m.CategoryMonthSummaryMonthCount[m.categoryMonthSummaryMonthOffset],
		Total:        m.CategoryMonthSummaryMonthTotal[m.categoryMonthSummaryMonthOffset],
	}
	m.categoryMonthSummaryMonthOffset = m.categoryMonthSummaryMonthOffset + 1
	return sum, nil
}
