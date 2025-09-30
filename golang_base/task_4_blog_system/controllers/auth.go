package controllers

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	// 1. 创建一个空的 RegisterInput 变量来存储用户提交的数据
	var input RegisterInput

	// 2. 从 HTTP 请求中提取 JSON 数据并验证
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果数据格式不正确，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 检查用户名或邮箱是否已经存在
	var existingUser models.User
	result := config.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser)
	if result.Error == nil {
		// 错误为 nil，说明找到了记录
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// 4. 对密码进行加密
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		// 如果加密失败，返回 500 服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 5. 创建新用户对象
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword, // 存储加密后的密码
	}

	// 6. 将用户保存到数据库
	if err := config.DB.Create(&user).Error; err != nil {
		// 如果保存失败，返回 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 7. 注册成功，返回 201 创建成功状态
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	result := config.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
