package model

import "strconv"

// 定义缓存使用的key，同一个redis集群的key收敛到同一文件

const (
	KeyWechatToken = "wx:tk" // 微信access_token

	keyBanners   = "banners:" // +city
	keyUserToken = "utk:"     // +token
	keyUserInfo  = "user:"    // +uid

	keyAdminSSO   = "asso:" // +id
	keyAdminToken = "atk:"  // +token
)

func BannersKey(city string) string {
	return keyBanners + city
}

func UserTokenKey(token string) string {
	return keyUserToken + token
}

func UserInfoKey(id int) string {
	return keyUserInfo + strconv.Itoa(id)
}

func AdminSSOKey(id int) string {
	return keyAdminSSO + strconv.Itoa(id)
}

func AdminTokenKey(token string) string {
	return keyAdminToken + token
}
