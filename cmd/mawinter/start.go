package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/factory"
	"mawinter-server/internal/server"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type StartOption struct {
	Logger *zap.Logger
	DBInfo struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
	BasicAuth struct {
		User string
		Pass string
	}
}

var (
	version  string
	revision string
	build    string
)

var startOpt StartOption

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
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

	db1, err := factory.NewDBRepositoryV1(startOpt.DBInfo.Host, startOpt.DBInfo.Port, startOpt.DBInfo.User, startOpt.DBInfo.Pass, startOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}

	db2, err := factory.NewDBRepositoryV2(startOpt.DBInfo.Host, startOpt.DBInfo.Port, startOpt.DBInfo.User, startOpt.DBInfo.Pass, startOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db2.CloseDB()

	ap1 := factory.NewServiceV1(l, db1)
	ap2 := factory.NewServiceV2(l, db2)
	srv := factory.NewServer(l, ap1, ap2)
	ctx := context.Background()

	l.Info("binary info", zap.String("version", version), zap.String("revision", revision), zap.String("build", build))
	server.Version = version
	server.Revision = revision
	server.Build = build

	return srv.Start(ctx)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	startCmd.Flags().StringVar(&startOpt.DBInfo.User, "db-user", "root", "DB User")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
}
