package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/factory"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type billOption struct {
	Logger          *zap.Logger
	BillAPIEndpoint string
	DBInfo          struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
	Date string // YYYYMM
}

var billOpt billOption

// importCmd represents the start command
var billCmd = &cobra.Command{
	Use:   "bill",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start()
	},
}

func start() (err error) {
	l, err := factory.NewLogger()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Sync()
	db, err := factory.NewDBRepository(billOpt.DBInfo.Host, billOpt.DBInfo.Port, billOpt.DBInfo.User, billOpt.DBInfo.Pass, billOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db.CloseDB()
	fet := factory.NewFetcherBill(billOpt.BillAPIEndpoint)
	ap := factory.NewRegisterService(l, db, fet)
	ctx := context.Background()
	return ap.MonthlyRegistBill(ctx, billOpt.Date)
}

func init() {
	rootCmd.AddCommand(billCmd)
	billCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	billCmd.Flags().StringVar(&billOpt.BillAPIEndpoint, "bill-endpoint", "http://localhost/bill/", "bill-mangager API endpoint")
	billCmd.Flags().StringVar(&billOpt.Date, "date", time.Now().Local().Format("200601"), "YYYYMM")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	billCmd.Flags().StringVar(&billOpt.DBInfo.User, "db-user", "root", "DB User")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
}
