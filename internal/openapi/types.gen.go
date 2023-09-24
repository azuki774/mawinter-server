// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package openapi

import (
	"time"
)

// CategoryYearSummary defines model for category_year_summary.
type CategoryYearSummary struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
	Price        []int  `json:"price"`
	Total        int    `json:"total"`
}

// ConfirmInfo defines model for confirm_info.
type ConfirmInfo struct {
	ConfirmDatetime *time.Time `json:"confirm_datetime,omitempty"`
	Status          *bool      `json:"status,omitempty"`
	Yyyymm          *string    `json:"yyyymm,omitempty"`
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

// PostV2RecordFixmonthParams defines parameters for PostV2RecordFixmonth.
type PostV2RecordFixmonthParams struct {
	Yyyymm *int `form:"yyyymm,omitempty" json:"yyyymm,omitempty"`
}

// GetV2RecordYyyymmParams defines parameters for GetV2RecordYyyymm.
type GetV2RecordYyyymmParams struct {
	CategoryId *int `form:"category_id,omitempty" json:"category_id,omitempty"`
}

// PutV2TableYyyymmConfirmJSONBody defines parameters for PutV2TableYyyymmConfirm.
type PutV2TableYyyymmConfirmJSONBody struct {
	Status *bool `json:"status,omitempty"`
}

// GetV2RecordYyyymmRecentParams defines parameters for GetV2RecordYyyymmRecent.
type GetV2RecordYyyymmRecentParams struct {
	// Num max record number
	Num *int `form:"num,omitempty" json:"num,omitempty"`
}

// GetVersionJSONBody defines parameters for GetVersion.
type GetVersionJSONBody struct {
	Build    *string `json:"build,omitempty"`
	Revision *string `json:"revision,omitempty"`
	Version  *string `json:"version,omitempty"`
}

// PostV2RecordJSONRequestBody defines body for PostV2Record for application/json ContentType.
type PostV2RecordJSONRequestBody = ReqRecord

// PutV2TableYyyymmConfirmJSONRequestBody defines body for PutV2TableYyyymmConfirm for application/json ContentType.
type PutV2TableYyyymmConfirmJSONRequestBody PutV2TableYyyymmConfirmJSONBody

// GetVersionJSONRequestBody defines body for GetVersion for application/json ContentType.
type GetVersionJSONRequestBody GetVersionJSONBody
