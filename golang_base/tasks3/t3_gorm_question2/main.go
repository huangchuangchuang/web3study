package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 一对多关系：一个用户可以发布多篇文章
	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	PostID    uint           `gorm:"not null" json:"post_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user"`
	Post Post `gorm:"foreignKey:PostID" json:"post"`
}

func main() {
	// 要求 ：
	// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	// 数据库连接配置
	dsn := "testuser:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	fmt.Println("数据库连接成功!")

	// 自动迁移模型到数据库表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	fmt.Println("数据库表创建成功!")

	// 查询示例
	queryUserPostsWithComments(db) // 查询用户发布的所有文章及评论
	queryPostWithMostComments(db)  // 查询评论数量最多的文章
	queryPostWithMostCommentsByRawSQL(db)
	queryTopNPostsWithMostComments(db, 3)
}

// 1. 查询某个用户发布的所有文章及其对应的评论信息
func queryUserPostsWithComments(db *gorm.DB) {
	fmt.Println("\n=== 查询用户发布的所有文章及评论 ===")

	// 假设查询用户名为"张三"的用户
	var user User
	result := db.Where("name = ?", "张三").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("未找到用户: 张三")
			return
		}
		log.Printf("查询用户失败: %v", result.Error)
		return
	}

	// 查询该用户的所有文章，并预加载每篇文章的评论
	var posts []Post
	result = db.Where("user_id = ?", user.ID).
		Preload("Comments").
		Preload("Comments.User"). // 同时加载评论的作者信息
		Find(&posts)

	if result.Error != nil {
		log.Printf("查询用户文章失败: %v", result.Error)
		return
	}

	fmt.Printf("用户: %s (ID: %d)\n", user.Name, user.ID)
	fmt.Printf("共发布 %d 篇文章:\n", len(posts))

	for _, post := range posts {
		fmt.Printf("\n  文章: %s\n", post.Title)
		fmt.Printf("  内容: %.50s...\n", post.Content)
		fmt.Printf("  评论数: %d\n", len(post.Comments))

		if len(post.Comments) > 0 {
			fmt.Println("  评论详情:")
			for _, comment := range post.Comments {
				fmt.Printf("    - %s (评论者: %s)\n",
					comment.Content, comment.User.Name)
			}
		}
	}
}

// 2. 查询评论数量最多的文章信息
func queryPostWithMostComments(db *gorm.DB) {
	fmt.Println("\n=== 查询评论数量最多的文章 ===")

	// 方法1: 使用子查询（推荐）
	var post Post
	result := db.
		Preload("User").          // 加载文章作者
		Preload("Comments").      // 加载文章评论
		Preload("Comments.User"). // 加载评论作者
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("COUNT(comments.id) DESC").
		First(&post)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("未找到任何文章")
			return
		}
		log.Printf("查询文章失败: %v", result.Error)
		return
	}

	fmt.Printf("评论数最多的文章:\n")
	fmt.Printf("标题: %s\n", post.Title)
	fmt.Printf("作者: %s\n", post.User.Name)
	fmt.Printf("评论数: %d\n", len(post.Comments))

	if len(post.Comments) > 0 {
		fmt.Println("评论列表:")
		for _, comment := range post.Comments {
			fmt.Printf("  - %s (评论者: %s)\n",
				comment.Content, comment.User.Name)
		}
	}
}

// 方法2: 使用原生SQL查询评论数最多的文章
func queryPostWithMostCommentsByRawSQL(db *gorm.DB) {
	fmt.Println("\n=== 使用原生SQL查询评论数最多的文章 ===")

	type PostCommentCount struct {
		PostID       uint
		Title        string
		UserName     string
		CommentCount int
	}

	var result PostCommentCount
	db.Raw(`
		SELECT p.id as post_id, p.title, u.name as user_name, COUNT(c.id) as comment_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		GROUP BY p.id, p.title, u.name
		ORDER BY comment_count DESC
		LIMIT 1
	`).Scan(&result)

	fmt.Printf("评论数最多的文章:\n")
	fmt.Printf("标题: %s\n", result.Title)
	fmt.Printf("作者: %s\n", result.UserName)
	fmt.Printf("评论数: %d\n", result.CommentCount)
}

// 方法3: 查询前N篇评论数最多的文章
func queryTopNPostsWithMostComments(db *gorm.DB, limit int) {
	fmt.Printf("\n=== 查询前%d篇评论数最多的文章 ===", limit)

	type PostWithCount struct {
		Post
		CommentCount int `gorm:"column:comment_count"`
	}

	var posts []PostWithCount
	db.Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(limit).
		Find(&posts)

	for _, post := range posts {
		fmt.Printf("文章: %s (作者: %s) - 评论数: %d\n",
			post.Title, post.User.Name, post.CommentCount)
	}
}
