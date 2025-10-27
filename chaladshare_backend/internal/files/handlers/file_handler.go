package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"chaladshare_backend/internal/files/models"
	"chaladshare_backend/internal/files/service"
)

type FileHandler struct {
	fileservice service.FileService
}

func NewFileHandler(fileservice service.FileService) *FileHandler {
	return &FileHandler{fileservice: fileservice}
}

// POST /api/v1/files/upload
func (h *FileHandler) UploadFile(c *gin.Context) {
	uid := c.GetInt("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาแนบไฟล์ PDF"})
		return
	}

	id := uuid.New().String()
	filename := id + ".pdf"
	abs := filepath.Join("./uploads", filename)
	publicURL := "/uploads/" + filename

	if err := c.SaveUploadedFile(fh, abs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถบันทึกไฟล์ได้"})
		return
	}

	req := &models.UploadRequest{
		UserID:          uid,
		DocumentName:    fh.Filename,
		DocumentURL:     publicURL,
		StorageProvider: "local",
	}
	resp, err := h.fileservice.UploadFile(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"document_id": resp.DocumentID,
		"pdf_url":     publicURL,
	})
}

// GET /api/v1/users/:user_id/files
func (h *FileHandler) GetFilesByUserID(c *gin.Context) {
	authUID := c.GetInt("user_id")
	targetUID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || targetUID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if authUID != targetUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	files, err := h.fileservice.GetFilesByUserID(targetUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบไฟล์ของผู้ใช้นี้"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

// GET /api/v1/documents/:document_id/summary
func (h *FileHandler) GetSummaryByDocumentID(c *gin.Context) {
	authUID := c.GetInt("user_id")
	docID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_id"})
		return
	}

	ok, err := h.fileservice.IsOwner(docID, authUID)
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
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// DELETE /api/v1/documents/:document_id
func (h *FileHandler) DeleteFile(c *gin.Context) {
	authUID := c.GetInt("user_id")
	if authUID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	docID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document_id"})
		return
	}

	ok, err := h.fileservice.IsOwner(docID, authUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

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
