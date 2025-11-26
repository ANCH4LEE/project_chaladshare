package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"chaladshare_backend/internal/config"
	"chaladshare_backend/internal/connectdb"
	"chaladshare_backend/internal/middleware"

	AuthHandler "chaladshare_backend/internal/auth/handlers"
	AuthRepo "chaladshare_backend/internal/auth/repository"
	AuthService "chaladshare_backend/internal/auth/service"

	FileHandler "chaladshare_backend/internal/files/handlers"
	FileRepo "chaladshare_backend/internal/files/repository"
	FileService "chaladshare_backend/internal/files/service"

	PostHandler "chaladshare_backend/internal/posts/handlers"
	PostRepo "chaladshare_backend/internal/posts/repository"
	PostService "chaladshare_backend/internal/posts/service"

	UserHandler "chaladshare_backend/internal/users/handlers"
	UserRepo "chaladshare_backend/internal/users/repository"
	UserService "chaladshare_backend/internal/users/service"

	FriendsHandler "chaladshare_backend/internal/friends/handlers"
	FriendsRepo "chaladshare_backend/internal/friends/repository"
	FriendsService "chaladshare_backend/internal/friends/service"
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
	authRepository := AuthRepo.NewAuthRepository(db.GetDB())
	authService := AuthService.NewAuthService(authRepository, []byte(cfg.JWTSecret), cfg.TokenTTLMinutes)
	authHandler := AuthHandler.NewAuthHandler(authService, cfg.CookieName, false)

	//file repo service handler
	fileRepository := FileRepo.NewFileRepository(db.GetDB())
	fileService := FileService.NewFileService(fileRepository)
	fileHandler := FileHandler.NewFileHandler(fileService)

	//post repo service handler
	postRepository := PostRepo.NewPostRepository(db.GetDB())
	postService := PostService.NewPostService(postRepository)
	postHandler := PostHandler.NewPostHandler(postService)

	friendsRepo := FriendsRepo.NewFriendRepository(db.GetDB())
	friendsService := FriendsService.NewFriendService(friendsRepo)
	friendsHandler := FriendsHandler.NewFriendHandler(friendsService)

	// user repo service handler
	userRepository := UserRepo.NewUserRepository(db.GetDB())
	userService := UserService.NewUserService(userRepository)
	userHandler := UserHandler.NewUserHandler(userService, postService, friendsService)

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
		AllowOrigins:     []string{cfg.AllowOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(TimeoutMiddleware(60 * time.Second))
	r.MaxMultipartMemory = 100 << 20
	r.Static("/uploads", "./uploads")

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatalf("cannot create uploads dir: %v", err)
	}

	r.GET("/health", func(c *gin.Context) {
		if err := connectdb.CheckDBConnection(db.GetDB()); err != nil {
			c.JSON(503, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy", "database": "connected"})
	})

	v1 := r.Group("/api/v1")

	//login register
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)

		authRoutes.GET("/users", authHandler.GetAllUsers)
		authRoutes.GET("/users/:id", authHandler.GetUserByID)
	}

	// --- Protected (ต้องมี JWT) ---
	protected := v1.Group("/")
	protected.Use(middleware.JWT([]byte(cfg.JWTSecret), cfg.CookieName))
	{
		// Post
		posts := protected.Group("/posts")
		{
			posts.GET("", postHandler.GetAllPosts)
			posts.GET("/:id", postHandler.GetPostByID)

			posts.POST("", postHandler.CreatePost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
		}

		// Files
		files := protected.Group("/files")
		{
			files.POST("/doc", fileHandler.UploadFile)
			files.GET("/user/:user_id", fileHandler.GetFilesByUserID)
			files.GET("/:document_id/summary", fileHandler.GetSummaryByDocumentID)
			files.DELETE("/:document_id", fileHandler.DeleteFile)

			files.POST("/cover", fileHandler.UploadCover)
			files.POST("/avatar", fileHandler.UploadAvatar)
		}

		//Profile
		profile := protected.Group("/profile")
		{
			profile.GET("", userHandler.GetOwnProfile)
			profile.PUT("", userHandler.UpdateOwnProfile)
			profile.GET("/:id", userHandler.GetViewedUserProfile)
		}

		social := protected.Group("/social")
		{
			// follow / unfollow
			social.POST("/follow", friendsHandler.FollowUser)         // body: {"followed_user_id": 123}
			social.DELETE("/follow/:id", friendsHandler.UnfollowUser) // :id = target user id

			// lists
			social.GET("/friends/:userID", friendsHandler.ListFriends)
			social.GET("/followers/:userID", friendsHandler.ListFollowers) // owner-only guard ใน service
			social.GET("/following/:userID", friendsHandler.ListFollowing) // owner-only guard ใน service

			// counters
			social.GET("/stats/:userID", friendsHandler.GetStats)

			// friend requests
			social.POST("/requests", friendsHandler.SendFriendRequest)
			social.GET("/requests/incoming", friendsHandler.ListIncomingRequests)
			social.GET("/requests/outgoing", friendsHandler.ListOutgoingRequests)
			social.POST("/requests/:id/accept", friendsHandler.AcceptFriendRequest)
			social.POST("/requests/:id/decline", friendsHandler.DeclineFriendRequest)
			social.DELETE("/requests/:id", friendsHandler.CancelFriendRequest)

			// unfriend
			social.DELETE("/friends/:userID", friendsHandler.Unfriend)
		}

	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
