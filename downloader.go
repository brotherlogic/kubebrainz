package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BASE = "https://data.metabrainz.org/pub/musicbrainz/data/fullexport/"
)

func getLatest() (string, error) {
	// Create a new bytes.Buffer
	var buf bytes.Buffer
	resp, err := http.Get(fmt.Sprintf("%vLATEST", BASE))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func downloadFile() error {
	latest, err := getLatest()
	if err != nil {
		return err
	}

	out, err := os.Create("download.tar.bz2")
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("Downloading %v", fmt.Sprintf("%v%v/mdbdump.tar.bz2", BASE, latest))
	resp, err := http.Get(fmt.Sprintf("%v%v/mdbdump.tar.bz2", BASE, latest))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}
