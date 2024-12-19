package auditerserver

import (
	"context"

	pb "github.com/tofraley/audity/rpc/auditer"
)

type Server struct{}

func (s *Server) RecordNpmAudit(ctx context.Context, req *pb.NpmAuditRequest) (res *pb.NpmAuditResponse, err error) {
	return &pb.NpmAuditResponse{Success: true}, nil
}
