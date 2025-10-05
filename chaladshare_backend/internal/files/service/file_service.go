package service

import (
	"fmt"

	"chaladshare_backend/internal/files/models"
	"chaladshare_backend/internal/files/repository"
)

type FileService interface {
	UploadFile(req *models.UploadRequest) (*models.UploadResponse, error) //upload file
	GetFilesByUserID(userID int) ([]models.Document, error)               //ดึงไฟล์สรุปของผู้ใช้
	SaveSummary(summary *models.Summary) error                            //บันทึกสรุปที่ได้มาจาก AI
	GetSummaryByDocumentID(docID int) (*models.Summary, error)            //ดึงสรุปจาก document_id
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo: repo}
}

// 1.อัปโหลดไฟล์
func (s *fileService) UploadFile(req *models.UploadRequest) (*models.UploadResponse, error) {
	if req.DocumentName == "" || req.DocumentURL == "" {
		return nil, fmt.Errorf("กรุณาใส่ชื่อไฟล์และ URL ให้ครบถ้วน")
	}

	//สร้าง strct doc เพื่อเตรียมบันทึกลง DB
	doc := &models.Document{
		DocumentUserID:  req.UserID,
		DocumentName:    req.DocumentName,
		DocumentURL:     req.DocumentURL,
		StorageProvider: req.Storage,
	}

	//repo เพื่อบันทึกไฟล์
	if err := s.repo.SaveDocument(doc); err != nil {
		return nil, err
	}

	//response to clint
	resp := &models.UploadResponse{
		Message: "อัปโหลดไฟล์สำเร็จ",
		File:    *doc,
	}
	return resp, nil
}

// 2.ดึงข้อมูลไฟล์ทั้งหมดของผู้ใช้
func (s *fileService) GetFilesByUserID(userID int) ([]models.Document, error) {
	doc, err := s.repo.GetDocumentByID(userID)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลไฟล์ได้: %v", err)
	}
	return doc, nil
}

// 3.บันทึกสรุปที่ได้มาจาก AI
func (s *fileService) SaveSummary(summary *models.Summary) error {
	if summary.DocumentID == 0 {
		return fmt.Errorf("ต้องระบุ ID ของไฟล์ก่อนบันทึกสรุป")
	}
	if summary.SummaryText == "" {
		return fmt.Errorf("ไม่มีข้อความสรุปให้บันทึก")
	}
	if err := s.repo.SaveSummary(summary); err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกสรุปได้: %v", err)
	}
	return nil
}

// 4.ดึงสรุปจาก document_id
func (s *fileService) GetSummaryByDocumentID(docID int) (*models.Summary, error) {
	summary, err := s.repo.GetSummaryByDocumentID(docID)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลสรุปของไฟล์นี้ได้: %v", err)
	}
	return summary, nil
}
