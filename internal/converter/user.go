package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/model"
)

// ToUserV1FromService converts user model to protobuf object
func ToUserV1FromService(user *model.User) *user_v1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &user_v1.User{
		Id:        user.ID,
		UserInfo:  ToUserInfoFromService(user),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserInfoFromService converts user info model to protobuf object
func ToUserInfoFromService(user *model.User) *user_v1.UserInfo {
	return &user_v1.UserInfo{
		Name:  wrapperspb.String(user.Name),
		Email: wrapperspb.String(user.Email),
		Role:  user_v1.Role(user_v1.Role_value[user.Role]),
	}
}

// ToNewUserFromUserV1 converts user protobuf object to model
func ToNewUserFromNewUserV1(newUser *user_v1.NewUser) *model.NewUser {
	return &model.NewUser{
		Name:            newUser.Name,
		Email:           newUser.Email,
		Password:        newUser.Password,
		PasswordConfirm: newUser.PasswordConfirm,
		Role:            newUser.Role.String(),
	}
}

func ToUpdateUserInfoFromV1(updateUserInfo *user_v1.UpdateUserInfo) *model.UpdateUserInfo {
	return &model.UpdateUserInfo{
		ID:   updateUserInfo.Id,
		Name: updateUserInfo.GetName().GetValue(),
		Role: updateUserInfo.GetRole().String(),
	}
}
