package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/access_v1"
)

const (
	authHeaderName = "authorization"
	authPrefix     = "Bearer "
)

// Check checks user access to resource
func (i *Implementation) Check(ctx context.Context, req *access_v1.CheckRequest) (
	*emptypb.Empty, error,
) {
	err := validateCheck(req)
	if err != nil {
		return nil, err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader := md.Get(authHeaderName)
	if len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	hasAccess, err := i.accessService.Check(ctx, accessToken, req.GetEndpointAddress())
	if err != nil {
		return nil, errors.New("unable to check access")
	}
	if !hasAccess {
		return nil, errors.New("access denied")
	}

	return &emptypb.Empty{}, nil
}
