package main

import (
	"fmt"
	"mawinter-server/internal/factory"
	"mawinter-server/internal/server"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DuplicateCheckOption struct {
	Logger *zap.Logger
	DBInfo struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
}

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
	// duplicateCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Run() (err error) {
	l, err := factory.NewLogger()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Sync()

	db, err := factory.NewDBRepositoryV2(startOpt.DBInfo.Host, startOpt.DBInfo.Port, startOpt.DBInfo.User, startOpt.DBInfo.Pass, startOpt.DBInfo.Name)
	if err != nil {
		l.Error("failed to connect DB", zap.Error(err))
		return err
	}
	defer db.CloseDB()

	l.Info("binary info", zap.String("version", version), zap.String("revision", revision), zap.String("build", build))
	server.Version = version
	server.Revision = revision
	server.Build = build

	return nil
}
