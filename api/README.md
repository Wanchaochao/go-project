### 示例接口
- GET/ping 连通测试
- POST/wechat/login 微信登录（小程序code2session）
- POST/wechat/phone 微信获取手机号（code换手机号）
- PUT/wechat/userinfo 更新头像昵称（更新DB和删缓存）
- GET/wechat/userinfo 获取用户信息（查询DB和设缓存）
- GET/example/banners 获取轮播广告（singleflight的使用）
- POST/example/message 投递消息到NSQ

