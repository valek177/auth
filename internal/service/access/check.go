package access

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

// Check checks user access permissions to resource
func (s *serv) Check(ctx context.Context, accessToken string, endpoint string) (bool, error) {
	fmt.Println("we are in check! token is ", accessToken)
	claims, err := s.tokenAccess.VerifyToken(ctx, accessToken)
	if err != nil {
		return false, err
	}
	fmt.Println("claims: ", claims)

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
	fmt.Println("access ", accessRule)

	if slices.Contains(accessRule.Roles, claims.Role) {
		return true, nil
	}

	return false, nil
}
