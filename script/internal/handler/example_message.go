package handler

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"project/model"
	"project/pkg/logger"
	"project/pkg/util/types"
	"project/script/internal/service"
)

type ExampleMessage struct {
	service *service.Service
}

func NewExampleMessage(srv *service.Service) *ExampleMessage {
	return &ExampleMessage{
		service: srv,
	}
}

func (h *ExampleMessage) Handle(msg *nsq.Message) error {
	reqid := string(msg.ID[:])
	var data model.MsgExample
	_ = json.Unmarshal(msg.Body, &data)
	_, l := logger.NewCtxLog(reqid, "Message", "Handle", types.Int2Str(msg.Timestamp))
	l.Info("msg.body", msg.Body, &data)
	return nil
}
