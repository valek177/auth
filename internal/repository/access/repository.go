package log

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/repository/access/converter"
	repoModel "github.com/valek177/auth/internal/repository/access/model"
	"github.com/valek177/platform-common/pkg/client/db"
)

const (
	tableName = "roles_users_access"

	idColumn       = "id"
	roleColumn     = "role"
	endpointColumn = "endpoint"
)

type repo struct {
	db db.Client
}

// NewRepository creates new repository
func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r *repo) GetAccessRuleByEndpoint(ctx context.Context, endpoint string) (
	*model.EndpointAccessRule, error,
) {
	builderSelect := sq.Select(idColumn, roleColumn, endpointColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{endpointColumn: endpoint})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "access_repository.GetAccessRuleByEndpoint",
		QueryRaw: query,
	}

	var rules []*repoModel.AccessRule
	err = r.db.DB().ScanAllContext(ctx, &rules, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("access rules not found")
		}
		return nil, err
	}

	return converter.ToEndpointAccessRuleFromRepo(endpoint, rules), nil
}
