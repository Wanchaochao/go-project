package proto

type UserToken struct {
	ID         int    `json:"i"`
	Openid     string `json:"o"`
	Unionid    string `json:"u"`
	SessionKey string `json:"s"`
}

type LoginArgs struct {
	JsCode string `json:"js_code" binding:"required"`
}

type LoginResp struct {
	Token   string `json:"token"`
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}

type WechatPhoneArgs struct {
	Code string `json:"code" binding:"required"`
}

type WechatPhoneResp struct {
	PhoneNumber string `json:"phone_number"`
}

type SaveUserInfoArgs struct {
	Nickname  string `json:"nickname" binding:"required"`
	AvatarURL string `json:"avatar_url" binding:"required"`
}

type GetUserInfoResp struct {
	PhoneNumber string `json:"phone_number"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
}
