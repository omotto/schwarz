package services

import (
	"context"
	"schwarz/models"
)

type Service interface {
	Create(ctx context.Context, request models.CreateRequest) (models.CreateResponse, error)
	Update(ctx context.Context, request models.UpdateRequest) error
	Delete(ctx context.Context, request models.DeleteRequest) error
}
