// internal/users/handlers/handler.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	postsvc "chaladshare_backend/internal/posts/service"
	"chaladshare_backend/internal/users/models"
	"chaladshare_backend/internal/users/service"
)

type UserHandler struct {
	svc     service.UserService
	postSvc postsvc.PostService
}

func NewUserHandler(s service.UserService, p postsvc.PostService) *UserHandler {
	return &UserHandler{svc: s, postSvc: p}
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

	prof, err := h.svc.GetOwnProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	withSet := map[string]bool{}
	if raw := strings.TrimSpace(c.Query("with")); raw != "" {
		for _, w := range strings.Split(raw, ",") {
			withSet[strings.ToLower(strings.TrimSpace(w))] = true
		}
	}
	want := func(k string) bool { return withSet["all"] || withSet[k] }

	// base response เดิม
	resp := gin.H{
		"user_id":         prof.UserID,
		"email":           prof.Email,
		"username":        prof.Username,
		"avatar_url":      prof.AvatarURL,
		"avatar_storage":  prof.AvatarStore,
		"bio":             prof.Bio,
		"user_status":     prof.Status,
		"user_created_at": prof.CreatedAt,
	}

	// ---- ตัวเลขที่ต้องการ (คำนวณเฉพาะสิ่งที่ขอ) ----
	if want("stats") {
		if cnt, err := h.postSvc.CountByUserID(uid); err == nil {
			resp["posts_count"] = cnt
		} else {
			resp["posts_count"] = 0
		}
	}
	if want("followers") {
		// ยังไม่มีระบบ follow → ให้ 0 ไปก่อน
		resp["followers_count"] = 0
		// ถ้ามีให้เรียก service จริงของคุณ เช่น: cnt, _ := h.svc.CountFollowers(uid)
	}
	if want("following") {
		resp["following_count"] = 0
		// ถ้ามี: cnt, _ := h.svc.CountFollowing(uid)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetViewedUserProfile(c *gin.Context) {
	if _, ok := getUID(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	prof, err := h.svc.GetViewedUserProfile(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	withSet := map[string]bool{}
	if raw := strings.TrimSpace(c.Query("with")); raw != "" {
		for _, w := range strings.Split(raw, ",") {
			withSet[strings.ToLower(strings.TrimSpace(w))] = true
		}
	}
	want := func(k string) bool { return withSet["all"] || withSet[k] }

	// base response (อย่าส่ง email/status ของคนอื่น)
	resp := gin.H{
		"user_id":        prof.UserID,
		"username":       prof.Username,
		"avatar_url":     prof.AvatarURL,
		"avatar_storage": prof.AvatarStore,
		"bio":            prof.Bio,
	}

	// ตัวเลขที่ขอ (อิง targetID)
	if want("stats") {
		if cnt, err := h.postSvc.CountByUserID(targetID); err == nil {
			resp["posts_count"] = cnt
		} else {
			resp["posts_count"] = 0
		}
	}
	if want("followers") {
		resp["followers_count"] = 0 // TODO: เปลี่ยนเป็น service จริงเมื่อพร้อม
	}
	if want("following") {
		resp["following_count"] = 0 // TODO: เปลี่ยนเป็น service จริงเมื่อพร้อม
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
