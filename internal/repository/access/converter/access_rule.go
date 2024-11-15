package converter

import (
	"github.com/valek177/auth/internal/model"
	repoModel "github.com/valek177/auth/internal/repository/access/model"
)

func ToEndpointAccessRuleFromRepo(endpoint string, rules []*repoModel.AccessRule,
) *model.EndpointAccessRule {
	if endpoint == "" {
		return nil
	}

	if len(rules) == 0 {
		return nil
	}

	endpointRoles := make([]string, len(rules))

	for i, rule := range rules {
		endpointRoles[i] = rule.Role
	}

	resRule := &model.EndpointAccessRule{
		Endpoint: endpoint,
		Roles:    endpointRoles,
	}

	return resRule
}
