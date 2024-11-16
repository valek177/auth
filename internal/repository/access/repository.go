package log

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/repository/access/converter"
	repoModel "github.com/valek177/auth/internal/repository/access/model"
	"github.com/valek177/platform-common/pkg/client/db"
)

const (
	tableName = "access_list"

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
	builderSelect := sq.Select(roleColumn, endpointColumn).
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
		fmt.Println("err in scan", err)
		fmt.Println("rules", rules)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("access rules not found")
		}
		return nil, err
	}
	fmt.Println("err in scan", err)
	fmt.Println("rules", rules)

	return converter.ToEndpointAccessRuleFromRepo(endpoint, rules), nil
}
