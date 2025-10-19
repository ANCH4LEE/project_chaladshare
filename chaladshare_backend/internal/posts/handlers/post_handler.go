package handlers

import (
	"net/http"
	"strconv"

	filesvc "chaladshare_backend/internal/files/service"
	"chaladshare_backend/internal/posts/models"
	postsvc "chaladshare_backend/internal/posts/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService postsvc.PostService
	fileService filesvc.FileService
}

func NewPostHandler(postService postsvc.PostService, fileService filesvc.FileService) *PostHandler {
	return &PostHandler{postService: postService, fileService: fileService}
}

// สร้างโพสต์ใหม่ (ต้องล็อกอิน)
func (h *PostHandler) CreatePost(c *gin.Context) {
	uid := c.GetInt("user_id") // เอาจาก JWT เท่านั้น
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Title       string   `json:"post_title" binding:"required"`
		Description string   `json:"post_description"`
		Visibility  string   `json:"post_visibility" binding:"required"` // ตอนนี้รองรับ "public" เท่านั้น
		DocumentID  *int     `json:"post_document_id"`
		SummaryID   *int     `json:"post_summary_id"`
		Tags        []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Visibility != "public" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported visibility"})
		return
	}

	post := &models.Post{
		AuthorUserID: uid, // ตั้งจาก JWT
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

// ดึงโพสต์ทั้งหมด (ต้องล็อกอิน)
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// รายละเอียดโพสต์ (ต้องล็อกอิน)
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

// แก้ไขโพสต์ (เฉพาะเจ้าของ)
func (h *PostHandler) UpdatePost(c *gin.Context) {
	uid := c.GetInt("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// เช็คสิทธิ์เจ้าของก่อน
	isOwner, err := h.postService.IsOwner(postID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		Title       string `json:"post_title"`
		Description string `json:"post_description"`
		Visibility  string `json:"post_visibility"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Visibility != "public" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported visibility"})
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

// ลบโพสต์ (เฉพาะเจ้าของ)
func (h *PostHandler) DeletePost(c *gin.Context) {
	uid := c.GetInt("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// เช็คสิทธิ์เจ้าของก่อน
	isOwner, err := h.postService.IsOwner(postID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.postService.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) GetPostPages(c *gin.Context) {
	viewerID := c.GetInt("user_id")
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}

	ok, reason, err := h.postService.ViewPost(viewerID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	if !ok {
		switch reason {
		case "not_found":
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case "friends_only":
			c.JSON(http.StatusForbidden, gin.H{"error": "post is friends only"})
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		}
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get post failed"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no ducumnet"})
		return
	}

	pages, err := h.fileService.GetDocumentPages(*post.DocumentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "load pages failed"})
		return
	}
	c.JSON(http.StatusOK, pages)
}
