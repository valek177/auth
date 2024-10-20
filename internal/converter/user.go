package converter

import (
	"database/sql"

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
		UserInfo:  ToUserInfoFromService(user.UserInfo),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserInfoFromService converts user info model to protobuf object
func ToUserInfoFromService(userInfo model.UserInfo) *user_v1.UserInfo {
	return &user_v1.UserInfo{
		Name:  wrapperspb.String(userInfo.Name),
		Email: wrapperspb.String(userInfo.Email),
		Role:  user_v1.Role(user_v1.Role_value[userInfo.Role]),
	}
}

// ToUserFromUserV1 converts user protobuf object to model
func ToUserFromUserV1(user *user_v1.User) *model.User {
	return &model.User{
		ID:        user.Id,
		UserInfo:  ToUserInfoFromUserInfoV1(user.UserInfo),
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: &sql.NullTime{Time: user.GetUpdatedAt().AsTime()},
	}
}

// ToUserInfoFromUserInfoV1 converts user info protobuf object to model
func ToUserInfoFromUserInfoV1(userInfo *user_v1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  userInfo.GetName().GetValue(),
		Email: userInfo.GetEmail().GetValue(),
		Role:  userInfo.GetRole().String(),
	}
}
