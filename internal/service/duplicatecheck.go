package service

import (
	"context"
	"mawinter-server/internal/openapi"
	"strings"
	"time"

	"go.uber.org/zap"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

type APIServiceDup interface {
	GetYYYYMMRecords(ctx context.Context, yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error)
}

type DuplicateCheckService struct {
	Logger *zap.Logger
	Ap     APIServiceDup
}

func judgeDuplicateRecords(d1 openapi.Record, d2 openapi.Record) bool {
	// CategoryID と Price と Datetime の日付までが一致していれば True（重複可能性あり）とする
	if d1.CategoryId != d2.CategoryId {
		return false
	}

	if d1.CategoryId != d2.CategoryId {
		return false
	}

	if d1.Datetime.Format("20060102") != d2.Datetime.Format("20060102") {
		return false
	}

	return true
}

func (d *DuplicateCheckService) DuplicateCheck(ctx context.Context, yyyymm string) (err error) {
	d.Logger.Info("DuplicateCheck start")
	var dupInt = 0

	recs, err := d.Ap.GetYYYYMMRecords(ctx, yyyymm, openapi.GetV2RecordYyyymmParams{})
	if err != nil {
		return err
	}

	var targets []openapi.Record // 重複判定する対象
	for _, r := range recs {
		if strings.Contains(r.Type, "D") {
			// Type に 'D' が入っているレコードは重複判定の対象外
			continue
		}
		targets = append(targets, r)
	}

	// targets 内全体に重複の判定をかける
	for i, u := range targets {
		for j, v := range targets {
			// ダブルカウントを防ぐため、i < j とする
			if i >= j {
				continue
			}

			if judgeDuplicateRecords(u, v) {
				// TODO: duplicate notification
				d.Logger.Info("detect duplicate data")
				dupInt++
			}
		}
	}

	d.Logger.Info("DuplicateCheck complete", zap.Int("rec_num", len(recs)), zap.Int("target_num", len(targets)), zap.Int("duplicate_num", dupInt))
	return nil
}
