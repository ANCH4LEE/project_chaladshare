package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"chaladshare_backend/internal/files/models"
	"chaladshare_backend/internal/files/service"
	"chaladshare_backend/internal/utils"
)

type FileHandler struct {
	fileservice service.FileService
}

func NewFileHandler(fileservice service.FileService) *FileHandler {
	return &FileHandler{fileservice: fileservice}
}

// upload file + แปลงภาพ + ลง db
func (h *FileHandler) UploadFile(c *gin.Context) {
	authUserID := c.GetInt("user_id")
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาแนบไฟล์ PDF"})
		return
	}

	savePath := fmt.Sprintf("./uploads/%s", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถบันทึกไฟล์ได้"})
		return
	}

	outputDir := fmt.Sprintf("./uploads/pages/%s", strings.TrimSuffix(file.Filename, ".pdf"))
	images, err := utils.ConvertPDFToImages(savePath, outputDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("แปลง PDF ไม่สำเร็จ: %v", err)})
		return
	}

	// เตรียมข้อมูลเพื่อส่งให้ service
	req := &models.UploadRequest{
		UserID:          authUserID,
		DocumentName:    file.Filename,
		DocumentURL:     savePath,
		StorageProvider: "local",
		Images:          images,
	}

	resp, err := h.fileservice.UploadFile(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ดึงไฟล์ทั้งหมดของผู้ใช้
func (h *FileHandler) GetFilesByUserID(c *gin.Context) {
	authUserID := c.GetInt("user_id")
	targetUserID, err := strconv.Atoi(c.Param("user_id"))

	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if authUserID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	files, err := h.fileservice.GetFilesByUserID(targetUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบไฟล์ของผู้ใช้นี้"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, files)
}

// ดึงภาพของแต่ละหน้า PDF
func (h *FileHandler) GetDocumentPages(c *gin.Context) {
	authUserID := c.GetInt("user_id")
	docID, err := strconv.Atoi(c.Param("document_id"))

	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_id"})
		return
	}

	ok, err := h.fileservice.IsOwner(docID, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	pages, err := h.fileservice.GetDocumentPages(docID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบหน้าเอกสารนี้"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, pages)
}

// ดึง summaries ตาม document_id
func (h *FileHandler) GetSummaryByDocumentID(c *gin.Context) {
	authUserID := c.GetInt("user_id")
	docID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_id"})
		return
	}
	ok, err := h.fileservice.IsOwner(docID, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	summary, err := h.fileservice.GetSummaryByDocumentID(docID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบสรุปของไฟล์นี้"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ลบไฟล์
func (h *FileHandler) DeleteFile(c *gin.Context) {
	// ดึง user_id จาก JWT (middleware ตั้งค่าไว้)
	authUserID := c.GetInt("user_id")
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// อ่าน document_id จาก path และตรวจความถูกต้อง
	docID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_id"})
		return
	}

	// อนุญาตเฉพาะเจ้าของไฟล์เท่านั้น
	ok, err := h.fileservice.IsOwner(docID, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// ลบไฟล์
	if err := h.fileservice.DeleteFile(docID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบไฟล์นี้"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบไฟล์สำเร็จ"})
}
