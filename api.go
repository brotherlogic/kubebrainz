package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/kubebrainz/proto"
)

type Server struct {
}

func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "haven't got to this yet")
}

func (s *Server) GetArtist(ctx context.Context, req *pb.GetArtistRequest) (*pb.GetArtistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "haven't got to this yet")
}
