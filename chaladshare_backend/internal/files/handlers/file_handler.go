package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
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
	var form models.UploadFrom
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุ user_id"})
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
		UserID:          form.UserID,
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
	userIDStr := c.Param("user_id")
	var userID int
	fmt.Sscanf(userIDStr, "%d", &userID)

	files, err := h.fileservice.GetFilesByUserID(userID)
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
	docIDStr := c.Param("document_id")
	var docID int
	fmt.Sscanf(docIDStr, "%d", &docID)

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
	docIDStr := c.Param("document_id")
	var docID int
	fmt.Sscanf(docIDStr, "%d", &docID)

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
	idStr := c.Param("document_id")
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	if err := h.fileservice.DeleteFile(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบไฟล์นี้"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบไฟล์สำเร็จ"})
}
