package service

import (
	"fmt"
	// "log"
	// "os"

	// "path/filepath"

	// "github.com/pdfcpu/pdfcpu"
	// "github.com/pdfcpu/pdfcpu/pkg/api"

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

// 1.อัปโหลดไฟล์ พร้อมแปลงภาพหลังบันทึกได้
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

	// แปลง PDF เป็น ภาพ
	// pdfPath := fmt.Sprintf("./uploads/%s", req.DocumentName)
	// outputDir := fmt.Sprintf("./uploads/pages/%d", doc.DocumentID)

	// pageURLs, err := convertPDFToImages(pdfPath, outputDir)
	// if err != nil {
	// 	log.Printf("[convert] error: %v", err)
	// 	pageURLs = []string{} // ถ้าแปลงไม่ได้ให้ยังอัปโหลดผ่าน
	// }

	// เก็บข้อมูลภาพและหน้า
	// thumbnailURL := ""
	// pageCount := len(pageURLs)
	// if len(pageURLs) > 0 {
	// 	thumbnailURL = pageURLs[0]
	// }

	// อัปเดตข้อมูล thumbnail และ page_count
	// if err := s.repo.UpdateDocumentThumbnailAndPages(doc.DocumentID, thumbnailURL, pageCount); err != nil {
	// 	log.Printf("[db update] error: %v", err)
	// }

	// บันทึก URLs ของแต่ละหน้า
	// if len(pageURLs) > 0 {
	// 	if err := s.repo.SaveDocumentPages(doc.DocumentID, pageURLs); err != nil {
	// 		log.Printf("[db pages] error: %v", err)
	// 	}
	// }

	//response to clint
	resp := &models.UploadResponse{
		Message:    "อัปโหลดไฟล์สำเร็จ",
		File:       *doc,
		FileURL:    doc.DocumentURL,
		DocumentID: doc.DocumentID,
	}
	return resp, nil
}

// 2.ดึงข้อมูลไฟล์ทั้งหมดของผู้ใช้
func (s *fileService) GetFilesByUserID(userID int) ([]models.Document, error) {
	doc, err := s.repo.GetDocumentByUserID(userID)
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

// แปลง PDF เป็น JPG ด้วย pdfcpu
// func convertPDFToImages(pdfPath string, outputDir string) ([]string, error) {
// 	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
// 		return nil, fmt.Errorf("ไม่สามารถสร้างโฟลเดอร์: %v", err)
// 	}

// 	conf := pdfcpu.NewDefaultConfiguration()

// 	// ดึงภาพทั้งหมดจาก PDF
// 	if err := api.ExtractImagesFile(pdfPath, outputDir, nil, conf); err != nil {
// 		return nil, fmt.Errorf("แปลง PDF เป็นภาพไม่สำเร็จ: %v", err)
// 	}

// 	files, err := os.ReadDir(outputDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("อ่านโฟลเดอร์ภาพไม่สำเร็จ: %v", err)
// 	}

// 	var pageURLs []string
// 	for _, f := range files {
// 		if !f.IsDir() && (filepath.Ext(f.Name()) == ".jpg" || filepath.Ext(f.Name()) == ".jpeg" || filepath.Ext(f.Name()) == ".png") {
// 			url := fmt.Sprintf("http://localhost:8080/uploads/pages/%s/%s", filepath.Base(outputDir), f.Name())
// 			pageURLs = append(pageURLs, url)
// 		}
// 	}

// 	return pageURLs, nil
// }
