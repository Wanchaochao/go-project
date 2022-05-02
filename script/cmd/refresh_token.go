package cmd

import (
	"github.com/spf13/cobra"
	"project/pkg/logger"
	"project/pkg/wechat"
	"project/script/internal/handler"
	"project/script/internal/service"
	"time"
)

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh:token",
	Short: "刷新小程序AccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		srv := service.NewService(service.NewRedis(&cfg.Redis))
		h := handler.NewRefreshToken(
			srv,
			wechat.NewBasicAPI(cfg.Wechat.Appid, cfg.Wechat.Secret, logger.NewHttpClient(30*time.Second)),
		)
		stop := make(chan struct{})
		done := make(chan struct{})
		go func() {
			tk := time.Tick(2 * time.Minute)
			for {
				select {
				case <-stop:
					close(done)
					return
				case <-tk:
					h.WechatServerToken()
				}
			}
		}()
		Notify()
		close(stop)
		<-done
	},
}

func init() {
	rootCmd.AddCommand(refreshTokenCmd)
}
