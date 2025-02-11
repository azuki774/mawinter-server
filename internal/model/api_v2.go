package model

// CategoryMidMonthSummary は openapi.CategoryYearSummary を作るための中間構造体
type CategoryMidMonthSummary struct {
	CategoryId int `json:"category_id"`
	Count      int `json:"count"`
	Price      int `json:"price"`
}

// GetRecordOption は GetRecord する際にAPI/DB共通で参照するためのオプション
type GetRecordOption struct {
	Num    int
	Offset int
}
