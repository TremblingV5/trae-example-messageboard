package router

import (
	"messageboard/controller"
	"messageboard/middleware"
	"messageboard/model"
	"messageboard/repository"
	"messageboard/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// Initialize repositories
	userRepo := repository.NewUserRepository(model.DB)
	postRepo := repository.NewPostRepository(model.DB)
	commentRepo := repository.NewCommentRepository(model.DB)
	voteRepo := repository.NewVoteRepository(model.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, voteRepo, postRepo)
	voteService := service.NewVoteService(voteRepo, commentRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(authService)
	postController := controller.NewPostController(postService)
	commentController := controller.NewCommentController(commentService, voteService)

	// Create router
	r := gin.New()

	// Middlewares
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Static files for uploads
	r.Static("/uploads", "./uploads")

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", middleware.AuthMiddleware(), userController.UpdateUser)
			users.POST("/:id/avatar", middleware.AuthMiddleware(), userController.UploadAvatar)
		}

		// Post routes
		posts := api.Group("/posts")
		{
			posts.GET("", postController.GetPostList)
			posts.GET("/search", postController.SearchPosts)
			posts.POST("", middleware.AuthMiddleware(), postController.CreatePost)
			posts.POST("/upload", middleware.AuthMiddleware(), postController.UploadPostImage)
			posts.GET("/:id", postController.GetPost)
		}

		// Comment routes (nested under posts)
		postsComments := api.Group("/posts")
		{
			postsComments.POST("/:id/comments", middleware.AuthMiddleware(), commentController.CreateComment)
			postsComments.GET("/:id/comments", middleware.OptionalAuth(), commentController.GetComments)
		}

		// Vote routes
		comments := api.Group("/comments")
		{
			comments.POST("/:id/vote", middleware.AuthMiddleware(), commentController.Vote)
			comments.GET("/:id/votes", middleware.OptionalAuth(), commentController.GetVotes)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}