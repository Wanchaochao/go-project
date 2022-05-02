package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) register(r *gin.Engine) {
	r.GET("ping", func(c *gin.Context) {
		c.String(OK, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(NotFound)
	})
	r.Use(Recover, SetContext) // 如nginx未添加跨域头，则此处应添加Cors中间件

	api := r.Group("", AccessLog)
	{
		api.POST("wechat/login", h.WechatLogin)
		api.GET("example/banners", h.GetBanners)
		api.POST("example/message", h.PushMessage)
	}

	{
		wx := api.Group("wechat", h.AuthCheck)
		wx.POST("phone", h.WechatPhone)
		wx.PUT("userinfo", h.SaveUserInfo)
		wx.GET("userinfo", h.GetUserInfo)
	}
}
