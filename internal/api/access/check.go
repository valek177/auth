package access

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "check access api")
	defer span.Finish()

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
		return nil, fmt.Errorf("check access error: %v", err.Error())
	}
	if !hasAccess {
		return nil, errors.New("access denied")
	}

	return &emptypb.Empty{}, nil
}
