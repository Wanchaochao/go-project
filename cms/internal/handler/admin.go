package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"project/cms/internal/acl"
	"project/cms/internal/proto"
	"project/model"
	"project/pkg/logger"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) captchaSign(exp, code string) string {
	mac := hmac.New(sha1.New, []byte(h.captcha))
	mac.Write([]byte(exp))
	mac.Write([]byte(strings.ToUpper(code)))
	return base32.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (h *Handler) Captcha(c *gin.Context) {
	code, bin := h.drawer.Generate(4)
	exp := strconv.FormatInt(time.Now().Unix()+65, 10)
	key := exp + "." + h.captchaSign(exp, code)
	c.JSON(OK, &proto.CaptchaResp{
		SessionKey:  key,
		Base64Image: bin,
	})
}

func (h *Handler) UserLogin(c *gin.Context) {
	var r proto.LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}

	sli := strings.Split(r.SessionKey, ".")
	if len(sli) != 2 {
		c.JSON(RespWithMsg(InvalidParam, "Invalid SessionKey"))
		return
	}
	exp, _ := strconv.ParseInt(sli[0], 10, 64)
	if time.Now().Unix() > exp {
		c.JSON(RespWithMsg(Unprocessable, "验证码过期"))
		return
	}
	sign := h.captchaSign(sli[0], r.Captcha)
	if sli[1] != sign {
		c.JSON(RespWithMsg(Unauthorized, "验证码错误"))
		return
	}
	user, err := h.service.FindAdminUserByName(c, r.Username)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminByUsername error", r.Username, err)
		c.JSON(RespWithErr(err))
		return
	}
	if user.ID == 0 || !acl.CheckPassword(r.Password, user.Password) {
		c.JSON(RespWithMsg(Unauthorized, "用户名或密码错误"))
		return
	}
	if user.Status != model.StatusOn {
		c.JSON(RespWithMsg(Unauthorized, "账号已禁用，请联系管理员"))
		return
	}
	pt := &acl.AdminToken{
		ID:       user.ID,
		Username: user.Username,
	}
	if user.AdminRole != nil {
		pt.Authority = user.AdminRole.Authority
	} else {
		pt.Authority = make(acl.Authority)
	}
	token, err := h.service.SetAdminToken(c, pt)
	if err != nil {
		logger.FromContext(c).Error("service.SetAdminToken error", user, err)
		c.JSON(RespWithErr(err))
		return
	}
	if user.Username == acl.Super {
		pt.Authority = acl.AllAuthority
	}
	c.JSON(OK, &proto.LoginResp{
		Token:     token,
		Username:  pt.Username,
		Authority: pt.Authority,
	})
}

func (h *Handler) UserLogout(c *gin.Context) {
	err := h.service.DelAdminToken(c, c.GetHeader("Authorization"))
	if err != nil {
		logger.FromContext(c).Error("service.DelAdminToken error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) UserPassword(c *gin.Context) {
	var r proto.UserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	v, _ := c.Get("user")
	user := v.(*acl.AdminToken)
	err := h.service.UpdateAdminUser(c, &acl.AdminUser{
		ID:       user.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdminUser error", user.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminRoleList(c *gin.Context) {
	var r proto.ListArgs
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	total, list, err := h.service.PaginateAdminRole(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.PaginateAdminRole error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	items := make([]*proto.AdminRoleItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminRoleItem{
			ID:         v.ID,
			Name:       v.Name,
			Authority:  v.Authority,
			CreateTime: v.CreateTime.Format(TimeFormat),
		})
	}
	c.JSON(OK, &proto.AdminRoleListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminRoleOption(c *gin.Context) {
	data, err := h.service.AllAdminRole(c)
	if err != nil {
		logger.FromContext(c).Error("service.AllAdminRole error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	opt := make([]*proto.OptionItem, 0, len(data))
	for _, d := range data {
		opt = append(opt, &proto.OptionItem{
			ID:   d.ID,
			Name: d.Name,
		})
	}
	c.JSON(OK, &proto.OptionResp{List: opt})
}

func (h *Handler) AdminRoleCreate(c *gin.Context) {
	var r proto.AdminRoleCreateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	auth := make(acl.Authority)
	for _, v := range r.Authority {
		if _, ok := acl.AllAuthority[v.Key]; ok && v.Code > 0 {
			auth[v.Key] = v.Code
		}
	}
	err := h.service.CreateAdminRole(c, &acl.AdminRole{
		Name:      r.Name,
		Authority: auth,
	})
	if err != nil {
		logger.FromContext(c).Error("service.CreateAdminRole error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminRoleUpdate(c *gin.Context) {
	var r proto.AdminRoleUpdateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	auth := make(acl.Authority)
	for _, v := range r.Authority {
		if _, ok := acl.AllAuthority[v.Key]; ok && v.Code > 0 {
			auth[v.Key] = v.Code
		}
	}
	role, err := h.service.FindAdminRoleByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminRoleByID error", r.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户角色"))
		return
	}
	if r.Name == role.Name && reflect.DeepEqual(auth, role.Authority) {
		c.JSON(OK, Empty)
		return
	}
	err = h.service.UpdateAdminRole(c, &acl.AdminRole{
		ID:        r.ID,
		Name:      r.Name,
		Authority: auth,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdminRole error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserList(c *gin.Context) {
	var r proto.AdminUserListArgs
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	total, list, err := h.service.PaginateAdminUser(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.PaginateAdminUser error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	items := make([]*proto.AdminUserItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminUserItem{
			ID:         v.ID,
			Username:   v.Username,
			RoleID:     v.RoleID,
			Status:     v.Status,
			CreateTime: v.CreateTime.Format(TimeFormat),
		})
	}
	c.JSON(OK, &proto.AdminUserListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminUserCreate(c *gin.Context) {
	var r proto.AdminUserCreateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	if r.Username == acl.Super {
		c.JSON(RespWithMsg(InvalidParam, "该用户名不能使用"))
		return
	}
	role, err := h.service.FindAdminRoleByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminRoleByID error", r.RoleID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户角色"))
		return
	}
	ok, err := h.service.CreateAdminUser(c, &acl.AdminUser{
		Username: r.Username,
		Password: r.Password,
		RoleID:   r.RoleID,
		Status:   model.StatusOn,
	})
	if err != nil {
		logger.FromContext(c).Error("service.CreateAdminUser error", r.RoleID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if !ok {
		c.JSON(RespWithMsg(Conflict, "用户名已存在"))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserPassword(c *gin.Context) {
	var r proto.AdminUserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	user, err := h.service.FindAdminUserByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminUserByID error", r.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户ID"))
		return
	}
	err = h.service.UpdateAdminUser(c, &acl.AdminUser{
		ID:       r.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdminUser error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserRole(c *gin.Context) {
	var r proto.AdminUserRoleArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	user, err := h.service.FindAdminUserByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminUserByID error", r.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户ID"))
		return
	}
	if user.RoleID == r.RoleID {
		c.JSON(OK, Empty)
		return
	}
	role, err := h.service.FindAdminRoleByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminRoleByID error", r.RoleID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户角色"))
		return
	}
	err = h.service.UpdateAdminUser(c, &acl.AdminUser{
		ID:     r.ID,
		RoleID: r.RoleID,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdminUser error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserStatus(c *gin.Context) {
	var r proto.SwitchStatusArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithErr(err))
		return
	}
	user, err := h.service.FindAdminUserByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminUserByID error", r.ID, err)
		c.JSON(RespWithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(RespWithMsg(InvalidParam, "无效的用户ID"))
		return
	}
	if user.Status == r.Status {
		c.JSON(OK, Empty)
		return
	}
	err = h.service.UpdateAdminUser(c, &acl.AdminUser{
		ID:     r.ID,
		Status: r.Status,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdminUser error", &r, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, Empty)
}
