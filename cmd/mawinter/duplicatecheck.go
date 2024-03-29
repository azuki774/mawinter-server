package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/factory"
	"mawinter-server/internal/server"
	"mawinter-server/internal/timeutil"
	"time"

	"github.com/spf13/cobra"
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

type DuplicateCheckOption struct {
	Logger *zap.Logger
	DBInfo struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
	Lastmonth bool // if true, process last month table
}

var duplicateCheckOpt DuplicateCheckOption

// duplicateCheckCmd represents the duplicateCheck command
var duplicateCheckCmd = &cobra.Command{
	Use:   "duplicateCheck",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("duplicateCheck called")
		Run()
	},
}

func init() {
	rootCmd.AddCommand(duplicateCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// duplicateCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	duplicateCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	duplicateCheckCmd.Flags().StringVar(&duplicateCheckOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	duplicateCheckCmd.Flags().StringVar(&duplicateCheckOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	duplicateCheckCmd.Flags().StringVar(&duplicateCheckOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	duplicateCheckCmd.Flags().StringVar(&duplicateCheckOpt.DBInfo.User, "db-user", "root", "DB User")
	duplicateCheckCmd.Flags().StringVar(&duplicateCheckOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
	duplicateCheckCmd.Flags().BoolVar(&duplicateCheckOpt.Lastmonth, "last-month", false, "if true, process last month table")
}

func Run() (err error) {
	var YYYYMM string // roc month table name (YYYYMM
	thisMonth := time.Date(timeutil.NowFunc().Year(), timeutil.NowFunc().Month(), 1, 0, 0, 0, 0, jst)

	if !duplicateCheckOpt.Lastmonth {
		YYYYMM = timeutil.NowFunc().Format("200601")
	} else {
		// last month
		YYYYMM = thisMonth.AddDate(0, -1, 0).Format("200601")
	}

	l, err := factory.NewLogger()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Sync()
	l.Info("binary info", zap.String("version", version), zap.String("revision", revision), zap.String("build", build))
	server.Version = version
	server.Revision = revision
	server.Build = build

	db, err := factory.NewDBRepositoryV2(startOpt.DBInfo.Host, startOpt.DBInfo.Port, startOpt.DBInfo.User, startOpt.DBInfo.Pass, startOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db.CloseDB()

	ctx := context.Background()
	ap := factory.NewServiceV2(l, db)
	svc := factory.NewDuplicateCheckService(l, ap)

	l.Info("proc month table name", zap.String("YYYYMM", YYYYMM))

	return svc.DuplicateCheck(ctx, YYYYMM)
}
