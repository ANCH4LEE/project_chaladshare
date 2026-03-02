// internal/recommend/handlers/recommend_handler.go
package handlers

import (
	"net/http"
	"strconv"

	recommendrepo "chaladshare_backend/internal/recommend/repository"
	recommendservice "chaladshare_backend/internal/recommend/service"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	svc  recommendservice.RecommendService
	repo recommendrepo.RecommendReadRepo
}

func NewRecommendHandler(svc recommendservice.RecommendService, repo recommendrepo.RecommendReadRepo) *RecommendHandler {
	return &RecommendHandler{svc: svc, repo: repo}
}

func (h *RecommendHandler) GetRecommend(c *gin.Context) {
	uid := c.GetInt("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 20
	offset := 0

	if s := c.Query("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if s := c.Query("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v >= 0 {
			offset = v
		}
	}

	items, err := h.repo.ListUserRecommendations(uid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *RecommendHandler) RecomputeNow(c *gin.Context) {
	uid := c.GetInt("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if h.svc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "recommend service is nil"})
		return
	}
	if err := h.svc.RecomputeFromLikes(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
