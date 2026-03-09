package binman

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

// extractTarXz extracts binaryName from a tar.xz archive, returning the temp path.
func extractTarXz(archivePath, binaryName string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	xzReader, err := xz.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("failed to create xz reader: %w", err)
	}

	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar entry: %w", err)
		}

		entryName := filepath.Base(header.Name)
		if entryName == binaryName && header.Typeflag == tar.TypeReg {
			tmpFile, err := os.CreateTemp("", binaryName+"-extracted-*")
			if err != nil {
				return "", fmt.Errorf("failed to create temp file: %w", err)
			}

			if _, err := io.Copy(tmpFile, tarReader); err != nil {
				tmpFile.Close()
				os.Remove(tmpFile.Name())
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			tmpFile.Close()
			return tmpFile.Name(), nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found in archive", binaryName)
}

// extractZip extracts binaryName from a ZIP archive, returning the temp path.
func extractZip(archivePath, binaryName string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		entryName := filepath.Base(f.Name)
		if entryName == binaryName {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open zip entry: %w", err)
			}
			defer rc.Close()

			tmpFile, err := os.CreateTemp("", binaryName+"-extracted-*")
			if err != nil {
				return "", fmt.Errorf("failed to create temp file: %w", err)
			}

			if _, err := io.Copy(tmpFile, rc); err != nil {
				tmpFile.Close()
				os.Remove(tmpFile.Name())
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			tmpFile.Close()
			return tmpFile.Name(), nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found in zip archive", binaryName)
}
