package service

import (
	"context"
	"errors"
	"time"

	"leslie-blog-server/internal/modules/auth/dto"
	"leslie-blog-server/internal/modules/user/repository"
	"leslie-blog-server/internal/pkg/jwt"
	"leslie-blog-server/internal/pkg/password"

	"gorm.io/gorm"
)

// AuthService 定义认证相关业务。
type AuthService interface {

	// Login 执行用户登录。
	Login(
		ctx context.Context,
		req *dto.LoginRequest,
	) (*dto.LoginResponse, error)
}

// authService 是 AuthService 的具体实现。
type authService struct {

	// userRepo 用于查询用户。
	//
	// Auth 不需要自己实现用户查询，
	// 直接复用 UserRepository。
	userRepo repository.UserRepository

	// jwtSecret 是 JWT 签名密钥。
	jwtSecret string

	// jwtIssuer 是 JWT 签发者。
	jwtIssuer string

	// jwtExpireHours 是 Token 有效时间。
	jwtExpireHours int
}

// NewAuthService 创建 AuthService。
func NewAuthService(
	userRepo repository.UserRepository,
	jwtSecret string,
	jwtIssuer string,
	jwtExpireHours int,
) AuthService {

	return &authService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtIssuer:      jwtIssuer,
		jwtExpireHours: jwtExpireHours,
	}
}

// Login 执行登录业务。
func (s *authService) Login(
	ctx context.Context,
	req *dto.LoginRequest,
) (*dto.LoginResponse, error) {

	// ----------------------------------------
	// 1. 参数校验
	// ----------------------------------------

	if req == nil {
		return nil, errors.New("login request cannot be nil")
	}

	if req.Username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if req.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// ----------------------------------------
	// 2. 根据 username 查询用户
	// ----------------------------------------

	user, err := s.userRepo.FindByUsername(
		ctx,
		req.Username,
	)

	if err != nil {

		// 用户不存在。
		//
		// 对登录接口来说，
		// “用户名不存在”和“密码错误”
		// 最终都应该表现成登录失败。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}

		// 数据库真正发生异常。
		return nil, err
	}

	// ----------------------------------------
	// 3. 检查用户状态
	// ----------------------------------------

	// 我们之前 users 表设计了：
	//
	// status = 1 → 正常
	// status = 0 → 禁用
	//
	// 如果账号被禁用，
	// 即使密码正确也不能登录。
	if user.Status != 1 {
		return nil, errors.New("user is disabled")
	}

	// ----------------------------------------
	// 4. 验证密码
	// ----------------------------------------

	// 注意这里：
	//
	// 第一个参数：
	// 数据库里的 password_hash
	//
	// 第二个参数：
	// 用户刚刚输入的明文 password
	if !password.Compare(
		user.PasswordHash,
		req.Password,
	) {
		return nil, errors.New(
			"invalid username or password",
		)
	}

	// ----------------------------------------
	// 5. 生成 JWT
	// ----------------------------------------

	accessToken, err := jwt.Generate(
		user.ID,
		user.Username,
		s.jwtSecret,
		s.jwtIssuer,
		s.jwtExpireHours,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------
	// 6. 计算 Token 有效时间
	// ----------------------------------------

	expiresIn := int64(
		time.Duration(s.jwtExpireHours) * time.Hour /
			time.Second,
	)

	// ----------------------------------------
	// 7. 返回登录结果
	// ----------------------------------------

	return &dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}, nil
}
