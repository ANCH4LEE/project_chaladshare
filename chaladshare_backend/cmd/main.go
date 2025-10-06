package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"chaladshare_backend/internal/config"
	// "chaladshare_backend/internal/middleware"
	"chaladshare_backend/internal/connectdb"

	authHandler "chaladshare_backend/internal/auth/handlers"
	authRepo "chaladshare_backend/internal/auth/repository"
	authService "chaladshare_backend/internal/auth/service"

	FileHandler "chaladshare_backend/internal/files/handlers"
	FileRepo "chaladshare_backend/internal/files/repository"
	FileService "chaladshare_backend/internal/files/service"

	postHandler "chaladshare_backend/internal/posts/handlers"
	postRepo "chaladshare_backend/internal/posts/repository"
	postSrv "chaladshare_backend/internal/posts/service"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {

	//config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//connet DB
	db, err := connectdb.NewPostgresDatabase(cfg.GetConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	//auth repo service handler
	authRepository := authRepo.NewUserRepository(db.GetDB())
	authService := authService.NewAuthService(authRepository)
	authHandler := authHandler.NewAuthHandler(authService)

	//file repo service handler
	fileRepository := FileRepo.NewFileRepository(db.GetDB())
	fileService := FileService.NewFileService(fileRepository)
	fileHandler := FileHandler.NewFileHandler(fileService)

	//post repo service handler
	postRepository := postRepo.NewPostRepository(db.GetDB())
	likeRepository := postRepo.NewLikeRepository(db.GetDB())
	saveRepository := postRepo.NewSaveRepository(db.GetDB())

	postService := postSrv.NewPostService(postRepository)
	likeService := postSrv.NewLikeService(likeRepository)
	saveService := postSrv.NewSaveService(saveRepository)

	postHandler := postHandler.NewPostHandler(postService, likeService, saveService)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost: %v", err)
				if reconnErr := db.Reconnect(cfg.GetConnectionString()); reconnErr != nil {
					log.Printf("Failed to reconnect: %v", reconnErr)
				} else {
					log.Printf("Successfully reconnected to the database")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(TimeoutMiddleware(5 * time.Second))

	r.GET("/health", func(c *gin.Context) {
		if err := connectdb.CheckDBConnection(db.GetDB()); err != nil {
			c.JSON(503, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy", "database": "connected"})
	})

	v1 := r.Group("/api/v1")

	//login register
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	// //Protected routes (ต้องตรวจ JWT ก่อนเข้า)
	// protected := v1.Group("/user")
	// protected.Use(middleware.JWTAuthMiddleware())
	// {
	// 	protected.GET("/profile", authHandler.GetProfile)
	// }

	//file
	fileRoutes := v1.Group("/files")
	{
		fileRoutes.POST("/upload", fileHandler.UploadFile)
		fileRoutes.GET("/user/:user_id", fileHandler.GetFilesByUserID)
		fileRoutes.POST("/summary", fileHandler.SaveSummary)
		fileRoutes.GET("/summary/:document_id", fileHandler.GetSummaryByDocumentID)
	}

	//post liked saved
	postRoutes := v1.Group("/posts")
	{
		postRoutes.POST("/", postHandler.CreatePost)
		postRoutes.GET("/", postHandler.GetAllPosts)
		postRoutes.GET("/:id", postHandler.GetPostByID)
		postRoutes.PUT("/:id", postHandler.UpdatePost)
		postRoutes.DELETE("/:id", postHandler.DeletePost)

		postRoutes.POST("/:id/like", postHandler.LikePost)
		postRoutes.DELETE("/:id/like", postHandler.UnlikePost)

		postRoutes.POST("/:id/save", postHandler.SavePost)
		postRoutes.DELETE("/:id/save", postHandler.UnsavePost)
	}

	//Run Server
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
