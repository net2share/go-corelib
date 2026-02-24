package binman

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// httpDownload downloads a URL to a temp file, reporting progress.
func httpDownload(url string, fn ProgressFunc) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %s for %s", resp.Status, url)
	}

	tmpFile, err := os.CreateTemp("", "binman-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	var reader io.Reader = resp.Body
	if fn != nil {
		reader = &progressReader{
			reader:     resp.Body,
			total:      resp.ContentLength,
			progressFn: fn,
		}
	}

	written, err := io.Copy(tmpFile, reader)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if written == 0 {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("downloaded file is empty")
	}

	return tmpFile.Name(), nil
}

type progressReader struct {
	reader     io.Reader
	total      int64
	downloaded int64
	progressFn ProgressFunc
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.downloaded += int64(n)
	if pr.progressFn != nil {
		pr.progressFn(pr.downloaded, pr.total)
	}
	return n, err
}
