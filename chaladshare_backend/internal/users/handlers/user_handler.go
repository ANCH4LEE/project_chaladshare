// internal/users/handlers/handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"chaladshare_backend/internal/users/models"
	"chaladshare_backend/internal/users/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func getUID(c *gin.Context) (int, bool) {
	for _, k := range []string{"user_id", "uid"} {
		if v, ok := c.Get(k); ok {
			switch t := v.(type) {
			case int:
				return t, true
			case int64:
				return int(t), true
			case float64:
				return int(t), true
			}
		}
	}
	return 0, false
}

func (h *UserHandler) GetOwnProfile(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.svc.GetOwnProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetViewedUserProfile(c *gin.Context) {
	if _, ok := getUID(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	resp, err := h.svc.GetViewedUserProfile(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) UpdateOwnProfile(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.UpdateOwnProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.svc.UpdateOwnProfile(c.Request.Context(), uid, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
