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

	file := fmt.Sprintf("%v%v/mbdump.tar.bz2", BASE, latest)
	log.Printf("Downloading %v", file)
	resp, err := http.Get(file)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}
