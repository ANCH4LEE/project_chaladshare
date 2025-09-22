package main

import (
	"log"
	"os"
	"time"

	"chaladshare_backend/db"
	"chaladshare_backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด .env
	_ = godotenv.Load()

	// เชื่อมต่อ PostgreSQL
	if err := db.InitDB(); err != nil {
		log.Fatal("ไม่สามารถเชื่อมต่อฐานข้อมูล:", err)
	}
	defer db.Pool.Close()

	// เริ่ม Gin
	router := gin.Default()

	// CORS สำหรับ frontend React ที่ localhost:3000
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// จัดกลุ่ม route ทั้งหมดไว้ใน /api/v1
	v1 := router.Group("/api/v1")
	routes.RegisterV1Routes(v1)

	// กำหนดพอร์ต
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("ChaladShare Backend รันที่พอร์ต :", port)
	router.Run(":" + port)
}
