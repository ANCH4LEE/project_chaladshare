package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(rg *gin.RouterGroup) {
	RegisterAuthRoutes(rg)

	// ต่อไปสามารถเพิ่ม:
	// RegisterSummaryRoutes(rg)
	// RegisterUserRoutes(rg)
}
