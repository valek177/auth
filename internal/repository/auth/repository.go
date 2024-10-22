package auth

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/valek177/auth/internal/client/db"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/repository/auth/converter"
	modelRepo "github.com/valek177/auth/internal/repository/auth/model"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates new repository object
func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

// CreateUser creates new user with specified parameters
func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn).
		Values(user.UserInfo.Name, user.UserInfo.Email, user.UserInfo.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "auth_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUser returns info about user
func (r *repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn,
		createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.GetUser",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser updates user info by id
func (r *repo) UpdateUser(ctx context.Context, user *model.User) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, user.UserInfo.Name).
		Set(roleColumn, user.UserInfo.Role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.UpdateUser",
		QueryRaw: query,
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

// DeleteUser removes user
func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.DeleteUser",
		QueryRaw: query,
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}
