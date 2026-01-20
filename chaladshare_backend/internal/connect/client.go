package connect

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client

	// timeout แยกตามงาน
	ExtractTimeout   time.Duration
	SummarizeTimeout time.Duration
}

func NewFromEnv() (*Client, error) {
	base := strings.TrimRight(os.Getenv("COLAB_URL"), "/")
	if base == "" {
		return nil, fmt.Errorf("COLAB_URL is empty")
	}

	key := os.Getenv("COLAB_API_KEY")

	return &Client{
		BaseURL:          base,
		APIKey:           key,
		HTTP:             &http.Client{},
		ExtractTimeout:   180 * time.Second, // เท่าของเดิม
		SummarizeTimeout: 10 * time.Minute,  // summarize นานกว่า
	}, nil
}

func (c *Client) postPDF(ctx context.Context, endpoint string, documentID int, pdfPath string) (*http.Response, error) {
	url := c.BaseURL + endpoint

	f, err := os.Open(pdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	// เขียน multipart ใน goroutine เพื่อ stream ไฟล์ (ไม่กิน RAM)
	go func() {
		defer pw.Close()
		defer w.Close()

		_ = w.WriteField("document_id", strconv.Itoa(documentID))
		_ = w.WriteField("file_name", filepath.Base(pdfPath))

		fw, err := w.CreateFormFile("pdf", filepath.Base(pdfPath))
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(fw, f); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, "POST", url, pr)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("ngrok-skip-browser-warning", "true")

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	return c.HTTP.Do(req)
}
