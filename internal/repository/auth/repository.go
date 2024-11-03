package auth

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/repository/auth/converter"
	modelRepo "github.com/valek177/auth/internal/repository/auth/model"
	"github.com/valek177/platform-common/pkg/client/db"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
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
func (r *repo) CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(newUser.Name, newUser.Email, newUser.Password, newUser.Role).
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
func (r *repo) UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: updateUserInfo.ID})

	isUpdated := false

	if updateUserInfo.Name != nil {
		builderUpdate = builderUpdate.Set(nameColumn, *updateUserInfo.Name)
		isUpdated = true
	}
	if updateUserInfo.Role != nil {
		builderUpdate = builderUpdate.Set(roleColumn, *updateUserInfo.Role)
		isUpdated = true
	}

	if !isUpdated {
		return errors.New("unable to update users: empty request")
	}

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
