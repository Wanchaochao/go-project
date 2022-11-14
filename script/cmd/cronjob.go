package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
	"project/pkg/logger"
	"project/pkg/wechat"
	"project/script/internal/handler"
	"project/script/internal/service"
	"time"
)

/*
spec格式:
	"Seconds Minutes Hours Day-of-month Month Day-of-week"

字段说明：
Field name   | Allowed values  | Allowed special characters
----------   | --------------  | --------------------------
Seconds      | 0-59            | * / , -
Minutes      | 0-59            | * / , -
Hours        | 0-23            | * / , -
Day-of-month | 1-31            | * / , - ?
Month        | 1-12 or JAN-DEC | * / , -
Day-of-week  | 0-6 or SUN-SAT  | * / , - ?

匹配符：
* 匹配任意值
/ 范围的增量，如在Minutes字段用"3-59/15"表示一小时中第3分钟开始到第59分钟，每隔15分钟
, 用于分隔列表中的项目，如在Day-of-week字段用"MON,WED,FRI"表示星期一三五
- 用于定义范围，如在Hours字段用"9-18"表示9:00~18:00之间的每小时
? 可用于代替*，将Day-of-month或Day-of-week留空

预定义简写：
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

默认标准parser只有5位：分时日月周，如需秒开始的6位，创建*cron.Cron需使用cron.New(cron.WithSeconds())
*/

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "定时任务",
	Long:  "精确定时执行的任务",
	Run: func(cmd *cobra.Command, args []string) {
		srv := service.NewService(service.NewMysql(&cfg.Mysql), service.NewRedis(&cfg.Redis))
		h := handler.NewCronjob(
			srv,
			wechat.NewServerAPI(logger.NewHttpClient(30*time.Second), srv.GetWechatToken),
			cfg.Robot.DingTalk,
			cfg.Robot.WechatWork,
		)

		c := cron.New()
		var err error

		_, err = c.AddFunc("1 0 * * *", h.LoadWechatAnalysis) // 每天0点1分拉取昨日微信小程序访问数据
		if err != nil {
			log.Fatal(err)
		}

		c.Start()
		Notify()
		ctx := c.Stop()
		<-ctx.Done()
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}
