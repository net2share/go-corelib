package binman

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// fetchAndVerifyChecksum downloads a checksum file and verifies the target file.
// Supports both multi-file checksum files ("hash  filename" lines) and
// single-hash files (e.g. shadowsocks per-file .sha256).
func fetchAndVerifyChecksum(checksumURL, filePath, expectedFilename string) error {
	resp, err := http.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("failed to fetch checksum from %s: %w", checksumURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch checksum: %s", resp.Status)
	}

	var expectedHash string
	scanner := bufio.NewScanner(resp.Body)
	lineCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineCount++

		// Xray .dgst format: "ALGO= hash" (e.g. "SHA2-256= abc123...")
		if strings.Contains(line, "= ") && !strings.Contains(line, "  ") {
			parts := strings.SplitN(line, "= ", 2)
			if len(parts) == 2 {
				algo := strings.TrimSpace(parts[0])
				hash := strings.TrimSpace(parts[1])
				// Prefer SHA2-256 / SHA-256
				if strings.Contains(strings.ToUpper(algo), "SHA2-256") ||
					strings.Contains(strings.ToUpper(algo), "SHA-256") ||
					strings.Contains(strings.ToUpper(algo), "SHA256") {
					expectedHash = hash
					break
				}
				// Fall back to first hash if no SHA-256 found yet
				if expectedHash == "" {
					expectedHash = hash
				}
				continue
			}
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// Standard format: "hash  filename"
			hash := parts[0]
			filename := parts[len(parts)-1]
			if filename == expectedFilename || strings.HasSuffix(filename, expectedFilename) {
				expectedHash = hash
				break
			}
		} else if len(parts) == 1 && lineCount == 1 {
			// Single-hash file (e.g. per-file .sha256)
			expectedHash = parts[0]
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to parse checksum file: %w", err)
	}

	if expectedHash == "" {
		return fmt.Errorf("checksum for %s not found in %s", expectedFilename, checksumURL)
	}

	return verifyFileSHA256(filePath, expectedHash)
}

// verifyFileSHA256 computes the SHA256 of a file and compares it to the expected hash.
func verifyFileSHA256(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expectedHash) {
		return fmt.Errorf("SHA256 checksum mismatch: expected %s, got %s", expectedHash, actual)
	}

	return nil
}
