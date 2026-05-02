package main

import (
	"database/sql"
	"sync"
)

type Server struct {
	db           *sql.DB
	version      string
	activeIssues map[string]bool
	mu           sync.Mutex
	githubridge  *GithubridgeClient
}

type GithubridgeClient struct{}

func (g *GithubridgeClient) PostIssue(artist, title, body string) error {
    return nil
}

type ConfiguredServer interface {
	CheckLatest() (string, error)
}

type TestServer struct {
	base *Server
}

func (ts *TestServer) CheckLatest() (string, error) {
	return "1", nil
}
