package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type colabResp struct {
	DocumentID     int       `json:"document_id"`
	StyleLabel     *string   `json:"style_label"`
	StyleVectorV16 []float64 `json:"style_vector_v16"`
	StyleVector    []float64 `json:"style_vector,omitempty"`
	ContentText    string    `json:"content_text"`
	Embedding      []float64 `json:"content_embedding"`
	ClusterID      *int      `json:"cluster_id"`
}

func callColabExtract(documentID int, pdfPath string) (*colabResp, error) {
	base := strings.TrimRight(os.Getenv("COLAB_URL"), "/")
	if base == "" {
		return nil, fmt.Errorf("COLAB_URL is empty")
	}

	url := base + "/extract_features"
	log.Printf("[COLAB] -> POST %s document_id=%d pdf=%s", url, documentID, pdfPath)
	start := time.Now()

	f, err := os.Open(pdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer w.Close()

		_ = w.WriteField("document_id", strconv.Itoa(documentID))

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

	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	req.Header.Set("ngrok-skip-browser-warning", "true")

	if key := os.Getenv("COLAB_API_KEY"); key != "" {
		req.Header.Set("X-API-Key", key)
	}

	httpc := &http.Client{Timeout: 180 * time.Second}
	resp, err := httpc.Do(req)
	if err != nil {
		log.Printf("[COLAB] !! request error after %s: %v", time.Since(start), err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("[COLAB] <- status=%d time=%s body=%s", resp.StatusCode, time.Since(start), string(b))
		return nil, fmt.Errorf("colab status %d: %s", resp.StatusCode, string(b))
	}

	var out colabResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Printf("[COLAB] !! decode error time=%s: %v", time.Since(start), err)
		return nil, err
	}

	vec := out.StyleVectorV16
	if len(vec) == 0 && len(out.StyleVector) > 0 {
		vec = out.StyleVector
	}
	if vec == nil {
		vec = []float64{}
	}
	out.StyleVectorV16 = vec

	log.Printf("[COLAB] <- OK status=%d time=%s label=%v vec_len=%d cluster=%v",
		resp.StatusCode, time.Since(start), out.StyleLabel, len(vec), out.ClusterID)

	return &out, nil

}
