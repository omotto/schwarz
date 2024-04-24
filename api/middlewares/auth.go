package middlewares

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authError              = "non valid authorization token header provided"
	headerAuthorizationKey = "authorization"
	tokenBearerPrefix      = "Bearer "
)

// AuthInterceptor middleware for each rpc request. This function verifies the client has the correct AUTH TOKEN.
func AuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, authError)
	}
	authHeaderValue, ok := meta[headerAuthorizationKey]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, authError)
	}
	bearerToken := authHeaderValue[0]
	if !strings.HasPrefix(bearerToken, tokenBearerPrefix) {
		return nil, status.Error(codes.Unauthenticated, authError)
	}
	// TODO: Perform the token validation here.
	return handler(ctx, req)
}
