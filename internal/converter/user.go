package converter

import (
	"log"

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
func ToNewUserFromNewUserV1(req *user_v1.CreateUserRequest) *model.NewUser {
	return &model.NewUser{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		Role:            req.Role.String(),
	}
}

func ToUpdateUserInfoFromV1(req *user_v1.UpdateUserRequest) *model.UpdateUserInfo {
	var ptrName, ptrRole *string

	if req.GetName() != nil {
		str := req.GetName().GetValue()
		ptrName = &str
	}
	log.Println("role is ", req.GetRole().String(), user_v1.Role_value[req.GetRole().String()])
	strrr := req.GetRole().String()
	log.Println("strrr ", strrr)

	if user_v1.Role_value[req.GetRole().String()] != 0 {
		str := req.GetRole().String()
		ptrRole = &str
	}

	return &model.UpdateUserInfo{
		ID:   req.Id,
		Name: ptrName,
		Role: ptrRole,
	}
}
