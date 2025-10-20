package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/kubebrainz/proto"
)

func InitTestServer() *Server {
	return &Server{}
}

func TestGetAritst(t *testing.T) {
	s := InitTestServer()

	resp, err := s.GetArtist(context.Background(), &pb.GetArtistRequest{Artist: "The Beatles"})
	if err != nil {
		t.Fatalf("Failure to get artist: %v", err)
	}

	if resp.GetArtistSort() != "Beatles, The" {
		t.Errorf("Wrong artist sort returned '%v' -> should have been 'Beatles, The'", resp.GetArtistSort())
	}
}
