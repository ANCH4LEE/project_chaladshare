package repository

import (
	"database/sql"
	"fmt"
	"time"

	"chaladshare_backend/internal/files/models"
)

type FileRepository interface {
	// documents
	CreateDocument(doc *models.Document, images []string) (*models.Document, error)
	GetDocumentByUserID(userID int) ([]models.Document, error)
	DeleteDocument(id int) error

	// document pages
	GetDocumentPages(docID int) ([]models.DocumentPage, error)

	// summaries
	GetSummaryByDocID(docID int) (*models.Summary, error)
	CreateSummary(summary *models.Summary) (*models.Summary, error)
}

type fileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &fileRepository{db: db}
}

// CreateDocument
func (r *fileRepository) CreateDocument(req *models.Document, images []string) (*models.Document, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO documents (document_user_id, document_name, document_url, storage_provider, page_count, uploaded_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING document_id, uploaded_at
	`,
		req.DocumentUserID, req.DocumentName, req.DocumentURL, req.StorageProvider, len(images), time.Now(),
	).Scan(&req.DocumentID, &req.UploadedAt)

	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถบันทึกไฟล์ได้: %v", err)
	}

	// บันทึกภาพแต่ละหน้า
	for i, img := range images {
		_, err := tx.Exec(`
			INSERT INTO document_pages (document_id, page_index, image_url)
			VALUES ($1,$2,$3)
		`, req.DocumentID, i+1, img)
		if err != nil {
			return nil, fmt.Errorf("ไม่สามารถเพิ่มภาพหน้า %d: %v", i+1, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return req, nil
}

// GetDocumentByUserID
func (r *fileRepository) GetDocumentByUserID(userID int) ([]models.Document, error) {
	rows, err := r.db.Query(`
		SELECT document_id, document_user_id, document_name, document_url, storage_provider, page_count, uploaded_at
		FROM documents
		WHERE document_user_id = $1
		ORDER BY uploaded_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(&d.DocumentID, &d.DocumentUserID, &d.DocumentName, &d.DocumentURL, &d.StorageProvider, &d.PageCount, &d.UploadedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// GetDocumentPages แต่ละหน้า
func (r *fileRepository) GetDocumentPages(docID int) ([]models.DocumentPage, error) {
	rows, err := r.db.Query(`
		SELECT doc_page_id, document_id, page_index, image_url, created_at
		FROM document_pages
		WHERE document_id = $1
		ORDER BY page_index ASC
	`, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []models.DocumentPage
	for rows.Next() {
		var p models.DocumentPage
		if err := rows.Scan(&p.DocPageID, &p.DocumentID, &p.PageIndex, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

// CreateSummary
func (r *fileRepository) CreateSummary(summary *models.Summary) (*models.Summary, error) {
	err := r.db.QueryRow(`
		INSERT INTO summaries (summary_text, summary_html, summary_pdf_url, summary_created_at, document_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING summary_id, summary_created_at
	`, summary.SummaryText, summary.SummaryHTML, summary.SummaryPDFURL, time.Now(), summary.DocumentID).
		Scan(&summary.SummaryID, &summary.SummaryCreatedAt)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

// GetSummaryByDocID
func (r *fileRepository) GetSummaryByDocID(docID int) (*models.Summary, error) {
	var s models.Summary
	err := r.db.QueryRow(`
		SELECT summary_id, summary_text, summary_html, summary_pdf_url, summary_created_at, document_id
		FROM summaries
		WHERE document_id = $1
	`, docID).Scan(&s.SummaryID, &s.SummaryText, &s.SummaryHTML, &s.SummaryPDFURL, &s.SummaryCreatedAt, &s.DocumentID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ไม่พบสรุปของเอกสารนี้")
	}
	return &s, err
}

// DeleteDocument
func (r *fileRepository) DeleteDocument(id int) error {
	res, err := r.db.Exec("DELETE FROM documents WHERE document_id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
