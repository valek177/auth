package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/valek177/auth/grpc/pkg/user_v1"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const grpcPort = 50051

type server struct {
	user_v1.UnimplementedUserV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser creates new user with specified parameters
func (s *server) CreateUser(_ context.Context, req *user_v1.CreateUserRequest) (
	*user_v1.CreateUserResponse, error,
) {
	id := gofakeit.Int64()
	log.Printf("Create new user with name %s and email %s", req.GetName(), req.GetEmail())

	return &user_v1.CreateUserResponse{
		Id: id,
	}, nil
}

// GetUser returns info about user
func (s *server) GetUser(_ context.Context, req *user_v1.GetUserRequest) (
	*user_v1.GetUserResponse, error,
) {
	id := req.GetId()
	log.Printf("Get user info by id: %d", id)

	return &user_v1.GetUserResponse{
		User: &user_v1.User{
			Id: req.GetId(),
			UserInfo: &user_v1.UserInfo{
				Name:  wrapperspb.String(gofakeit.Name()),
				Email: wrapperspb.String(gofakeit.Email()),
				Role:  user_v1.Role_USER,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// UpdateUser updates user info by id
func (s *server) UpdateUser(_ context.Context, req *user_v1.UpdateUserRequest) (
	*emptypb.Empty, error,
) {
	id := req.GetId()
	log.Printf("Update user info by id: %d", id)

	return &emptypb.Empty{}, nil
}

// DeleteUser removes user
func (s *server) DeleteUser(_ context.Context, req *user_v1.DeleteUserRequest) (
	*emptypb.Empty, error,
) {
	id := req.GetId()
	log.Printf("Delete user by id: %d", id)

	return &emptypb.Empty{}, nil
}
