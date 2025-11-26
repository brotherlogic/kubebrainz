package main

import (
	"context"
	"database/sql"
	"log"
	"testing"

	pb "github.com/brotherlogic/kubebrainz/proto"
	"github.com/stapelberg/postgrestest"
)

var pgt *postgrestest.Server

func TestMain(m *testing.M) {
	var err error
	pgt, err = postgrestest.Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer pgt.Cleanup()

	m.Run()
}

func InitTestServer() (*Server, error) {
	pgurl, err := pgt.CreateDatabase(context.Background())
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", pgurl)
	if err != nil {
		return nil, err
	}

	s := &Server{db: db}
	err = s.initDB()
	if err != nil {
		return nil, err
	}

	err = s.loadFile(context.Background(), "artist", "testdata/artist-test.sql")

	return s, err
}

func DoubleInitTestServer() (*Server, error) {
	pgurl, err := pgt.CreateDatabase(context.Background())
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", pgurl)
	if err != nil {
		return nil, err
	}

	s := &Server{db: db}
	err = s.initDB()
	if err != nil {
		return nil, err
	}

	err = s.loadFile(context.Background(), "artist", "testdata/artist-test.sql")

	if err != nil {
		return nil, err
	}

	err = s.initDB()

	return s, err
}

func TestGetAritst(t *testing.T) {
	s, err := InitTestServer()
	if err != nil {
		log.Fatalf("Unable to init test server: %v", err)
	}

	resp, err := s.GetArtist(context.Background(), &pb.GetArtistRequest{Artist: "The Beatles"})
	if err != nil {
		t.Fatalf("Failure to get artist: %v", err)
	}

	if resp.GetArtistSort() != "Beatles, The" {
		t.Errorf("Wrong artist sort returned '%v' -> should have been 'Beatles, The'", resp.GetArtistSort())
	}
}
func TestGetAritstWithDouble(t *testing.T) {
	s, err := DoubleInitTestServer()
	if err != nil {
		log.Fatalf("Unable to init test server: %v", err)
	}

	resp, err := s.GetArtist(context.Background(), &pb.GetArtistRequest{Artist: "The Beatles"})
	if err != nil {
		t.Fatalf("Failure to get artist: %v", err)
	}

	if resp.GetArtistSort() != "Beatles, The" {
		t.Errorf("Wrong artist sort returned '%v' -> should have been 'Beatles, The'", resp.GetArtistSort())
	}
}
