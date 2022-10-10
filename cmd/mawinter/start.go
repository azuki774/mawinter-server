package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/factory"

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
		return start(&startOpt)
	},
}

func start(opts *StartOption) error {
	l, err := factory.NewLogger()
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		return err
	}

	api, err := factory.NewAPIService(opts.DBInfo.User, opts.DBInfo.Pass, opts.DBInfo.Host, opts.DBInfo.Port, opts.DBInfo.Name)
	if err != nil {
		return err
	}
	l.Info("loaded api service")

	srv, err := factory.NewServer(api)
	if err != nil {
		l.Error("failed to load server", zap.Error(err))
		return err
	}

	defer api.DBRepo.CloseDB()

	return srv.Start(context.Background())
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