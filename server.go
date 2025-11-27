package main

import "database/sql"

type Server struct {
	db      *sql.DB
	version string
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
