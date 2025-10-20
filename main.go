package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func checkLatest() (string, error) {
	res, err := http.Get("https://data.metabrainz.org/pub/musicbrainz/data/fullexport/LATEST")
	if err != nil {
		return "", fmt.Errorf("unable to get data: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}

	return strings.TrimSpace(string(resBody)), nil
}

func checkDBVersion() (string, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DBNAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	return "", nil
}

func runLoop() {
	for {
		log.Printf("Checking musicbrainz")

		val, err := checkLatest()
		log.Printf("Versioned returned %v -> %v", val, err)

		dbvla, err := checkDBVersion()
		log.Printf("DB version returned %v -> %v", dbvla, err)

		time.Sleep(time.Hour)
	}
}

func main() {
	runLoop()
}
