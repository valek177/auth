package access

import (
	"context"
	"slices"
)

// Check checks user access permissions to resource
func (s *serv) Check(ctx context.Context, accessToken string, endpoint string) (bool, error) {
	claims, err := s.tokenAccess.VerifyToken(ctx, accessToken)
	if err != nil {
		return false, err
	}

	accessRule, err := s.accessRepository.GetAccessRuleByEndpoint(ctx, endpoint)
	if err != nil {
		return false, err
	}

	if slices.Contains(accessRule.Roles, claims.Role) {
		return true, nil
	}

	return false, nil
}
