package main

import (
	"archive/tar"
	"compress/bzip2"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const tables = `
CREATE TABLE artist ( -- replicate (verbose)
    id                  SERIAL,
    gid                 UUID NOT NULL,
    name                VARCHAR NOT NULL,
    sort_name           VARCHAR NOT NULL,
    begin_date_year     SMALLINT,
    begin_date_month    SMALLINT,
    begin_date_day      SMALLINT,
    end_date_year       SMALLINT,
    end_date_month      SMALLINT,
    end_date_day        SMALLINT,
    type                INTEGER, -- references artist_type.id
    area                INTEGER, -- references area.id
    gender              INTEGER, -- references gender.id
    comment             VARCHAR(255) NOT NULL DEFAULT '',
    edits_pending       INTEGER NOT NULL DEFAULT 0 CHECK (edits_pending >= 0),
    last_updated        TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended               BOOLEAN NOT NULL DEFAULT FALSE
      CONSTRAINT artist_ended_check CHECK (
        (
          -- If any end date fields are not null, then ended must be true
          (end_date_year IS NOT NULL OR
           end_date_month IS NOT NULL OR
           end_date_day IS NOT NULL) AND
          ended = TRUE
        ) OR (
          -- Otherwise, all end date fields must be null
          (end_date_year IS NULL AND
           end_date_month IS NULL AND
           end_date_day IS NULL)
        )
      ),
    begin_area          INTEGER, -- references area.id
    end_area            INTEGER -- references area.id
);
`

func (s *Server) unzipFile(archivePath, outputPath string) error {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return err
	}

	// Open the .tar.bz2 file
	file, err := os.Open(archivePath)
	if err != nil {
		fmt.Printf("Error opening archive: %v\n", err)
		return err
	}
	defer file.Close()

	// Create a bzip2 decompressor
	bz2Reader := bzip2.NewReader(file)

	// Create a tar reader from the decompressed stream
	tarReader := tar.NewReader(bz2Reader)

	// Iterate through the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			fmt.Printf("Error reading tar header: %v\n", err)
			return err
		}

		// Clean and validate the extracted file path to prevent Zip Slip (directory traversal)
		cleanHeaderName := filepath.Clean(header.Name)
		// Reject absolute paths or traversal outside the target dir
		if strings.HasPrefix(cleanHeaderName, ".."+string(os.PathSeparator)) || cleanHeaderName == ".." || filepath.IsAbs(header.Name) {
			fmt.Printf("Skipping potentially unsafe file outside target dir: %s\n", header.Name)
			continue
		}
		targetPath := filepath.Join(outputPath, cleanHeaderName)
		absOutputPath, err := filepath.Abs(outputPath)
		if err != nil {
			fmt.Printf("Error getting absolute output path: %v\n", err)
			return err
		}
		absTargetPath, err := filepath.Abs(targetPath)
		if err != nil {
			fmt.Printf("Error getting absolute target path: %v\n", err)
			return err
		}
		if !strings.HasPrefix(absTargetPath, absOutputPath+string(os.PathSeparator)) && absTargetPath != absOutputPath {
			fmt.Printf("Skipping potentially unsafe file outside target dir: %s\n", header.Name)
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", targetPath, err)
				return err
			}
		case tar.TypeReg:
			// Create file
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", targetPath, err)
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Printf("Error writing file %s: %v\n", targetPath, err)
				return err
			}
		default:
			fmt.Printf("Skipping unsupported tar entry type: %v for %s\n", header.Typeflag, header.Name)
		}
	}

	fmt.Printf("Successfully extracted %s to %s\n", archivePath, outputPath)

	return nil
}

func (s *Server) loadDatabase(ctx context.Context, file string) error {
	err := s.initDB()
	if err != nil {
		return err
	}

	//Unzip. the tarball
	err = s.unzipFile(file, "data_out")
	if err != nil {
		return err
	}

	return s.loadFile(ctx, "artist", "data_out/artist")
}

func (s *Server) initDB() error {
	_, err := s.db.Exec(tables)
	return err
}

func (s *Server) loadFile(ctx context.Context, table string, file string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, fmt.Sprintf("COPY %v FROM '%v/%v' (DELIMITER('\t'))", table, path, file))
	return err
}
