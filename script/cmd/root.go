package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"project/pkg/cache"
	"project/pkg/db"
	"project/pkg/logger"
	"syscall"
)

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Println("start ...")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Println("stop ...")
	},
}

var cfg struct {
	App struct {
		IsProd bool
		Logger string
	}
	Cdn    string
	Wechat struct {
		Appid  string
		Secret string
	}
	Robot struct {
		DingTalk   string
		WechatWork string
	}
	Mysql db.Mysql
	Redis cache.Redis
	Nsq   struct {
		Producer string
		Consumer string
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./script")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("viper.ReadInConfig error", err)
		}
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal("viper.Unmarshal error: ", err)
		}
		logger.SetOutput(cfg.App.Logger)
	})
}

// Notify 阻塞主进程，监听退出信息
func Notify() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
