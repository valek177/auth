package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/config/env"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser creates new user with specified parameters
func (s *server) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (
	*user_v1.CreateUserResponse, error,
) {
	log.Printf("Create new user with name %s and email %s", req.GetName(), req.GetEmail())

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(req.GetName(), req.GetEmail(), req.GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, errors.Errorf("failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, errors.Errorf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &user_v1.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser returns info about user
func (s *server) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (
	*user_v1.GetUserResponse, error,
) {
	id := req.GetId()
	log.Printf("Get user info by id: %d", id)

	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, errors.Errorf("failed to build query: %v", err)
	}

	var name, email, role string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt,
		&updatedAt)
	if err != nil {
		return nil, errors.Errorf("failed to select users: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n",
		id, name, email, role, createdAt, updatedAt)

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &user_v1.GetUserResponse{
		User: &user_v1.User{
			Id: id,
			UserInfo: &user_v1.UserInfo{
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  user_v1.Role(user_v1.Role_value[role]),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

// UpdateUser updates user info by id
func (s *server) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (
	*emptypb.Empty, error,
) {
	id := req.GetId()

	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.GetName().GetValue()).
		Set("role", req.GetRole()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, errors.Errorf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.Errorf("failed to update user: %v", err)
	}

	log.Printf("updated user with id: %d", id)

	return &emptypb.Empty{}, nil
}

// DeleteUser removes user
func (s *server) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (
	*emptypb.Empty, error,
) {
	id := req.GetId()

	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, errors.Errorf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.Errorf("failed to delete user: %v", err)
	}

	log.Printf("deleted user with id: %d", id)

	return &emptypb.Empty{}, nil
}
