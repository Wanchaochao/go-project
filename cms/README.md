### 示例接口
- GET/ping 连通测试
- GET/modules 所有模块列表(仅非生产环境)
- GET/test/:code 模拟错误码(仅非生产环境)
- GET/captcha 获取登录验证码
- POST/user/login 登入
- DELETE/user/logout 登出
- POST/user/password 修改密码
- GET/admin/role/list 角色分页列表
- GET/admin/role/options 所有角色(仅包含id和name)
- POST/admin/role 创建角色
- PUT/admin/role 更新角色配置
- GET/admin/users/list 管理员分页列表
- POST/admin/user 创建管理员账号
- PUT/admin/user/password 重置账号密码
- PUT/admin/user/role 分配账号角色
- PUT/admin/user/status 切换账号状态
- POST/upload/file 通用文件上传(2M)
- POST/upload/image 图片上传(500KB)

### 权限管理设计
> - 以模块为单位，给每个角色指定各个模块的权限（0无权限，1只读，2操作）
> - 前端根据登录返回的账号权限模块加载菜单，根据是否有写权限显示操作按钮。
> - 后端根据账号模块权限判断接口权限，无权限的返回403。
