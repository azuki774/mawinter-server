package main

import (
	"context"
	"fmt"
	"mawinter-server/internal/factory"
	"mawinter-server/internal/server"
	"os"

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

	// Override with environment variables for security
	dbHost := getEnvOrDefault("DB_HOST", startOpt.DBInfo.Host)
	dbPort := getEnvOrDefault("DB_PORT", startOpt.DBInfo.Port)
	dbName := getEnvOrDefault("DB_NAME", startOpt.DBInfo.Name)
	dbUser := getEnvOrDefault("DB_USER", startOpt.DBInfo.User)
	dbPass := getEnvOrDefault("DB_PASS", startOpt.DBInfo.Pass)

	db2, err := factory.NewDBRepositoryV2(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db2.CloseDB()

	ap2 := factory.NewServiceV2(l, db2)
	srv := factory.NewServer(l, ap2)
	ctx := context.Background()

	l.Info("binary info", zap.String("version", version), zap.String("revision", revision), zap.String("build", build))
	server.Version = version
	server.Revision = revision
	server.Build = build

	return srv.Start(ctx)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host (can be overridden by DB_HOST env var)")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Port, "db-port", "3306", "DB Port (can be overridden by DB_PORT env var)")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Name, "db-name", "mawinter", "DB Name (can be overridden by DB_NAME env var)")
	startCmd.Flags().StringVar(&startOpt.DBInfo.User, "db-user", "root", "DB User (can be overridden by DB_USER env var)")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Pass, "db-pass", "password", "DB Pass (can be overridden by DB_PASS env var - DEPRECATED: use env var for security)")
}
