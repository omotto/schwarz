package middlewares

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GIVEN AuthInterceptor
func TestAuthInterceptor(t *testing.T) {
	tcs := []struct {
		description string
		incoming    context.Context
		expectedErr error
	}{
		{
			description: "WHEN there are no fields THEN Unauthenticated error",
			incoming:    context.Background(),
			expectedErr: status.Error(codes.Unauthenticated, authError),
		},
		{
			description: "WHEN there is no authorization field THEN Unauthenticated error",
			incoming:    metadata.NewIncomingContext(context.Background(), metadata.Pairs("auth", "random")),
			expectedErr: status.Error(codes.Unauthenticated, authError),
		},
		{
			description: "WHEN there is authorization field but not Bearer prefix THEN Unauthenticated error",
			incoming:    metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "random")),
			expectedErr: status.Error(codes.Unauthenticated, authError),
		},
		{
			description: "WHEN there is authorization field with Bearer prefix THEN no error",
			incoming:    metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer ***********")),
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			_, err := AuthInterceptor(tc.incoming, struct{}{}, &grpc.UnaryServerInfo{}, func(context.Context, any) (any, error) { return nil, nil })
			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedErr == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedErr, err)
			}
		})
	}
}
