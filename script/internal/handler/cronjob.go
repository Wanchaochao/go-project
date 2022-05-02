package handler

import (
	"project/model"
	"project/pkg/dingtalk"
	"project/pkg/logger"
	"project/pkg/util/random"
	"project/pkg/wechat"
	"project/pkg/wechatwork"
	"project/script/internal/service"
	"strconv"
	"strings"
	"time"
)

type Cronjob struct {
	service     *service.Service
	wechat      wechat.ServerAPI
	robotDing   string
	robotWechat string
}

func NewCronjob(srv *service.Service, api wechat.ServerAPI, robotDing, robotWechat string) *Cronjob {
	return &Cronjob{
		service:     srv,
		wechat:      api,
		robotDing:   robotDing,
		robotWechat: robotWechat,
	}
}

func (h *Cronjob) LoadWechatAnalysis() {
	ctx, l := logger.NewCtxLog(random.UUID(), "Cronjob", "LoadWechatAnalysis", "")
	yesterday := time.Now().AddDate(0, 0, -1).Format(wechat.DateFormat)
	args := &wechat.DatacubeArgs{
		BeginDate: yesterday,
		EndDate:   yesterday,
	}

	trend, err := h.wechat.GetDailyVisitTrend(ctx, args)
	if err != nil {
		trend, err = h.wechat.GetDailyVisitTrend(ctx, args)
	}
	if err != nil {
		l.Error("wechat.GetDailyVisitTrend error", args, err)
		return
	}
	if trend.Errcode != 0 || len(trend.List) == 0 {
		l.Warn("wechat.GetDailyVisitTrend fail", args, trend)
		return
	}

	summary, err := h.wechat.GetDailySummary(ctx, args)
	if err != nil {
		summary, err = h.wechat.GetDailySummary(ctx, args)
	}
	if err != nil {
		l.Error("wechat.GetDailySummary error", args, err)
		return
	}
	if summary.Errcode != 0 || len(summary.List) == 0 {
		l.Warn("wechat.GetDailySummary fail", args, summary)
		return
	}

	data := &model.WechatAnalysis{
		RefDate:    yesterday,
		SessionCnt: trend.List[0].SessionCnt,
		VisitPv:    trend.List[0].VisitPv,
		VisitUv:    trend.List[0].VisitUv,
		VisitUvNew: trend.List[0].VisitUvNew,
		SharePv:    summary.List[0].SharePv,
		ShareUv:    summary.List[0].ShareUv,
	}
	err = h.service.SaveWechatAnalysis(ctx, data)
	if err != nil {
		l.Error("service.SaveAnalysisDailyTrend error", data, err)
	}
	h.notifyWechatAnalysis(summary.List[0].VisitTotal, data)
}

func (h *Cronjob) notifyWechatAnalysis(total int, data *model.WechatAnalysis) {
	text := &strings.Builder{}
	text.WriteString(data.RefDate)
	text.WriteString("\n打开次数:")
	text.WriteString(strconv.Itoa(data.SessionCnt))
	text.WriteString("\n访问PV:")
	text.WriteString(strconv.Itoa(data.VisitPv))
	text.WriteString("\n访问UV:")
	text.WriteString(strconv.Itoa(data.VisitUv))
	text.WriteString("\n新增用户:")
	text.WriteString(strconv.Itoa(data.VisitUvNew))
	text.WriteString("\n转发次数:")
	text.WriteString(strconv.Itoa(data.SharePv))
	text.WriteString("\n转发人数:")
	text.WriteString(strconv.Itoa(data.ShareUv))
	text.WriteString("\n=====================\n累计访问人数:")
	text.WriteString(strconv.Itoa(total))

	content := text.String()
	_, _ = wechatwork.SendText(h.robotWechat, &wechatwork.Text{Content: content})
	_, _ = dingtalk.SendText(h.robotDing, &dingtalk.Text{Content: content}, nil)
}
