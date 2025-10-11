package handlers

import (
	"net/http"
	"strconv"

	"chaladshare_backend/internal/posts/models"
	"chaladshare_backend/internal/posts/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

// ทำแค่ post service ก่อน
func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// สร้างโพสต์ใหม่
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		AuthorUserID int      `json:"author_user_id" binding:"required"`
		Title        string   `json:"post_title" binding:"required"`
		Description  string   `json:"post_description"`
		Visibility   string   `json:"post_visibility" binding:"required"`
		DocumentID   *int     `json:"post_document_id"`
		SummaryID    *int     `json:"post_summary_id"`
		Tags         []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	post := &models.Post{
		AuthorUserID: req.AuthorUserID,
		Title:        req.Title,
		Description:  req.Description,
		Visibility:   req.Visibility,
		DocumentID:   req.DocumentID,
		SummaryID:    req.SummaryID,
	}

	postID, err := h.postService.CreatePost(post, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", "/posts/"+strconv.Itoa(postID))
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"post_id": postID}})
}

// ดึงโพสต์ทั้งหมด
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// detail แต่ละโพสต์
func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// แก้ไขโพสต์
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Title       string `json:"post_title" binding:"required"`
		Description string `json:"post_description"`
		Visibility  string `json:"post_visibility" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	post := &models.Post{
		PostID:      postID,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
	}

	if err := h.postService.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
}

// ลบโพสต์
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))

	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.postService.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
