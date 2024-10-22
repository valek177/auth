package converter

import (
	"github.com/valek177/auth/internal/model"
	modelRepo "github.com/valek177/auth/internal/repository/auth/model"
)

// ToUserFromRepo converts user from repository model to service model
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
