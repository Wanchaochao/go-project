package handler

import (
	"github.com/gin-gonic/gin"
	"project/api/internal/proto"
	"project/model"
	"project/pkg/logger"
)

func (h *Handler) WechatLogin(c *gin.Context) {
	var r proto.LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	resp, err := h.wechat.JsCode2Session(c, r.JsCode)
	if err != nil {
		logger.FromContext(c).Error("wechat.JsCode2Session error", r.JsCode, err)
		c.JSON(RespWithErr(err))
		return
	}
	if resp.Openid == "" {
		logger.FromContext(c).Warn("wechat.JsCode2Session fail", r.JsCode, resp)
		c.JSON(RespWithMsg(Unprocessable, "Invalid Or Expired"))
		return
	}
	c.Set("v2", resp.Openid)
	c.Set("v3", resp.Unionid)
	uid, err := h.service.SaveUser(c, &model.User{
		Openid:  resp.Openid,
		Unionid: resp.Unionid,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUser error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	token, err := h.service.SetUserToken(c, &proto.UserToken{
		ID:         uid,
		Openid:     resp.Openid,
		Unionid:    resp.Unionid,
		SessionKey: resp.SessionKey,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SetUserToken error", uid, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, &proto.LoginResp{
		Token:   token,
		Openid:  resp.Openid,
		Unionid: resp.Unionid,
	})
}

func (h *Handler) WechatPhone(c *gin.Context) {
	var r proto.WechatPhoneArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	resp, err := h.wechat.GetUserPhoneNumber(c, r.Code)
	if err != nil {
		logger.FromContext(c).Error("wechat.GetUserPhoneNumber error", r.Code, err)
		c.JSON(RespWithErr(err))
		return
	}
	if resp.PhoneInfo == nil || resp.PhoneInfo.PhoneNumber == "" {
		logger.FromContext(c).Warn("wechat.GetUserPhoneNumber fail", r.Code, resp)
		c.JSON(RespWithMsg(Unprocessable, "Invalid Or Expired"))
		return
	}
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	err = h.service.UpdateUser(c, &model.User{
		ID:          user.ID,
		PhoneNumber: resp.PhoneInfo.PhoneNumber,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUserPhone error", resp.PhoneInfo, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, &proto.WechatPhoneResp{
		PhoneNumber: resp.PhoneInfo.PhoneNumber,
	})
}

func (h *Handler) SaveUserInfo(c *gin.Context) {
	var r proto.SaveUserInfoArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	err := h.service.UpdateUser(c, &model.User{
		ID:        user.ID,
		Nickname:  r.Nickname,
		AvatarURL: r.AvatarURL,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUserInfo error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	info, err := h.service.FindUserByID(c, user.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindUserByID error", user.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, &proto.GetUserInfoResp{
		PhoneNumber: info.PhoneNumber,
		Nickname:    info.Nickname,
		AvatarURL:   info.AvatarURL,
	})
}
