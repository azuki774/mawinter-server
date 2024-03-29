package register

import (
	"context"
	"errors"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"os"

	"go.uber.org/zap"
)

var (
	ErrAlreadyRegisted = errors.New("already registed")
)

type DBRepository interface {
	InsertUniqueCatIDRecord(req openapi.Record) (res openapi.Record, err error)
	GetMonthlyFixDone(yyyymm string) (flag bool, err error)
	GetMonthlyFixBilling() (fixBills []model.MonthlyFixBilling, err error)
	InsertMonthlyFixBilling(yyyymm string) (recs []openapi.Record, err error)
}

type MailClient interface {
	Send(ctx context.Context, to string, title string, body string) (err error)
}
type RegisterService struct {
	Logger     *zap.Logger
	DB         DBRepository
	MailClient MailClient
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

	// Insert
	recs, err := r.DB.InsertMonthlyFixBilling(yyyymm)
	if err != nil {
		lg.Error("failed to insert fix billing records", zap.Error(err))
		return err
	}

	lg.Info("insert fix billing records to DB")

	// 環境変数 MAIL_TO に何か入ったときのみ通知メールを送信する。
	if os.Getenv("MAIL_TO") != "" {
		err = notifyMailInsertMonthlyFixBilling(ctx, r.MailClient, recs)
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

func notifyMailInsertMonthlyFixBilling(ctx context.Context, MailClient MailClient, recs []openapi.Record) (err error) {
	to := os.Getenv("MAIL_TO")
	title := "[Mawinter] 月次固定費の登録完了"
	body := model.NewMailMonthlyFixBilling(recs)

	return MailClient.Send(ctx, to, title, body)
}
