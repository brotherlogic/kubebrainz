package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/kubebrainz/proto"
)

var (
	port        = flag.Int("port", 8080, "The server port.")
	metricsPort = flag.Int("mmetricsport", 8081, "Serves prometheus metrics")
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

func (s *Server) runLoop() {
	for {
		log.Printf("Checking musicbrainz")

		val, err := checkLatest()
		log.Printf("Versioned returned %v -> %v", val, err)

		dbvla, err := checkDBVersion()
		log.Printf("DB version returned %v -> %v", dbvla, err)

		time.Sleep(time.Hour)
	}
}

func (s *Server) checkCount(ctx context.Context) int {
	res, err := s.db.QueryContext(ctx, "SELECT COUNT(*) FROM artist")
	if err != nil {
		log.Printf("error counting artist: %v", err)
		return 0
	}
	res.Next()
	var count int
	res.Scan(&count)
	return count
}

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DBNAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("")
	}
	defer db.Close()

	s := &Server{
		db: db,
	}
	err = s.initDB()
	if err != nil {
		log.Fatalf("unable to init db: %v", err)
	}

	count := s.checkCount(context.Background())

	if count == 0 {
		t := time.Now()
		// Download on startup
		version, err := downloadFile()
		if err != nil {
			log.Fatalf("unable to download file: %v", err)
		}
		s.version = version
		fi, err := os.Stat("download.tar.bz2")
		if err != nil {
			log.Fatalf("Unable to stat downloadedfile: %v", err)
		}
		// get the size
		size := fi.Size()
		log.Printf("Downloaded in %v (%v)", time.Since(t), size)

		if size < 1000 {
			data, err := os.ReadFile("download.tar.bz2")
			if err != nil {
				log.Fatalf("Error reading file: %v", err)
				return
			}

			// Convert the byte slice to a string and print it
			log.Println(string(data))
		}

		t2 := time.Now()
		err = s.loadDatabase(context.Background(), "download.tar.bz2")
		if err != nil {
			log.Fatalf("Unable to load database: %v", err)
		}
		log.Printf("Database loaded in %v", time.Since(t2))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("pstore failed to listen on the serving port %v: %v", *port, err)
	}

	gs := grpc.NewServer()
	pb.RegisterKubeBrainzServiceServer(gs, s)
	log.Printf("kubebrainz is listening on %v", lis.Addr())

	// Setup prometheus export
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	}()

	//go s.runLoop()

	go func() {
		time.Sleep(time.Minute * 5)
		res, err := s.GetArtist(context.Background(), &pb.GetArtistRequest{
			Artist: "The Beatles",
		})
		if err != nil {
			log.Fatalf("Unable to get artist: %v", err)
		}
		log.Printf("Artist: %v", res)
	}()

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("kubebrainz failed to serve: %v", err)
	}

}
