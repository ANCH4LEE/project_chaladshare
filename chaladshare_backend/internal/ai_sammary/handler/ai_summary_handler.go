package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AISummaryHandler struct {
	colabURL string
	apiKey   string
	client   *http.Client
}

func NewAISummaryHandler() *AISummaryHandler {
	return &AISummaryHandler{
		colabURL: os.Getenv("COLAB_URL"),
		apiKey:   os.Getenv("COLAB_API_KEY"),
		client: &http.Client{
			Timeout: 10 * time.Minute, // กันสรุปนาน
		},
	}
}

func (h *AISummaryHandler) Summarize(c *gin.Context) {
	if h.colabURL == "" {
		c.JSON(500, gin.H{"error": "COLAB_URL is not set"})
		return
	}

	// รับไฟล์จาก React: form-data key = "file"
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "missing file (key: file)"})
		return
	}

	src, err := fh.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot open uploaded file"})
		return
	}
	defer src.Close()

	// สร้าง multipart ส่งต่อไป Colab
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", fh.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot create multipart"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "cannot write multipart"})
		return
	}
	writer.Close()

	req, err := http.NewRequestWithContext(c.Request.Context(), "POST", h.colabURL, &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot create request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if h.apiKey != "" {
		req.Header.Set("X-API-KEY", h.apiKey)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to call colab", "detail": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// forward status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.Data(resp.StatusCode, "application/json; charset=utf-8", body)
		return
	}

	var js map[string]any
	if err := json.Unmarshal(body, &js); err != nil {
		// ถ้าไม่ใช่ json ก็ส่ง raw กลับ
		c.Data(200, "text/plain; charset=utf-8", body)
		return
	}
	c.JSON(200, js)
}
