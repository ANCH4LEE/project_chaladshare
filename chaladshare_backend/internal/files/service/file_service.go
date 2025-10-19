package service

import (
	"errors"
	"fmt"
	"strings"

	"chaladshare_backend/internal/files/models"
	"chaladshare_backend/internal/files/repository"
)

type FileService interface {
	UploadFile(req *models.UploadRequest) (*models.UploadResponse, error)
	GetFilesByUserID(userID int) ([]models.Document, error)
	DeleteFile(documentID int) error

	GetDocumentOwnerID(documentID int) (int, error)
	GetDocumentPages(docID int) ([]models.DocumentPage, error)

	SaveSummary(summary *models.Summary) (*models.Summary, error)
	GetSummaryByDocumentID(docID int) (*models.Summary, error)

	IsOwner(documentID int, userID int) (bool, error)
}

type fileService struct {
	filerepo repository.FileRepository
}

func NewFileService(filerepo repository.FileRepository) FileService {
	return &fileService{filerepo: filerepo}
}

func (s *fileService) UploadFile(req *models.UploadRequest) (*models.UploadResponse, error) {
	if strings.TrimSpace(req.DocumentName) == "" {
		return nil, errors.New("ต้องระบุชื่อไฟล์")
	}
	if strings.TrimSpace(req.DocumentURL) == "" {
		return nil, errors.New("ต้องระบุ URL ของไฟล์")
	}
	if len(req.Images) == 0 {
		return nil, errors.New("ต้องมีภาพของไฟล์อย่างน้อย 1 หน้า")
	}

	doc := &models.Document{
		DocumentUserID:  req.UserID,
		DocumentName:    req.DocumentName,
		DocumentURL:     req.DocumentURL,
		StorageProvider: req.StorageProvider,
	}

	savedDoc, err := s.filerepo.CreateDocument(doc, req.Images)
	if err != nil {
		return nil, fmt.Errorf("บันทึกไฟล์ไม่สำเร็จ: %v", err)
	}

	resp := &models.UploadResponse{
		Message:    "อัปโหลดไฟล์สำเร็จ",
		File:       *savedDoc,
		FileURL:    savedDoc.DocumentURL,
		DocumentID: savedDoc.DocumentID,
	}
	return resp, nil
}

func (s *fileService) GetFilesByUserID(userID int) ([]models.Document, error) {
	files, err := s.filerepo.GetDocumentByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลไฟล์ได้: %v", err)
	}
	return files, nil
}

func (s *fileService) DeleteFile(documentID int) error {
	if documentID <= 0 {
		return errors.New("document_id ไม่ถูกต้อง")
	}

	if err := s.filerepo.DeleteDocument(documentID); err != nil {
		return fmt.Errorf("ไม่สามารถลบไฟล์ได้: %v", err)
	}
	return nil
}

func (s *fileService) GetDocumentPages(docID int) ([]models.DocumentPage, error) {
	if docID <= 0 {
		return nil, errors.New("document_id ไม่ถูกต้อง")
	}

	pages, err := s.filerepo.GetDocumentPages(docID)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงภาพของไฟล์ได้: %v", err)
	}
	return pages, nil
}

func (s *fileService) SaveSummary(summary *models.Summary) (*models.Summary, error) {
	if summary.DocumentID == 0 {
		return nil, errors.New("ต้องระบุ document_id")
	}
	if strings.TrimSpace(summary.SummaryText) == "" {
		return nil, errors.New("ต้องมีข้อความสรุปก่อนบันทึก")
	}

	saved, err := s.filerepo.CreateSummary(summary)
	if err != nil {
		return nil, fmt.Errorf("บันทึกสรุปไม่สำเร็จ: %v", err)
	}
	return saved, nil
}

func (s *fileService) GetSummaryByDocumentID(docID int) (*models.Summary, error) {
	if docID <= 0 {
		return nil, errors.New("document_id ไม่ถูกต้อง")
	}

	summary, err := s.filerepo.GetSummaryByDocID(docID)
	if err != nil {
		return nil, fmt.Errorf("ไม่พบสรุปของไฟล์นี้: %v", err)
	}
	return summary, nil
}

// ดึง owner_id ของเอกสารจาก repository
func (s *fileService) GetDocumentOwnerID(documentID int) (int, error) {
	if documentID <= 0 {
		return 0, errors.New("document_id ไม่ถูกต้อง")
	}
	ownerID, err := s.filerepo.GetDocumentOwnerID(documentID)
	if err != nil {
		return 0, fmt.Errorf("ตรวจสอบเจ้าของไฟล์ล้มเหลว: %v", err)
	}
	return ownerID, nil
}

// เช็คว่า userID เป็นเจ้าของไฟล์ documentID หรือไม่
func (s *fileService) IsOwner(documentID int, userID int) (bool, error) {
	if documentID <= 0 || userID <= 0 {
		return false, errors.New("document_id หรือ user_id ไม่ถูกต้อง")
	}
	ownerID, err := s.filerepo.GetDocumentOwnerID(documentID)
	if err != nil {
		return false, fmt.Errorf("ตรวจสอบเจ้าของไฟล์ล้มเหลว: %v", err)
	}
	return ownerID == userID, nil
}
