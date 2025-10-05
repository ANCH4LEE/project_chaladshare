package repository

import (
	"database/sql"
	"fmt"
	"time"

	"chaladshare_backend/internal/files/models"
)

type FileRepository interface {
	SaveDocument(doc *models.Document) error
	GetDocumentByID(userID int) ([]models.Document, error)
	SaveSummary(summary *models.Summary) error
	GetSummaryByDocumentID(docID int) (*models.Summary, error)
}

type fileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &fileRepository{db: db}
}

// บันทึกข้อมูลไฟล์ PDF ที่อัปโหลด
func (r *fileRepository) SaveDocument(doc *models.Document) error {
	query := `
		INSERT INTO documents(document_user_id, document_name, document_url, storage_provider, upload_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING document_id
	`
	err := r.db.QueryRow(
		query,
		doc.DocumentUserID,
		doc.DocumentName,
		doc.DocumentURL,
		doc.StorageProvider,
		time.Now(),
	).Scan(&doc.DocumentUserID, &doc.UploadedAt)

	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกไฟล์ได้: %v", err)
	}
	return nil
}

// ดึงข้อมูลไฟล์ PDF ทั้งหมดของผู้ใช้
func (r *fileRepository) GetDocumentByID(userID int) ([]models.Document, error) {
	query := `
		SELECT document_id, document_user_id, document_name, document_url, storage_provider, updated_at
		FROM documents
		WHERE document_user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลไฟล์ได้: %v", err)
	}
	defer rows.Close()

	var doc []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(
			&d.DocumentID,
			&d.DocumentUserID,
			&d.DocumentName,
			&d.DocumentURL,
			&d.StorageProvider,
			&d.UploadedAt,
		); err != nil {
			return nil, err
		}
		doc = append(doc, d)
	}

	return doc, nil
}

// บันทึกข้อมูลสรุปเนื้อหาที่ได้จาก AI
func (r *fileRepository) SaveSummary(summary *models.Summary) error {
	query := `
		INSERT INTO summaries (summary_text, summary_html, summary_pdf_url, summary_created_at, document_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING summary_id, summary_created_at;
	`

	err := r.db.QueryRow(
		query,
		summary.SummaryText,
		summary.SummaryHTML,
		summary.SummaryPDFURL,
		time.Now(),
		summary.DocumentID,
	).Scan(&summary.SummaryID, &summary.SummaryCreatedAt)

	if err != nil {
		return fmt.Errorf("ไม่สามารถเพิ่มสรุปได้: %v", err)
	}
	return nil
}

// ดึงข้อมูลสรุป DocumentID
func (r *fileRepository) GetSummaryByDocumentID(docID int) (*models.Summary, error) {
	query := `
		SELECT summary_id, summary_text, summary_html, summary_pdf_url, summary_created_at, document_id
		FROM summaries
		WHERE document_id = $1
		LIMIT 1
	`

	var s models.Summary
	err := r.db.QueryRow(query, docID).Scan(
		&s.SummaryID,
		&s.SummaryText,
		&s.SummaryHTML,
		&s.SummaryPDFURL,
		&s.SummaryCreatedAt,
		&s.DocumentID,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ไม่พบข้อมูลสรุปของไฟล์นี้")
	}
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลสรุปได้: %v", err)
	}

	return &s, nil
}
