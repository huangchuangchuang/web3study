package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Book 结构体，对应 books 表的字段
type Book struct {
	ID     int     `db:"id"`     // 主键ID
	Title  string  `db:"title"`  // 书名
	Author string  `db:"author"` // 作者
	Price  float64 `db:"price"`  // 价格
}

func main() {
	// 初始化数据库连接
	db, err := sqlx.Connect("mysql", "testuser:password@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 确保books表存在
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL
	)`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建表失败: %v", err)
		return
	}

	// 插入一些测试数据
	booksData := []Book{
		{Title: "Go语言编程", Author: "张三", Price: 68.50},
		{Title: "Python核心编程", Author: "李四", Price: 98.00},
		{Title: "Java从入门到放弃", Author: "王五", Price: 45.00},
		{Title: "C++ Primer", Author: "Stanley", Price: 128.00},
		{Title: "JavaScript权威指南", Author: "David", Price: 158.00},
		{Title: "算法导论", Author: "Thomas", Price: 89.50},
	}

	for _, book := range booksData {
		insertSQL := `INSERT INTO books (title, author, price) VALUES (?, ?, ?)`
		_, err = db.Exec(insertSQL, book.Title, book.Author, book.Price)
		if err != nil {
			log.Printf("插入数据失败: %v", err)
		}
	}

	// 执行复杂查询：查询价格大于50元的书籍
	// 使用sqlx.Select方法确保类型安全映射
	querySQL := `SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC`
	var books []Book

	err = db.Select(&books, querySQL, 50.0)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}

	// 输出查询结果
	fmt.Println("价格大于50元的书籍列表：")
	fmt.Println("========================")
	for _, book := range books {
		fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %.2f元\n",
			book.ID, book.Title, book.Author, book.Price)
	}

	fmt.Printf("\n总共找到 %d 本符合条件的书籍\n", len(books))
}
