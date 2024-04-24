package kubernetes

import (
	"context"
	"fmt"
	"schwarz/models"
	"strings"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	invalidDBNameLengthError   = "invalid db_name length of %d chars"
	invalidUUIDError           = "invalid %s UUID format"
	invalidNumReplicasError    = "invalid number %d of replicas "
	invalidUserNameLengthError = "invalid user_name length of %d chars"
	invalidUserPassLengthError = "invalid user_pass length of %d chars"
	invalidPortNumError        = "invalid port_num %d value"
	invalidAccessModeError     = "invalid access mode %s"
	invalidCapacityError       = "invalid capacity %s format"

	minDBNameLength   = 4
	maxDBNameLength   = 100
	minUserNameLength = 2
	maxUserNameLength = 100
	minUserPassLength = 8
	maxUserPassLength = 64
	minPortNum        = 1024
	maxPortNum        = 65353
	minReplicas       = 1
	maxReplicas       = 10
)

type Validator struct {
	service Service
}

func NewValidator(service Service) Service {
	return &Validator{
		service: service,
	}
}

func (v *Validator) Create(ctx context.Context, request models.CreateRequest) (models.CreateResponse, error) {
	if len(request.DBName) < minDBNameLength || len(request.DBName) > maxDBNameLength {
		return models.CreateResponse{}, fmt.Errorf(invalidDBNameLengthError, len(request.DBName))
	}
	if len(request.UserName) < minUserNameLength || len(request.UserName) > maxUserNameLength {
		return models.CreateResponse{}, fmt.Errorf(invalidUserNameLengthError, len(request.UserName))
	}
	if len(request.UserPass) < minUserPassLength || len(request.UserPass) > maxUserPassLength {
		return models.CreateResponse{}, fmt.Errorf(invalidUserPassLengthError, len(request.UserPass))
	}
	if request.PortNum < minPortNum || request.PortNum > maxPortNum {
		return models.CreateResponse{}, fmt.Errorf(invalidPortNumError, request.PortNum)
	}
	if request.Replicas < minReplicas || request.Replicas > maxReplicas {
		return models.CreateResponse{}, fmt.Errorf(invalidNumReplicasError, request.Replicas)
	}
	if _, err := resource.ParseQuantity(request.Capacity); err != nil {
		return models.CreateResponse{}, fmt.Errorf(invalidCapacityError, request.Capacity)
	}
	if !isValidAccessMode(request.AccessMode) {
		return models.CreateResponse{}, fmt.Errorf(invalidAccessModeError, request.AccessMode)
	}
	return v.service.Create(ctx, request)
}

func (v *Validator) Delete(ctx context.Context, request models.DeleteRequest) error {
	if !isValidUUID(request.ID) {
		return fmt.Errorf(invalidUUIDError, request.ID)
	}
	return v.service.Delete(ctx, request)
}

func (v *Validator) Update(ctx context.Context, request models.UpdateRequest) error {
	if !isValidUUID(request.ID) {
		return fmt.Errorf(invalidUUIDError, request.ID)
	}
	if request.Replicas < minReplicas || request.Replicas > maxReplicas {
		return fmt.Errorf(invalidNumReplicasError, request.Replicas)
	}
	return v.service.Update(ctx, request)
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func isValidAccessMode(a string) bool {
	validAccessModes := []string{"readwriteonce", "readonlymany", "readwritemany", "readwriteoncepod"}
	for _, validAccessMode := range validAccessModes {
		if strings.EqualFold(validAccessMode, a) {
			return true
		}
	}
	return false
}
