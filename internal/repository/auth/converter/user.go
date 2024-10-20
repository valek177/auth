package converter

import (
	"github.com/valek177/auth/internal/model"
	modelRepo "github.com/valek177/auth/internal/repository/auth/model"
)

// ToUserFromRepo converts user from repository model to internal model
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		UserInfo:  ToUserInfoFromRepo(&user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo converts user info from repository model to internal model
func ToUserInfoFromRepo(userInfo *modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  userInfo.Role,
	}
}
