package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/kubebrainz/proto"
	"log"
)

func (s *Server) reportFailure(ctx context.Context, artist string) {
	s.mu.Lock()
	if s.activeIssues[artist] {
		s.mu.Unlock()
		return
	}
	s.activeIssues[artist] = true
	s.mu.Unlock()

	log.Printf("Reporting failure for artist: %v", artist)
	err := s.githubridge.PostIssue(artist, "Resolution Failure", "Unable to resolve artist "+artist)
	if err != nil {
		log.Printf("Error reporting failure: %v", err)
	}
}

func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	// In this version we just serve
	return &pb.GetStatusResponse{
		Version: s.version,
	}, nil
}

func (s *Server) GetArtist(ctx context.Context, req *pb.GetArtistRequest) (*pb.GetArtistResponse, error) {
	rows, err := s.db.Query("SELECT sort_name FROM artist WHERE name = $1",
		req.GetArtist())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sortName string
	for rows.Next() {
		if err := rows.Scan(&sortName); err == nil {
			return &pb.GetArtistResponse{ArtistSort: sortName}, nil
		}
	}

	go s.reportFailure(ctx, req.GetArtist())
	return nil, status.Errorf(codes.NotFound, "Could not locate %v in db", req.GetArtist())
}
