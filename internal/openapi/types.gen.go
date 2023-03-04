// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package openapi

import (
	"time"
)

// CategoryYearSummary defines model for category_year_summary.
type CategoryYearSummary struct {
	CategoryId   *int       `json:"category_id,omitempty"`
	CategoryName *string    `json:"category_name,omitempty"`
	Price        *[]float32 `json:"price,omitempty"`
	Total        *int       `json:"total,omitempty"`
}

// Record defines model for record.
type Record struct {
	CategoryId   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Datetime     time.Time `json:"datetime"`
	From         string    `json:"from"`
	Id           int       `json:"id"`
	Memo         string    `json:"memo"`
	Price        int       `json:"price"`
	Type         string    `json:"type"`
}

// ReqRecord defines model for req_record.
type ReqRecord struct {
	CategoryId int     `json:"category_id"`
	Datetime   *string `json:"datetime,omitempty"`
	From       *string `json:"from,omitempty"`
	Memo       *string `json:"memo,omitempty"`
	Price      int     `json:"price"`
	Type       *string `json:"type,omitempty"`
}

// PostRecordJSONBody defines parameters for PostRecord.
type PostRecordJSONBody = map[string]interface{}

// PostV2RecordFixmonthParams defines parameters for PostV2RecordFixmonth.
type PostV2RecordFixmonthParams struct {
	Yyyymm *string `form:"yyyymm,omitempty" json:"yyyymm,omitempty"`
}

// GetV2RecordYyyymmParams defines parameters for GetV2RecordYyyymm.
type GetV2RecordYyyymmParams struct {
	From *string `form:"from,omitempty" json:"from,omitempty"`
}

// PostRecordJSONRequestBody defines body for PostRecord for application/json ContentType.
type PostRecordJSONRequestBody = PostRecordJSONBody

// PostV2RecordJSONRequestBody defines body for PostV2Record for application/json ContentType.
type PostV2RecordJSONRequestBody = ReqRecord
