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
	PostCount int            `gorm:"default:0" json:"post_count"` // 新增：文章数量统计字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 一对多关系：一个用户可以发布多篇文章
	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(200);not null" json:"title"`
	Content       string         `gorm:"type:text" json:"content"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`                       // 新增：评论数量字段
	CommentStatus string         `gorm:"type:varchar(50);default:'无评论'" json:"comment_status"` // 新增：评论状态字段
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

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

// ==================== 钩子函数实现 ====================

// 1. Post 模型的钩子函数：在文章创建后更新用户的文章数量统计
func (p *Post) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("钩子函数: 文章 '%s' 创建完成\n", p.Title)

	// 更新用户的文章数量统计
	result := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("钩子函数: 用户ID %d 的文章数量已更新\n", p.UserID)
	return nil
}

// 2. Post 模型的钩子函数：在文章删除后更新用户的文章数量统计
func (p *Post) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("钩子函数: 文章ID %d 删除完成\n", p.ID)

	// 更新用户的文章数量统计
	result := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("钩子函数: 用户ID %d 的文章数量已减少\n", p.UserID)
	return nil
}

// 3. Comment 模型的钩子函数：在评论创建后更新文章的评论数量和状态
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("钩子函数: 评论ID %d 创建完成\n", c.ID)

	// 更新文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}

	// 更新文章评论状态为"有评论"
	result = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论")
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("钩子函数: 文章ID %d 的评论数量和状态已更新\n", c.PostID)
	return nil
}

// 4. Comment 模型的钩子函数：在评论删除前检查文章的评论数量
func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("钩子函数: 准备删除评论ID %d\n", c.ID)

	// 获取文章当前的评论数量
	var post Post
	result := tx.Select("id, comment_count").Where("id = ?", c.PostID).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil // 文章不存在，直接返回
		}
		return result.Error
	}

	// 如果这是文章的最后一条评论
	if post.CommentCount <= 1 {
		// 更新文章评论状态为"无评论"
		result := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论")
		if result.Error != nil {
			return result.Error
		}
		fmt.Printf("钩子函数: 文章ID %d 的评论状态已更新为'无评论'\n", c.PostID)
	}

	return nil
}

// 5. Comment 模型的钩子函数：在评论删除后更新文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("钩子函数: 评论ID %d 删除完成\n", c.ID)

	// 更新文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("钩子函数: 文章ID %d 的评论数量已减少\n", c.PostID)
	return nil
}

func main() {
	// 	为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
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

	// 演示钩子函数功能
	demonstrateHooks(db)
}

// 演示钩子函数功能
func demonstrateHooks(db *gorm.DB) {
	fmt.Println("\n=== 演示钩子函数功能 ===")

	// 1. 创建用户
	user := User{Name: "张三疯", Email: "zhangsanfeng@example.com"}
	if result := db.Create(&user); result.Error != nil {
		log.Printf("创建用户失败: %v", result.Error)
		return
	}
	fmt.Printf("创建用户: %s (ID: %d)\n", user.Name, user.ID)

	// 2. 创建文章（触发 AfterCreate 钩子）
	post1 := Post{Title: "Go语言入门", Content: "Go语言是一门现代化的编程语言...", UserID: user.ID}
	post2 := Post{Title: "GORM使用指南", Content: "GORM是Go语言的ORM库...", UserID: user.ID}

	if result := db.Create(&post1); result.Error != nil {
		log.Printf("创建文章1失败: %v", result.Error)
		return
	}
	fmt.Printf("\n创建文章1: %s (ID: %d)\n", post1.Title, post1.ID)

	if result := db.Create(&post2); result.Error != nil {
		log.Printf("创建文章2失败: %v", result.Error)
		return
	}
	fmt.Printf("创建文章2: %s (ID: %d)\n", post2.Title, post2.ID)

	// 检查用户的文章数量
	var updatedUser User
	db.First(&updatedUser, user.ID)
	fmt.Printf("用户 %s 的文章数量: %d\n", updatedUser.Name, updatedUser.PostCount)

	// 3. 创建评论（触发 AfterCreate 钩子）
	comment1 := Comment{Content: "很好的文章！", UserID: user.ID, PostID: post1.ID}
	comment2 := Comment{Content: "学到了很多知识", UserID: user.ID, PostID: post1.ID}

	if result := db.Create(&comment1); result.Error != nil {
		log.Printf("创建评论1失败: %v", result.Error)
		return
	}
	fmt.Printf("\n创建评论1: %s (ID: %d)\n", comment1.Content, comment1.ID)

	if result := db.Create(&comment2); result.Error != nil {
		log.Printf("创建评论2失败: %v", result.Error)
		return
	}
	fmt.Printf("创建评论2: %s (ID: %d)\n", comment2.Content, comment2.ID)

	// 检查文章的评论数量和状态
	var updatedPost Post
	db.First(&updatedPost, post1.ID)
	fmt.Printf("文章 '%s' 的评论数量: %d, 评论状态: %s\n",
		updatedPost.Title, updatedPost.CommentCount, updatedPost.CommentStatus)

	// 4. 删除评论（触发 BeforeDelete 和 AfterDelete 钩子）
	fmt.Println("\n=== 删除评论测试 ===")

	// 删除第一条评论
	if result := db.Delete(&comment1); result.Error != nil {
		log.Printf("删除评论1失败: %v", result.Error)
		return
	}
	fmt.Printf("删除评论1 (ID: %d)\n", comment1.ID)

	// 检查文章状态
	db.First(&updatedPost, post1.ID)
	fmt.Printf("删除评论后，文章 '%s' 的评论数量: %d, 评论状态: %s\n",
		updatedPost.Title, updatedPost.CommentCount, updatedPost.CommentStatus)

	// 删除最后一条评论
	if result := db.Delete(&comment2); result.Error != nil {
		log.Printf("删除评论2失败: %v", result.Error)
		return
	}
	fmt.Printf("删除评论2 (ID: %d)\n", comment2.ID)

	// 检查文章状态（应该变为"无评论"）
	db.First(&updatedPost, post1.ID)
	fmt.Printf("删除所有评论后，文章 '%s' 的评论数量: %d, 评论状态: %s\n",
		updatedPost.Title, updatedPost.CommentCount, updatedPost.CommentStatus)

	// 5. 删除文章（触发 AfterDelete 钩子）
	fmt.Println("\n=== 删除文章测试 ===")

	if result := db.Delete(&post1); result.Error != nil {
		log.Printf("删除文章失败: %v", result.Error)
		return
	}
	fmt.Printf("删除文章: %s (ID: %d)\n", post1.Title, post1.ID)

	// 检查用户的文章数量
	db.First(&updatedUser, user.ID)
	fmt.Printf("删除文章后，用户 %s 的文章数量: %d\n", updatedUser.Name, updatedUser.PostCount)
}
