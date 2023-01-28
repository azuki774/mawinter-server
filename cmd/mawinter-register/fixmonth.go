package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/client"
	"mawinter-server/internal/factory"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type fixMonthlyOption struct {
	Logger *zap.Logger
	DBInfo struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
}

var fixMonthlyOpt fixMonthlyOption

// fixMonthlyOptCmd represents the start command
var fixMonthlyOptCmd = &cobra.Command{
	Use:   "fixmonth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startFixMonthly()
	},
}

func startFixMonthly() (err error) {
	l, err := factory.NewLogger()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Sync()
	db, err := factory.NewDBRepository(fixMonthlyOpt.DBInfo.Host, fixMonthlyOpt.DBInfo.Port, fixMonthlyOpt.DBInfo.User, fixMonthlyOpt.DBInfo.Pass, fixMonthlyOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db.CloseDB()
	mc := factory.NewMailClient()
	ap := factory.NewRegisterService(l, db, &client.BillFetcher{}, mc)
	ctx := context.Background()
	return ap.InsertMonthlyFixBilling(ctx, time.Now().Local().Format("200601"))
}

func init() {
	rootCmd.AddCommand(fixMonthlyOptCmd)
	fixMonthlyOptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fixMonthlyOptCmd.Flags().StringVar(&fixMonthlyOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	fixMonthlyOptCmd.Flags().StringVar(&fixMonthlyOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	fixMonthlyOptCmd.Flags().StringVar(&fixMonthlyOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	fixMonthlyOptCmd.Flags().StringVar(&fixMonthlyOpt.DBInfo.User, "db-user", "root", "DB User")
	fixMonthlyOptCmd.Flags().StringVar(&fixMonthlyOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
}
