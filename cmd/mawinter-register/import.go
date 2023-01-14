package main

import (
	"fmt"
	"mawinter-server/internal/factory"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type billOption struct {
	Logger *zap.Logger
	DBInfo struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
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
	return nil
}

func init() {
	rootCmd.AddCommand(billCmd)
	billCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	billCmd.Flags().StringVar(&billOpt.DBInfo.User, "db-user", "root", "DB User")
	billCmd.Flags().StringVar(&billOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
}
