package servers

import (
	"context"
	"errors"
	pb "schwarz/api/proto"
	"schwarz/models"
	"testing"
)

// GIVEN CreatePostgres
func TestCreatePostgres(t *testing.T) {
	tcs := []struct {
		description    string
		incoming       *pb.CreatePostgresRequest
		forcedResult   string
		forcedError    error
		expectedResult string
		expectedError  error
	}{
		{
			description: "WHEN incoming data is set without error THEN current data is processed and result given",
			incoming: &pb.CreatePostgresRequest{
				DbName:     "dbName",
				UserName:   "user_name",
				UserPass:   "user_pass",
				PortNum:    10,
				Replicas:   1,
				Capacity:   "10Mi",
				AccessMode: "ReadOnlyOnce",
			},
			forcedResult:   "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
			forcedError:    nil,
			expectedError:  nil,
			expectedResult: "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
		},
		{
			description: "WHEN incoming data is set with error THEN current data is processed and error given",
			incoming: &pb.CreatePostgresRequest{
				DbName:     "dbName",
				UserName:   "user_name",
				UserPass:   "user_pass",
				PortNum:    10,
				Replicas:   1,
				Capacity:   "10Mi",
				AccessMode: "ReadOnlyOnce",
			},
			forcedResult:   "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
			forcedError:    errors.New("random"),
			expectedError:  errors.New("random"),
			expectedResult: "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			postgresService := &mockPostgresService{
				create: func(_ context.Context, request models.CreateRequest) (models.CreateResponse, error) {
					if request.Replicas != tc.incoming.Replicas {
						t.Errorf("expected Replicas = %d, received = %d", tc.incoming.Replicas, request.Replicas)
					}
					if request.PortNum != tc.incoming.PortNum {
						t.Errorf("expected PortNum = %d, received = %d", tc.incoming.PortNum, request.PortNum)
					}
					if request.AccessMode != tc.incoming.AccessMode {
						t.Errorf("expected AccessMode = %s, received = %s", tc.incoming.AccessMode, request.AccessMode)
					}
					if request.Capacity != tc.incoming.Capacity {
						t.Errorf("expected Capacity = %s, received = %s", tc.incoming.Capacity, request.Capacity)
					}
					if request.DBName != tc.incoming.DbName {
						t.Errorf("expected DbName = %s, received = %s", tc.incoming.DbName, request.DBName)
					}
					if request.UserName != tc.incoming.UserName {
						t.Errorf("expected UserName = %s, received = %s", tc.incoming.UserName, request.UserName)
					}
					if request.UserPass != tc.incoming.UserPass {
						t.Errorf("expected UserPass = %s, received = %s", tc.incoming.UserPass, request.UserPass)
					}
					return models.CreateResponse{
						ID: tc.forcedResult,
					}, tc.forcedError
				},
			}
			postgresServer := NewPostgres(postgresService)
			result, err := postgresServer.CreatePostgres(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedError != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedError == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedError, err)
			} else if result.Id != tc.expectedResult {
				t.Errorf("expected result = %s, got %s", tc.expectedResult, result)
			}
		})
	}
}

// GIVEN DeletePostgres
func TestDeletePostgres(t *testing.T) {
	tcs := []struct {
		description   string
		incoming      *pb.DeletePostgresRequest
		forcedError   error
		expectedError error
	}{
		{
			description: "WHEN incoming data is set without error THEN current data is processed and no error given",
			incoming: &pb.DeletePostgresRequest{
				Id: "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
			},
			forcedError:   nil,
			expectedError: nil,
		},
		{
			description: "WHEN incoming data is set with error THEN current data is processed and error given",
			incoming: &pb.DeletePostgresRequest{
				Id: "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
			},
			forcedError:   errors.New("random"),
			expectedError: errors.New("random"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			postgresService := &mockPostgresService{
				delete: func(_ context.Context, request models.DeleteRequest) error {
					if request.ID != tc.incoming.Id {
						t.Errorf("expected ID = %s, received = %s", tc.incoming.Id, request.ID)
					}
					return tc.forcedError
				},
			}
			postgresServer := NewPostgres(postgresService)
			_, err := postgresServer.DeletePostgres(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedError != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedError == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedError, err)
			}
		})
	}
}

// GIVEN UpdatePostgres
func TestUpdatePostgres(t *testing.T) {
	tcs := []struct {
		description   string
		incoming      *pb.UpdatePostgresRequest
		forcedError   error
		expectedError error
	}{
		{
			description: "WHEN incoming data is set without error THEN current data is processed and no error given",
			incoming: &pb.UpdatePostgresRequest{
				Id:       "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
				Replicas: 1,
			},
			forcedError:   nil,
			expectedError: nil,
		},
		{
			description: "WHEN incoming data is set with error THEN current data is processed and error given",
			incoming: &pb.UpdatePostgresRequest{
				Id:       "ac5eaeee-0c5a-4478-9af0-e7e1d94d9e2b",
				Replicas: 1,
			},
			forcedError:   errors.New("random"),
			expectedError: errors.New("random"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			postgresService := &mockPostgresService{
				update: func(_ context.Context, request models.UpdateRequest) error {
					if request.ID != tc.incoming.Id {
						t.Errorf("expected ID = %s, received = %s", tc.incoming.Id, request.ID)
					}
					if request.Replicas != tc.incoming.Replicas {
						t.Errorf("expected Replicas = %d, received = %d", tc.incoming.Replicas, request.Replicas)
					}
					return tc.forcedError
				},
			}
			postgresServer := NewPostgres(postgresService)
			_, err := postgresServer.UpdatePostgres(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedError != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedError == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedError, err)
			}
		})
	}
}

// Mocked Postgres Service
type mockPostgresService struct {
	create func(context.Context, models.CreateRequest) (models.CreateResponse, error)
	delete func(context.Context, models.DeleteRequest) error
	update func(context.Context, models.UpdateRequest) error
}

func (m *mockPostgresService) Create(ctx context.Context, request models.CreateRequest) (models.CreateResponse, error) {
	return m.create(ctx, request)
}

func (m *mockPostgresService) Delete(ctx context.Context, request models.DeleteRequest) error {
	return m.delete(ctx, request)
}

func (m *mockPostgresService) Update(ctx context.Context, request models.UpdateRequest) error {
	return m.update(ctx, request)
}
