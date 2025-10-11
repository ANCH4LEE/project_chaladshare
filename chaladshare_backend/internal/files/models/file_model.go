package models

import "time"

// ข้อมูลไฟล์ที่อัปโหลด
type Document struct {
	DocumentID      int       `json:"document_id"`
	DocumentUserID  int       `json:"document_user_id"`
	DocumentName    string    `json:"document_name"`
	DocumentURL     string    `json:"document_url"`
	StorageProvider string    `json:"storage_provider"`
	UploadedAt      time.Time `json:"uploaded_at"`
	PageCount       int       `json:"page_count"`
}

//
type DocumentPage struct {
	DocPageID  int       `json:"doc_page_id"`
	DocumentID int       `json:"document_id"`
	PageIndex  int       `json:"page_index"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
}

// เก็บข้อมูลจากไฟล์ที่สรุปเนื้อหาด้วย AI
type Summary struct {
	SummaryID        int       `json:"summary_id"`
	SummaryText      string    `json:"summary_text"`
	SummaryHTML      string    `json:"summary_html"`
	SummaryPDFURL    string    `json:"summary_pdf_url"`
	SummaryCreatedAt time.Time `json:"summary_created_at"`
	DocumentID       int       `json:"document_id"`
}

// ผู้ใช้อัปโหลดไฟล์
type UploadFrom struct {
	UserID int `form:"user_id"`
}

type UploadRequest struct {
	UserID          int      `json:"user_id" binding:"required"`
	DocumentName    string   `json:"document_name" binding:"required"`
	DocumentURL     string   `json:"document_url" binding:"required"`
	StorageProvider string   `json:"storage_provider" binding:"required"`
	Images          []string `json:"images" binding:"required"` // path ของภาพแต่ละหน้า
}

// response ตอนอัปโหลดสำเร็จ
type UploadResponse struct {
	Message    string   `json:"message"`
	File       Document `json:"file"`
	FileURL    string   `json:"file_url"`
	DocumentID int      `json:"document_id"`
}
