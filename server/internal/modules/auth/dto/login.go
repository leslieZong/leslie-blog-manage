package dto

// LoginRequest 是登录接口的请求参数。
//
// 对应 HTTP：
//
// POST /api/admin/v1/auth/login
//
// JSON：
//
//	{
//	    "username": "admin",
//	    "password": "123456"
//	}
type LoginRequest struct {

	// Username 登录用户名。
	Username string `json:"username"`

	// Password 登录密码。
	//
	// 注意：
	//
	// 这个字段只存在于 HTTP 请求中。
	// 我们绝对不会把它保存到数据库。
	Password string `json:"password"`
}

// LoginResponse 是登录成功后的响应。
type LoginResponse struct {

	// AccessToken 是登录成功后返回给前端的 JWT。
	AccessToken string `json:"accessToken"`

	// ExpiresIn 表示 Token 有效时间，单位秒。
	//
	// 前端可以根据这个值判断 Token 大概什么时候过期。
	ExpiresIn int64 `json:"expiresIn"`
}
