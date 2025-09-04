package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

func runLoop() {
	for {
		log.Printf("Checking musicbrainz")

		val, err := checkLatest()
		log.Printf("Returned %v -> %v", val, err)

		time.Sleep(time.Hour)
	}
}

func main() {
	runLoop()
}
