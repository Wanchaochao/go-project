package handler

import (
	"github.com/gin-gonic/gin"
	"project/api/internal/proto"
	"project/pkg/logger"
	"strings"
	"time"
)

func (h *Handler) GetBanners(c *gin.Context) {
	city := c.Query("city")
	data, err := h.service.GetBannersByCity(c, city)
	if err != nil {
		logger.FromContext(c).Error("service.GetBannersByCity error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	list := make([]*proto.BannerItem, 0, len(data))
	now := time.Now().Unix()
	for _, d := range data {
		if d.BeginTime <= now && now < d.EndTime {
			if !strings.HasPrefix(d.Img, "http") { //相对路径拼上cdn域名
				d.Img = h.cdn + d.Img
			}
			list = append(list, &proto.BannerItem{
				Title: d.Title,
				Img:   d.Img,
				Type:  d.Type,
				Link:  d.Link,
			})
		}
	}
	c.JSON(OK, &proto.BannersResp{
		List: list,
	})
}

func (h *Handler) PushMessage(c *gin.Context) {
	//err := h.service.PushMessage(c, &model.MsgExample{
	//	UUID:   random.UUID(),
	//	Number: time.Now().UnixMicro(),
	//})
	//if err != nil {
	//	logger.FromContext(c).Error("service.PushMessage error", nil, err)
	//	c.JSON(RespWithErr(err))
	//	return
	//}
	c.JSON(OK, Empty)
}
