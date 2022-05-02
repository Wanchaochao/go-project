package handler

import (
	"project/pkg/logger"
	"project/pkg/util/random"
	"project/pkg/wechat"
	"project/script/internal/service"
	"time"
)

type RefreshToken struct {
	service *service.Service
	wechat  wechat.BasicAPI
}

func NewRefreshToken(srv *service.Service, api wechat.BasicAPI) *RefreshToken {
	return &RefreshToken{
		service: srv,
		wechat:  api,
	}

}

func (s *RefreshToken) WechatServerToken() {
	ctx, l := logger.NewCtxLog(random.UUID(), "RefreshToken", "WechatServerToken", "")
	ttl, err := s.service.TtlWechatToken(ctx)
	if err != nil {
		l.Error("service.TtlWechatToken error", nil, err)
		return
	}
	if ttl > 10*time.Minute {
		return
	}
	resp, err := s.wechat.GetAccessToken(ctx)
	if err != nil {
		l.Error("wechat.AccessToken error", nil, err)
		return
	}
	if resp.Errcode == 0 && resp.AccessToken != "" {
		err = s.service.SetWechatToken(ctx, resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second)
		if err != nil {
			l.Error("service.SetWechatToken error", nil, err)
		}
	} else {
		l.Warn("wechat.AccessToken fail", nil, resp)
	}
}
