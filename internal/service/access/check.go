package access

import "context"

// Check checks user access permissions to resource
func (s *serv) Check(ctx context.Context, accessToken string, endpoint string) (bool, error) {
	return false, nil
}
