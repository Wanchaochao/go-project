package cmd

import (
	"project/model"
	"project/pkg/mq"
	"project/script/internal/handler"
	"project/script/internal/service"

	"github.com/spf13/cobra"
)

var exampleMessageCmd = &cobra.Command{
	Use:   "example:message",
	Short: "nsq消费消息示例",
	Run: func(cmd *cobra.Command, args []string) {
		srv := service.NewService()
		h := handler.NewExampleMessage(srv)
		c := mq.NewNsqConsumer(cfg.Nsq.Consumer, model.TopicExample, "default", 4, h.Handle)
		Notify()
		c.Stop()
	},
}

func init() {
	rootCmd.AddCommand(exampleMessageCmd)
}
