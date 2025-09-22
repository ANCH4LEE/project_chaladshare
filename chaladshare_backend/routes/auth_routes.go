package routes

import (
	"chaladshare_backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterHandler)
	}
}
