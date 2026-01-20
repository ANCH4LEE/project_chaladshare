package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

type SummarizeResp struct {
	DocumentID    int    `json:"document_id"`
	SummaryText   string `json:"summary_text"`
	SummaryHTML   string `json:"summary_html,omitempty"`
	SummaryPDFURL string `json:"summary_pdf_url,omitempty"`
}

func (c *Client) Summarize(documentID int, pdfPath string) (*SummarizeResp, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.SummarizeTimeout)
	defer cancel()

	resp, err := c.postPDF(ctx, "/summarize", documentID, pdfPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("summarize status %d: %s", resp.StatusCode, string(b))
	}

	var out SummarizeResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	log.Printf("[COLAB][SUM] OK time=%s doc=%d len(summary)=%d",
		time.Since(start), out.DocumentID, len(out.SummaryText))

	return &out, nil
}
