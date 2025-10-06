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
	likeService service.LikeService
	saveService service.SaveService
}

func NewPostHandler(
	postService service.PostService,
	likeService service.LikeService,
	saveService service.SaveService,
) *PostHandler {
	return &PostHandler{
		postService: postService,
		likeService: likeService,
		saveService: saveService,
	}
}

// 1.สร้างโพสต์ใหม่
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		AuthorUserID int      `json:"author_user_id"`
		Title        string   `json:"post_title"`
		Description  string   `json:"post_description"`
		Visibility   string   `json:"post_visibility"`
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

	c.JSON(http.StatusCreated, gin.H{"message": "post created successfully", "post_id": postID})
}

// 2. ดึงโพสต์ทั้งหมด
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// 3. detail แต่ละโพสต์
func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := h.postService.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// 4.แก้ไขโพสต์
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title       string `json:"post_title"`
		Description string `json:"post_description"`
		Visibility  string `json:"post_visibility"`
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

// 5.ลบโพสต์
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))

	if err := h.postService.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

// 6.จัดการการกดถูกใจ
func (h *PostHandler) LikePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id") // ดึงจาก JWT ภายหลัง (ตอนนี้ mock ไว้ได้)

	if err := h.likeService.LikePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

func (h *PostHandler) UnlikePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := h.likeService.UnlikePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

// 7.จัดการการบันทึกโพสต์
func (h *PostHandler) SavePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := h.saveService.SavePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

func (h *PostHandler) UnsavePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := h.saveService.UnsavePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unsaved"})
}
