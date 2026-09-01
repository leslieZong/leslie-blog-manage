package service

import (
	"errors"

	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
	"leslie-blog-server/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepo
	jwt      *utils.JWT
}

func NewAuthService(userRepo *repository.UserRepo, jwt *utils.JWT) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt}
}

// Login 校验并签发 token
func (s *AuthService) Login(req model.LoginReq) (*model.LoginResp, error) {
	u, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if !utils.CheckPassword(u.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	token, err := s.jwt.Generate(u.ID, u.Username, u.Role)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}
	return &model.LoginResp{
		Token: token,
		UserInfo: model.UserInfo{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Email:    u.Email,
			Role:     u.Role,
		},
	}, nil
}

// CurrentUser 根据 claims 获取用户信息
func (s *AuthService) CurrentUser(userID uint) (*model.UserInfo, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	info := &model.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Email:    u.Email,
		Role:     u.Role,
	}
	return info, nil
}
