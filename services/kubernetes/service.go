package kubernetes

import (
	"context"
	"schwarz/models"
)

type Service interface {
	Create(ctx context.Context, request models.CreateRequest) (models.CreateResponse, error)
	Update(ctx context.Context, request models.UpdateRequest) error
	Delete(ctx context.Context, request models.DeleteRequest) error
}

type DefaultService struct{}

func NewDefault() Service {
	return &DefaultService{}
}

func (d *DefaultService) Create(context.Context, models.CreateRequest) (models.CreateResponse, error) {
	return models.CreateResponse{}, nil
}

func (d *DefaultService) Delete(context.Context, models.DeleteRequest) error {
	return nil
}

func (d *DefaultService) Update(context.Context, models.UpdateRequest) error {
	return nil
}
