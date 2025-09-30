package main

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	config.ConnectDatabase()
	// 自动迁移数据库表
	err := config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建Gin路由器
	router := gin.Default()

	// 设置路由
	routes.SetupRoutes(router)

	// 启动服务器
	log.Println("Server starting on port 8080...")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
