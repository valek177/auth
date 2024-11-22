package access

import (
	"context"
	"errors"
	"slices"

	"github.com/opentracing/opentracing-go"
)

// Check checks user access permissions to resource
func (s *serv) Check(ctx context.Context, accessToken string, endpoint string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "check access (service)")
	defer span.Finish()

	claims, err := s.tokenAccess.VerifyToken(ctx, accessToken)
	if err != nil {
		return false, err
	}

	accessRule, err := s.accessRepository.GetAccessRuleByEndpoint(ctx, endpoint)
	if err != nil {
		return false, err
	}
	if accessRule == nil {
		return false, errors.New("unable to find access rule")
	}
	if len(accessRule.Roles) == 0 {
		return false, errors.New("no roles in access rule")
	}

	if slices.Contains(accessRule.Roles, claims.Role) {
		return true, nil
	}

	return false, nil
}
