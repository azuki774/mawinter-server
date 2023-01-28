package register

import (
	"context"
	"errors"
	"mawinter-server/internal/model"
	"os"

	"go.uber.org/zap"
)

var (
	ErrAlreadyRegisted = errors.New("already registed")
)

type DBRepository interface {
	InsertUniqueCatIDRecord(req model.Recordstruct) (res model.Recordstruct, err error)
	GetMonthlyFixDone(yyyymm string) (flag bool, err error)
	GetMonthlyFixBilling() (fixBills []model.MonthlyFixBilling, err error)
	InsertMonthlyFixBilling(yyyymm string, fixBills []model.MonthlyFixBilling) (err error)
}
type BillFetcher interface {
	FetchBills(ctx context.Context, yyyymm string) (ress []model.BillAPIResponse, err error)
}

type MailClient interface {
	Send(ctx context.Context, to string, title string, body string) (err error)
}
type RegisterService struct {
	Logger      *zap.Logger
	DB          DBRepository
	BillFetcher BillFetcher
	MailClient  MailClient
}

// InsertMonthlyFixBilling は 固定費を登録する
func (r *RegisterService) InsertMonthlyFixBilling(ctx context.Context, yyyymm string) (err error) {
	// すでに処理済なら skip
	done, err := r.DB.GetMonthlyFixDone(yyyymm)
	if err != nil {
		r.Logger.Error("failed to get done status", zap.Error(err))
		return err
	}

	if done {
		r.Logger.Warn("this month is processed")
		return model.ErrAlreadyRecorded
	}
	lg := r.Logger.With(zap.String("yyyymm", yyyymm))

	// Record テーブルに挿入するデータを取得
	fixBills, err := r.DB.GetMonthlyFixBilling()
	if err != nil {
		lg.Error("failed to get fix billing recoreds", zap.Error(err))
		return err
	}

	// Insert
	err = r.DB.InsertMonthlyFixBilling(yyyymm, fixBills)
	if err != nil {
		lg.Error("failed to insert fix billing records", zap.Error(err))
		return err
	}

	lg.Info("insert fix billing records to DB")

	// 環境変数 MAIL_TO に何か入ったときのみ通知メールを送信する。
	if os.Getenv("MAIL_TO") != "" {
		err = notifyMailInsertMonthlyFixBilling(ctx, r.MailClient, fixBills)
		if err != nil {
			// send error
			r.Logger.Error("notify mail send error", zap.Error(err))
			return err
		}
		r.Logger.Info("send notify mail", zap.String("mail_address", os.Getenv("MAIL_TO")))
	} else {
		r.Logger.Info("MAIL_TO is not set. sending a notify mail skipped.")
	}

	return nil
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

	// 環境変数 MAIL_TO に何か入ったときのみ通知メールを送信する。
	if os.Getenv("MAIL_TO") != "" {
		err = notifyMailMonthlyRegistBill(ctx, r.MailClient, ress)
		if err != nil {
			// send error
			r.Logger.Error("notify mail send error", zap.Error(err))
			return err
		}
		r.Logger.Info("send notify mail", zap.String("mail_address", os.Getenv("MAIL_TO")))
	} else {
		r.Logger.Info("MAIL_TO is not set. sending a notify mail skipped.")
	}

	return nil
}

func notifyMailInsertMonthlyFixBilling(ctx context.Context, MailClient MailClient, fbs []model.MonthlyFixBilling) (err error) {
	to := os.Getenv("MAIL_TO")
	title := "[Mawinter] 月次固定費の登録完了"
	body := model.NewMailMonthlyFixBilling(fbs)

	return MailClient.Send(ctx, to, title, body)
}

func notifyMailMonthlyRegistBill(ctx context.Context, MailClient MailClient, ress []model.BillAPIResponse) (err error) {
	to := os.Getenv("MAIL_TO")
	title := "[Mawinter] 月次光熱費等の登録完了"
	body := model.NewMailMonthlyRegistBill(ress)

	return MailClient.Send(ctx, to, title, body)
}
