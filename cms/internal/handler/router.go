package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/acl"
	"strconv"
)

func (h *Handler) register(r *gin.Engine) {
	r.GET("ping", func(c *gin.Context) {
		c.String(OK, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(NotFound)
	})
	r.Use(Recover, SetContext) // 如nginx未添加跨域头，则此处应添加Cors中间件

	if gin.Mode() != gin.ReleaseMode {
		r.GET("modules", func(c *gin.Context) {
			c.JSON(OK, gin.H{"list": acl.Modules})
		})
		r.GET("test/:code", func(c *gin.Context) {
			code, _ := strconv.Atoi(c.Param("code"))
			if code < 200 && code > 599 {
				code = 400
			}
			c.Status(code)
		})
	}
	r.GET("captcha", h.Captcha)

	{
		user := r.Group("user")
		user.POST("login", h.UserLogin)
		user.DELETE("logout", h.AuthCheck(""), h.UserLogout)
		user.PUT("password", h.AuthCheck(""), h.UserPassword)
	}

	{
		admin := r.Group("admin", h.AuthCheck(acl.ModuleAdmin), AccessLog)
		admin.GET("role/list", h.AdminRoleList)
		admin.GET("role/option", h.AdminRoleOption)
		admin.POST("role", h.AdminRoleCreate)
		admin.PUT("role", h.AdminRoleUpdate)
		admin.GET("user/list", h.AdminUserList)
		admin.POST("user", h.AdminUserCreate)
		admin.PUT("user/password", h.AdminUserPassword)
		admin.PUT("user/role", h.AdminUserRole)
		admin.PUT("user/status", h.AdminUserStatus)
	}

	{
		upload := r.Group("upload", h.AuthCheck(""))
		upload.POST("file", h.UploadFile)
		upload.POST("image", h.UploadImage)
	}

}
