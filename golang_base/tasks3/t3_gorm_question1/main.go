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
	// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	// 要求 ：
	// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
	// Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
	// 编写Go代码，使用Gorm创建这些模型对应的数据库表。

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

	// 显示表结构信息
	showTableInfo()

	// 插入示例数据
	err = insertSampleData(db)
	if err != nil {
		log.Printf("插入示例数据失败: %v", err)
	} else {
		fmt.Println("示例数据插入完成!")
	}

	// 查询示例
	queryExamples(db)
}

// 显示表结构信息
func showTableInfo() {
	fmt.Println("\n=== 数据库表结构 ===")
	fmt.Println("1. Users 表:")
	fmt.Println("   - ID: 主键，自增")
	fmt.Println("   - Name: 用户名，varchar(100)，非空")
	fmt.Println("   - Email: 邮箱，varchar(100)，唯一索引，非空")
	fmt.Println("   - CreatedAt: 创建时间")
	fmt.Println("   - UpdatedAt: 更新时间")
	fmt.Println("   - DeletedAt: 软删除时间")
	fmt.Println("   - 关系: 与 Posts 一对多")

	fmt.Println("\n2. Posts 表:")
	fmt.Println("   - ID: 主键，自增")
	fmt.Println("   - Title: 标题，varchar(200)，非空")
	fmt.Println("   - Content: 内容，text")
	fmt.Println("   - UserID: 外键，关联 users 表")
	fmt.Println("   - CreatedAt: 创建时间")
	fmt.Println("   - UpdatedAt: 更新时间")
	fmt.Println("   - DeletedAt: 软删除时间")
	fmt.Println("   - 关系: 与 User 多对一，与 Comments 一对多")

	fmt.Println("\n3. Comments 表:")
	fmt.Println("   - ID: 主键，自增")
	fmt.Println("   - Content: 评论内容，text，非空")
	fmt.Println("   - UserID: 外键，关联 users 表")
	fmt.Println("   - PostID: 外键，关联 posts 表")
	fmt.Println("   - CreatedAt: 创建时间")
	fmt.Println("   - UpdatedAt: 更新时间")
	fmt.Println("   - DeletedAt: 软删除时间")
	fmt.Println("   - 关系: 与 User 多对一，与 Post 多对一")
}

// 插入示例数据
func insertSampleData(db *gorm.DB) error {
	fmt.Println("\n=== 插入示例数据 ===")

	// 创建用户
	user1 := User{Name: "张三", Email: "zhangsan@example.com"}
	user2 := User{Name: "李四", Email: "lisi@example.com"}

	if result := db.Create(&user1); result.Error != nil {
		return fmt.Errorf("创建用户1失败: %v", result.Error)
	}
	if result := db.Create(&user2); result.Error != nil {
		return fmt.Errorf("创建用户2失败: %v", result.Error)
	}

	fmt.Printf("创建用户: %s (ID: %d)\n", user1.Name, user1.ID)
	fmt.Printf("创建用户: %s (ID: %d)\n", user2.Name, user2.ID)

	// 创建文章
	posts := []Post{
		{Title: "Go语言入门", Content: "Go语言是一门现代化的编程语言...", UserID: user1.ID},
		{Title: "GORM使用指南", Content: "GORM是Go语言的ORM库...", UserID: user1.ID},
		{Title: "数据库设计", Content: "良好的数据库设计是系统成功的关键...", UserID: user2.ID},
	}

	for i := range posts {
		if result := db.Create(&posts[i]); result.Error != nil {
			return fmt.Errorf("创建文章失败: %v", result.Error)
		}
		fmt.Printf("创建文章: %s (ID: %d)\n", posts[i].Title, posts[i].ID)
	}

	// 创建评论
	comments := []Comment{
		{Content: "很好的文章，学到了很多！", UserID: user2.ID, PostID: posts[0].ID},
		{Content: "感谢分享，期待更多内容", UserID: user2.ID, PostID: posts[1].ID},
		{Content: "这篇文章对我帮助很大", UserID: user1.ID, PostID: posts[2].ID},
	}

	for i := range comments {
		if result := db.Create(&comments[i]); result.Error != nil {
			return fmt.Errorf("创建评论失败: %v", result.Error)
		}
		fmt.Printf("创建评论: %s (ID: %d)\n", comments[i].Content[:10]+"...", comments[i].ID)
	}

	return nil
}

// 查询示例
func queryExamples(db *gorm.DB) {
	fmt.Println("\n=== 查询示例 ===")

	// 查询所有用户及其文章
	var users []User
	// 只执行 2 次查询，无论有多少用户 总查询次数: 2 次
	db.Preload("Posts").Find(&users)
	fmt.Println("用户及其文章:")
	for _, user := range users {
		fmt.Printf("用户: %s, 文章数: %d\n", user.Name, len(user.Posts))
	}

	// 查询特定文章及其评论
	var post Post
	db.Preload("Comments").Preload("User").First(&post)
	fmt.Printf("\n文章: %s (作者: %s)\n", post.Title, post.User.Name)
	fmt.Printf("评论数: %d\n", len(post.Comments))
}
