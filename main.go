package main

import (
	"github.com/XM-GO/PandaKit/ginx"
	"github.com/XM-GO/PandaKit/logger"
	gStarter "github.com/XM-GO/PandaKit/starter"
	"github.com/spf13/cobra"
	"os"
	"pandax/apps/job/jobs"
	"pandax/pkg/config"
	"pandax/pkg/global"
	"pandax/pkg/initialize"
	"pandax/pkg/middleware"
	"pandax/pkg/starter"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "panda is the main component in the panda.",
	Short: `panda is go gin frame`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if configFile != "" {
			global.Conf = config.InitConfig(configFile)
			global.Log = logger.InitLog(global.Conf.Log.File.GetFilename(), global.Conf.Log.Level)

			dbGorm := gStarter.DbGorm{Type: global.Conf.Server.DbType}
			if dbGorm.Type == "mysql" {
				dbGorm.Dsn = global.Conf.Mysql.Dsn()
				dbGorm.MaxIdleConns = global.Conf.Mysql.MaxIdleConns
				dbGorm.MaxOpenConns = global.Conf.Mysql.MaxOpenConns
			} else {
				dbGorm.Dsn = global.Conf.Postgresql.PgDsn()
				dbGorm.MaxIdleConns = global.Conf.Postgresql.MaxIdleConns
				dbGorm.MaxOpenConns = global.Conf.Postgresql.MaxOpenConns
			}
			global.Db = dbGorm.GormInit()

			initialize.InitTable()
		} else {
			global.Log.Panic("请配置config")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ginx.UseAfterHandlerInterceptor(middleware.OperationHandler)
		// gin前置 函数
		ginx.UseBeforeHandlerInterceptor(middleware.PermissionHandler)
		// gin后置 函数
		ginx.UseAfterHandlerInterceptor(middleware.LogHandler)
		go func() {
			// 启动系统调度任务
			jobs.InitJob()
			jobs.Setup()
		}()

		starter.RunWebServer(initialize.InitRouter())
	},
}

func init() {
	rootCmd.Flags().StringVar(&configFile, "config", getEnvStr("PANDA_CONFIG", "./config.yml"), "panda config file path.")
}

func getEnvStr(env string, defaultValue string) string {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue
	}
	return v
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("panda root cmd execute: %s", err)
		os.Exit(1)
	}
}
