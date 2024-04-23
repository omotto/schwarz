package services

import (
	"context"
	"fmt"
	"schwarz/models"
	"testing"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/rand"
)

// GIVEN CreateValidator
func TestCreateValidator(t *testing.T) {
	validator := NewValidator(NewDefault())
	tcs := []struct {
		description string
		incoming    models.CreateRequest
		expectedErr error
	}{
		{
			description: "WHEN DBName is less than minDBNameLength (4) THEN invalidDBNameLengthError",
			incoming: models.CreateRequest{
				DBName: generateString(minDBNameLength - 1),
			},
			expectedErr: fmt.Errorf(invalidDBNameLengthError, minDBNameLength-1),
		},
		{
			description: "WHEN DBName is higher than maxDBNameLength (100) THEN invalidDBNameLengthError",
			incoming: models.CreateRequest{
				DBName: generateString(maxDBNameLength + 1),
			},
			expectedErr: fmt.Errorf(invalidDBNameLengthError, maxDBNameLength+1),
		},
		{
			description: "WHEN UserName is less than minUserNameLength (2) THEN invalidUserNameLengthError",
			incoming: models.CreateRequest{
				DBName:   generateString(minDBNameLength),
				UserName: generateString(minUserNameLength - 1),
			},
			expectedErr: fmt.Errorf(invalidUserNameLengthError, minUserNameLength-1),
		},
		{
			description: "WHEN UserName is higher than maxUserNameLength (100) THEN invalidUserNameLengthError",
			incoming: models.CreateRequest{
				DBName:   generateString(maxDBNameLength),
				UserName: generateString(maxUserNameLength + 1),
			},
			expectedErr: fmt.Errorf(invalidUserNameLengthError, maxUserNameLength+1),
		},
		{
			description: "WHEN UserPass is less than minUserPassLength (8) THEN invalidUserPassLengthError",
			incoming: models.CreateRequest{
				DBName:   generateString(minDBNameLength),
				UserName: generateString(minUserNameLength),
				UserPass: generateString(minUserPassLength - 1),
			},
			expectedErr: fmt.Errorf(invalidUserPassLengthError, minUserPassLength-1),
		},
		{
			description: "WHEN UserPass is higher than maxUserPassLength (64) THEN invalidUserPassLengthError",
			incoming: models.CreateRequest{
				DBName:   generateString(maxDBNameLength),
				UserName: generateString(maxUserNameLength),
				UserPass: generateString(maxUserPassLength + 1),
			},
			expectedErr: fmt.Errorf(invalidUserPassLengthError, maxUserPassLength+1),
		},
		{
			description: "WHEN PortNum is less than minPortNum (1024) THEN invalidPortNumError",
			incoming: models.CreateRequest{
				DBName:   generateString(minDBNameLength),
				UserName: generateString(minUserNameLength),
				UserPass: generateString(minUserPassLength),
				PortNum:  minPortNum - 1,
			},
			expectedErr: fmt.Errorf(invalidPortNumError, minPortNum-1),
		},
		{
			description: "WHEN PortNum is higher than maxPortNum (65353) THEN invalidPortNumError",
			incoming: models.CreateRequest{
				DBName:   generateString(maxDBNameLength),
				UserName: generateString(maxUserNameLength),
				UserPass: generateString(maxUserPassLength),
				PortNum:  maxPortNum + 1,
			},
			expectedErr: fmt.Errorf(invalidPortNumError, maxPortNum+1),
		},
		{
			description: "WHEN Replicas is less than minReplicas (1) THEN invalidNumReplicasError",
			incoming: models.CreateRequest{
				DBName:   generateString(minDBNameLength),
				UserName: generateString(minUserNameLength),
				UserPass: generateString(minUserPassLength),
				PortNum:  minPortNum,
				Replicas: minReplicas - 1,
			},
			expectedErr: fmt.Errorf(invalidNumReplicasError, minReplicas-1),
		},
		{
			description: "WHEN Replicas is higher than maxReplicas (10) THEN invalidNumReplicasError",
			incoming: models.CreateRequest{
				DBName:   generateString(maxDBNameLength),
				UserName: generateString(maxUserNameLength),
				UserPass: generateString(maxUserPassLength),
				PortNum:  maxPortNum,
				Replicas: maxReplicas + 1,
			},
			expectedErr: fmt.Errorf(invalidNumReplicasError, maxReplicas+1),
		},
		{
			description: "WHEN Capacity has no valid Quantity format THEN invalidCapacityError",
			incoming: models.CreateRequest{
				DBName:   generateString(maxDBNameLength),
				UserName: generateString(maxUserNameLength),
				UserPass: generateString(maxUserPassLength),
				PortNum:  maxPortNum,
				Replicas: maxReplicas,
				Capacity: "M10",
			},
			expectedErr: fmt.Errorf(invalidCapacityError, "M10"),
		},
		{
			description: "WHEN AccessMode has no valid value THEN invalidAccessModeError",
			incoming: models.CreateRequest{
				DBName:     generateString(maxDBNameLength),
				UserName:   generateString(maxUserNameLength),
				UserPass:   generateString(maxUserPassLength),
				PortNum:    maxPortNum,
				Replicas:   maxReplicas,
				Capacity:   "10Mi",
				AccessMode: "random",
			},
			expectedErr: fmt.Errorf(invalidAccessModeError, "random"),
		},
		{
			description: "WHEN all values are valid THEN error is nil",
			incoming: models.CreateRequest{
				DBName:     generateString(maxDBNameLength),
				UserName:   generateString(maxUserNameLength),
				UserPass:   generateString(maxUserPassLength),
				PortNum:    maxPortNum,
				Replicas:   maxReplicas,
				Capacity:   "10Mi",
				AccessMode: "ReadOnlyMany",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			_, err := validator.Create(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedErr == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedErr, err)
			}
		})
	}
}

// GIVEN UpdateValidator
func TestUpdateValidator(t *testing.T) {
	validator := NewValidator(NewDefault())
	tcs := []struct {
		description string
		incoming    models.UpdateRequest
		expectedErr error
	}{
		{
			description: "WHEN ID has no valid UUID format THEN invalidUUIDError",
			incoming: models.UpdateRequest{
				ID: "random",
			},
			expectedErr: fmt.Errorf(invalidUUIDError, "random"),
		},
		{
			description: "WHEN Replicas is less than minReplicas (1) THEN invalidNumReplicasError",
			incoming: models.UpdateRequest{
				ID:       uuid.New().String(),
				Replicas: minReplicas - 1,
			},
			expectedErr: fmt.Errorf(invalidNumReplicasError, minReplicas-1),
		},
		{
			description: "WHEN Replicas is higher than maxReplicas (10) THEN invalidNumReplicasError",
			incoming: models.UpdateRequest{
				ID:       uuid.New().String(),
				Replicas: maxReplicas + 1,
			},
			expectedErr: fmt.Errorf(invalidNumReplicasError, maxReplicas+1),
		},
		{
			description: "WHEN all values are valid THEN error is nil",
			incoming: models.UpdateRequest{
				ID:       uuid.New().String(),
				Replicas: maxReplicas,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			err := validator.Update(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedErr == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedErr, err)
			}
		})
	}
}

// GIVEN UpdateValidator
func TestDeleteValidator(t *testing.T) {
	validator := NewValidator(NewDefault())
	tcs := []struct {
		description string
		incoming    models.DeleteRequest
		expectedErr error
	}{
		{
			description: "WHEN ID has no valid UUID format THEN invalidUUIDError",
			incoming: models.DeleteRequest{
				ID: "random",
			},
			expectedErr: fmt.Errorf(invalidUUIDError, "random"),
		},
		{
			description: "WHEN all values are valid THEN error is nil",
			incoming: models.DeleteRequest{
				ID: uuid.New().String(),
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			err := validator.Delete(context.Background(), tc.incoming)
			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedErr == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedErr, err)
			}
		})
	}
}

func generateString(size int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
