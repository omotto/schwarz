package servers

import (
	"context"
	"schwarz/models"
	"schwarz/services/kubernetes"

	pb "schwarz/api/proto"
)

type PostgresServer struct {
	postgresService kubernetes.Service
}

func NewPostgres(postgresService kubernetes.Service) pb.PostgresServiceServer {
	return &PostgresServer{postgresService}
}

func (s *PostgresServer) CreatePostgres(ctx context.Context, req *pb.CreatePostgresRequest) (*pb.CreatePostgresResponse, error) {
	resp, err := s.postgresService.Create(ctx, models.CreateRequest{
		DBName:     req.GetDbName(),
		UserName:   req.GetUserName(),
		UserPass:   req.GetUserPass(),
		PortNum:    req.GetPortNum(),
		Replicas:   req.GetReplicas(),
		Capacity:   req.GetCapacity(),
		AccessMode: req.GetAccessMode(),
	})
	return &pb.CreatePostgresResponse{
		Id: resp.ID,
	}, err
}

func (s *PostgresServer) UpdatePostgres(ctx context.Context, req *pb.UpdatePostgresRequest) (*pb.UpdatePostgresResponse, error) {
	err := s.postgresService.Update(ctx, models.UpdateRequest{
		ID:       req.GetId(),
		Replicas: req.GetReplicas(),
	})
	return &pb.UpdatePostgresResponse{}, err
}

func (s *PostgresServer) DeletePostgres(ctx context.Context, req *pb.DeletePostgresRequest) (*pb.DeletePostgresResponse, error) {
	err := s.postgresService.Delete(ctx, models.DeleteRequest{
		ID: req.GetId(),
	})
	return &pb.DeletePostgresResponse{}, err
}
