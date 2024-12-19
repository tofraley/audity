package auditerserver

import (
	"context"

	// Assume this is where we'll put our database service

	auditerdb "github.com/tofraley/audity/internal/auditordb"
	pb "github.com/tofraley/audity/rpc/auditer"
)

type Server struct {
	AuditerService *auditerdb.AuditerService
}

func NewServer(dbPath string) (*Server, error) {
	auditerService, err := auditerdb.NewAuditerService(dbPath)
	if err != nil {
		return nil, err
	}
	return &Server{AuditerService: auditerService}, nil
}

func (s *Server) RecordNpmAudit(ctx context.Context, req *pb.NpmAuditRequest) (*pb.NpmAuditResponse, error) {
	err := s.AuditerService.RecordNpmAudit(req)
	if err != nil {
		return &pb.NpmAuditResponse{Success: false}, err
	}
	return &pb.NpmAuditResponse{Success: true}, nil
}
