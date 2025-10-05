package handlers

import (
	"net/http"
	"strconv"

	"chaladshare_backend/internal/files/models"
	"chaladshare_backend/internal/files/service"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileservice service.FileService
}

func NewFileHandler(fileservice service.FileService) *FileHandler {
	return &FileHandler{fileservice: fileservice}
}

// 1.อัปโหลดไฟล์ post /api/files/upload
func (h *FileHandler) UploadFile(c *gin.Context) {
	var req models.UploadRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ไม่ถูกต้อง"})
		return
	}

	//เรียก service เพื่อบันทึกไฟล์
	resp, err := h.fileservice.UploadFile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp) //response to client
}

// 2.ดึงข้อมูลไฟล์ทั้งหมดของผู้ใช้ get /api/files/user/:user_id
func (h *FileHandler) GetFilesByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id ไม่ถูกต้อง"})
		return
	}

	files, err := h.fileservice.GetFilesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

// 3.บันทึกสรุปที่ได้มาจาก AI post /api/files/summary
func (h *FileHandler) SaveSummary(c *gin.Context) {
	var req models.Summary

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลสรุปไม่ถูกต้อง"})
		return
	}

	if err := h.fileservice.SaveSummary(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "บันทึกสรุปสำเร็จ"})
}

// 4.ดึงสรุปจาก document_id get /api/files/summary/:document_id
func (h *FileHandler) GetSummaryByDocumentID(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_id ไม่ถูกต้อง"})
		return
	}

	summary, err := h.fileservice.GetSummaryByDocumentID(docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}
