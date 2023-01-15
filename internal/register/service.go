package register

import (
	"context"
	"errors"
	"mawinter-server/internal/model"

	"go.uber.org/zap"
)

var (
	ErrAlreadyRegisted = errors.New("already registed")
)

type DBRepository interface {
	InsertUniqueCatIDRecord(req model.Recordstruct) (res model.Recordstruct, err error)
}
type BillFetcher interface {
	FetchBills(ctx context.Context, yyyymm string) (ress []model.BillAPIResponse, err error)
}
type RegisterService struct {
	Logger      *zap.Logger
	DB          DBRepository
	BillFetcher BillFetcher
}

// MonthlyRegistBill fetch from bill-mangager-api and add to mawinter DB.
func (r *RegisterService) MonthlyRegistBill(ctx context.Context, yyyymm string) (err error) {
	if err := model.ValidYYYYMM(yyyymm); err != nil {
		r.Logger.Error("invalid args", zap.String("yyyymm", yyyymm), zap.Error(err))
		return err
	}

	r.Logger.Info("get bill data from API", zap.String("yyyymm", yyyymm))
	ress, err := r.BillFetcher.FetchBills(ctx, yyyymm)
	if err != nil {
		r.Logger.Error("failed to get bill records", zap.Error(err))
		return err
	}

	r.Logger.Info("response from bili fetcher", zap.String("yyyymm", yyyymm))

	for _, res := range ress {
		req, err := res.NewRecordstruct()
		if err != nil {
			r.Logger.Error("failed to new bill records", zap.Error(err))
			return err
		}

		_, err = r.DB.InsertUniqueCatIDRecord(req)
		if errors.Is(err, ErrAlreadyRegisted) {
			// already category_id registed
			r.Logger.Warn("this category_id is already registed", zap.Int("category_id", req.CategoryID), zap.Error(err))
			continue
		}
		if err != nil {
			r.Logger.Error("failed to insert records", zap.Error(err))
			return err
		}

		r.Logger.Info("insert bill record", zap.String("billname", res.BillName))
	}

	r.Logger.Info("insert bill records sucessfully")
	return nil
}
