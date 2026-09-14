package service

import (
	"context"
	"errors"
	"net/http"
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
			http.StatusBadRequest,
			appErrors.ErrLoginRequestCannotBeNil,
		)
	}

	if req.Username == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			appErrors.ErrUsernameCannotBeEmpty,
		)
	}

	if req.Password == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			appErrors.ErrPasswordCannotBeEmpty,
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
				appErrors.ErrInvalidParams,
				http.StatusBadRequest,
				appErrors.ErrInvalidCredentialsMessage,
			)
		}

		// 真正的数据库错误不能伪装成 401。
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			appErrors.ErrFailedToFindUserMessage,
			err,
		)
	}

	// ==================================================
	// 3. 检查用户状态
	// ==================================================

	if user.Status != 1 {
		return nil, appErrors.New(
			appErrors.ErrUnauthorized,
			http.StatusBadRequest,
			appErrors.ErrUserDisabledMessage,
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
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			appErrors.ErrInvalidCredentialsMessage,
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
			http.StatusInternalServerError,
			appErrors.ErrFailedToGenerateAccessTokenMessage,
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
