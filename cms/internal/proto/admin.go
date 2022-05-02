package proto

import "project/cms/internal/acl"

type CaptchaResp struct {
	SessionKey  string `json:"session_key"`
	Base64Image []byte `json:"base64_image"`
}

type LoginArgs struct {
	SessionKey string `json:"session_key" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token     string        `json:"token"`
	Username  string        `json:"username"`
	Authority acl.Authority `json:"authority"`
}

type UserPasswordArgs struct {
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type AdminRoleListResp struct {
	Total int64            `json:"total"`
	List  []*AdminRoleItem `json:"list"`
}

type AdminRoleItem struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Authority  map[string]int8 `json:"authority"`
	CreateTime string          `json:"create_time"`
}

type AdminRoleCreateArgs struct {
	Name      string           `json:"name" binding:"required"`
	Authority []*AuthorityItem `json:"authority" binding:"required,dive"`
}

type AuthorityItem struct {
	Key  string `json:"key" binding:"required"`
	Code int8   `json:"code" binding:"oneof=0 1 2"`
}

type AdminRoleUpdateArgs struct {
	ID        int              `json:"id" binding:"min=1"`
	Name      string           `json:"name" binding:"required"`
	Authority []*AuthorityItem `json:"authority" binding:"required,dive"`
}

type AdminUserListArgs struct {
	Page     int    `form:"page" binding:"min=1"`
	Size     int    `form:"size" binding:"min=10,max=100"`
	RoleID   int    `form:"role_id"`
	Username string `form:"username" binding:"max=32"`
	Status   int8   `form:"status" binding:"min=-1,max=1"`
}

type AdminUserListResp struct {
	Total int64            `json:"total"`
	List  []*AdminUserItem `json:"list"`
}

type AdminUserItem struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	RoleID     int    `json:"role_id"`
	Status     int8   `json:"status"`
	CreateTime string `json:"create_time"`
}

type AdminUserCreateArgs struct {
	Username string `json:"username" binding:"min=2,max=32"`
	Password string `json:"password" binding:"min=6,max=32"`
	RoleID   int    `json:"role_id" binding:"min=1"`
}

type AdminUserPasswordArgs struct {
	ID       int    `json:"id" binding:"min=1"`
	Password string `json:"password" binding:"min=6,max=32"`
}

type AdminUserRoleArgs struct {
	ID     int `json:"id" binding:"min=1"`
	RoleID int `json:"role_id" binding:"min=1"`
}
