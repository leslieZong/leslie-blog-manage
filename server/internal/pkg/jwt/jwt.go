package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 表示 Leslie Blog 自己定义的 JWT Payload。
//
// JWT 可以简单理解成：
//
// Header
//
//	+
//
// Payload
//
//	+
//
// # Signature
//
// Claims 就是 Payload 中我们自己关心的数据。
type Claims struct {

	// UserID 表示当前登录用户的 ID。
	//
	// 后续访问 Admin API 时，
	// 我们可以从 JWT 中拿到这个 ID。
	UserID string `json:"userId"`

	// Username 表示当前登录用户的用户名。
	Username string `json:"username"`

	// StandardClaims 是 JWT 标准字段。
	//
	// 里面包含：
	//
	// Issuer
	// ExpiresAt
	// IssuedAt
	// NotBefore
	// 等。
	jwt.RegisteredClaims
}

// Generate 生成 JWT Token。
//
// userID：用户 ID
// username：用户名
// secret：JWT Secret
// issuer：JWT 签发者
// expireHours：Token 有效时间
func Generate(
	userID string,
	username string,
	secret string,
	issuer string,
	expireHours int,
) (string, error) {

	// 当前时间。
	now := time.Now()

	// 计算 Token 过期时间。
	expireTime := now.Add(
		time.Duration(expireHours) * time.Hour,
	)

	// 创建 Claims。
	claims := Claims{
		UserID:   userID,
		Username: username,

		RegisteredClaims: jwt.RegisteredClaims{

			// Issuer：
			// Token 是由谁签发的。
			Issuer: issuer,

			// IssuedAt：
			// Token 创建时间。
			IssuedAt: jwt.NewNumericDate(now),

			// ExpiresAt：
			// Token 过期时间。
			ExpiresAt: jwt.NewNumericDate(expireTime),

			// NotBefore：
			// Token 从什么时候开始生效。
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// 创建 JWT。
	//
	// SigningMethodHS256 表示：
	//
	// HMAC + SHA256
	//
	// 这是一个对称签名算法。
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// 使用 Secret 对 JWT 进行签名。
	signedToken, err := token.SignedString(
		[]byte(secret),
	)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Parse 解析并验证 JWT。
//
// 这个函数后面会被 JWT Middleware 使用。
func Parse(
	tokenString string,
	secret string,
) (*Claims, error) {

	// jwt.ParseWithClaims 会：
	//
	// 1. 解析 Token
	// 2. 解析 Claims
	// 3. 验证签名
	// 4. 验证过期时间等标准字段
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {

			// 防止攻击者修改 Token 的签名算法。
			//
			// 我们生成 Token 时使用 HS256，
			// 所以验证时也必须要求 HS256。
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			// 返回 Secret。
			//
			// JWT 库会使用这个 Secret
			// 验证 Token 的 Signature。
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// 类型断言。
	//
	// ParseWithClaims 返回的是 jwt.Claims 接口，
	// 我们需要把它转换成自己的 Claims。
	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// token.Valid 表示 Token 是否通过验证。
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
