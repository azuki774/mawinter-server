package main

import (
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

	return srv.Start()

	// l, err := logger.NewSugarLogger()
	// defer l.Sync()
	// if err != nil {
	// 	fmt.Printf("logger failed")
	// 	os.Exit(1)
	// }

	// l.Infow("Program Start")

	// if os.Getenv("BASIC_AUTH_PASSWORD") == "" {
	// 	l.Warnw("No Basic Auth password set")
	// } else {
	// 	l.Infow("Basic Authentication info", "username", os.Getenv("BASIC_AUTH_USERNAME"), "password", os.Getenv("BASIC_AUTH_PASSWORD"))
	// }

	// l.Infow("database info", "name", os.Getenv("MYSQL_DATABASE"))

	// var DBSleepTime time.Duration
	// if os.Getenv("DB_WAITTIME") == "" {
	// 	DBSleepTime = 0
	// } else {
	// 	t, err := strconv.Atoi(os.Getenv("DB_WAITTIME"))
	// 	if err != nil {
	// 		l.Errorw("DB_WAITTIME is invalid")
	// 		DBSleepTime = 0
	// 	} else {
	// 		l.Infof("DB_WAITTIME is %d s", t)
	// 		DBSleepTime = time.Duration(t)
	// 	}
	// }

	// time.Sleep(time.Second * DBSleepTime)
	// gormdb, err := repository.DBConnect("root", "password", "mawinter-db", os.Getenv("MYSQL_DATABASE"))
	// if err != nil {
	// 	l.Errorw(err.Error())
	// 	os.Exit(1)
	// }
	// sqlDB, err := gormdb.DB()
	// if err != nil {
	// 	l.Error(err.Error())
	// 	os.Exit(1)
	// }
	// defer sqlDB.Close()

	// dbR := repository.NewDBRepository(gormdb)
	// as := api.NewAPIService(dbR, l)
	// server.Start(as, l)
	// l.Info("Program End")

}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Host, "db-host", "mawinter-db", "DB Host")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Port, "db-port", "3306", "DB Port")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Name, "db-name", "mawinter", "DB Name")
	startCmd.Flags().StringVar(&startOpt.DBInfo.User, "db-user", "root", "DB User")
	startCmd.Flags().StringVar(&startOpt.DBInfo.Pass, "db-pass", "password", "DB Pass")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
