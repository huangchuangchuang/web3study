package controllers

import (
	"blog-system/config"
	"blog-system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查文章是否存在
	var post models.Post
	result := config.DB.First(&post, input.PostID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := models.Comment{
		Content: input.Content,
		UserID:  userID,
		PostID:  input.PostID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 预加载用户和文章信息
	config.DB.Preload("User").Preload("Post").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// GetComments 获取某篇文章的所有评论
func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 检查文章是否存在
	var post models.Post
	result := config.DB.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var comments []models.Comment
	result = config.DB.Where("post_id = ?", postID).Preload("User").Find(&comments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
