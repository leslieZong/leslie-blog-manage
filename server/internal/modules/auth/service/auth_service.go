package service

import (
	"context"
	"errors"
	"time"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/auth/dto"
	"leslie-blog-server/internal/modules/user/repository"
	"leslie-blog-server/internal/pkg/jwt"
	"leslie-blog-server/internal/pkg/password"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(
		ctx context.Context,
		req *dto.LoginRequest,
	) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtIssuer      string
	jwtExpireHours int
}

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

// Login 用户登录。
func (s *authService) Login(
	ctx context.Context,
	req *dto.LoginRequest,
) (*dto.LoginResponse, error) {

	// ==================================================
	// 1. 参数校验
	// ==================================================

	if req == nil {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"login request cannot be nil",
		)
	}

	if req.Username == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"username cannot be empty",
		)
	}

	if req.Password == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"password cannot be empty",
		)
	}

	// ==================================================
	// 2. 查询用户
	// ==================================================

	user, err := s.userRepo.FindByUsername(
		ctx,
		req.Username,
	)

	if err != nil {

		// 为了避免用户名枚举，
		// “用户不存在”和“密码错误”
		// 对外统一返回认证失败。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(
				appErrors.ErrInvalidCredentials,
				401,
				"invalid username or password",
			)
		}

		// 真正的数据库错误不能伪装成 401。
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find user",
			err,
		)
	}

	// ==================================================
	// 3. 检查用户状态
	// ==================================================

	if user.Status != 1 {
		return nil, appErrors.New(
			appErrors.ErrUserDisabled,
			403,
			"user is disabled",
		)
	}

	// ==================================================
	// 4. 验证密码
	// ==================================================

	if !password.Compare(
		user.PasswordHash,
		req.Password,
	) {
		return nil, appErrors.New(
			appErrors.ErrInvalidCredentials,
			401,
			"invalid username or password",
		)
	}

	// ==================================================
	// 5. 生成 JWT
	// ==================================================

	accessToken, err := jwt.Generate(
		user.ID,
		user.Username,
		s.jwtSecret,
		s.jwtIssuer,
		s.jwtExpireHours,
	)

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to generate access token",
			err,
		)
	}

	// ==================================================
	// 6. 计算 Token 有效秒数
	// ==================================================

	expiresIn := int64(
		time.Duration(s.jwtExpireHours) *
			time.Hour /
			time.Second,
	)

	// ==================================================
	// 7. 返回登录结果
	// ==================================================

	return &dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}, nil
}
