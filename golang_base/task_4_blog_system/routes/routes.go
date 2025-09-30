package routes

import (
	"blog-system/controllers"
	"blog-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 控制器实例
	authController := &controllers.AuthController{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	// 公开路由
	public := router.Group("/api")
	{
		// 用户认证
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)

		// 文章相关（公开）
		public.GET("/posts", postController.GetPosts)
		public.GET("/posts/:id", postController.GetPost)

		// 评论相关（公开）
		public.GET("/posts/:id/comments", commentController.GetComments)
	}

	// 需要认证的路由
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 文章相关（需要认证）
		protected.POST("/posts", postController.CreatePost)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)

		// 评论相关（需要认证）
		protected.POST("/comments", commentController.CreateComment)
	}
}
