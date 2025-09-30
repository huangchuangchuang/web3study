package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Student represents the students table structure
type Student struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Age   int    `db:"age"`
	Grade string `db:"grade"`
}

func main() {
	// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	// 要求 ：
	// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "李四"，年龄为 20，年级为 "三年级"。
	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	// 编写SQL语句将 students 表中姓名为 "李四" 的学生年级更新为 "四年级"。
	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

	// Initialize database connection (example with MySQL)
	db, err := sqlx.Connect("mysql", "testuser:password@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Insert a new record: Add student "李四", age 20, grade "三年级"
	insertQuery := `INSERT INTO students (name, age, grade) VALUES (?, ?, ?)`
	_, err = db.Exec(insertQuery, "李四", 20, "三年级")
	if err != nil {
		log.Printf("Insert error: %v", err)
	} else {
		fmt.Println("Successfully inserted student 李四")
	}

	// 2. Query all students with age > 18
	selectQuery := `SELECT id, name, age, grade FROM students WHERE age > ?`
	var students []Student
	err = db.Select(&students, selectQuery, 18)
	if err != nil {
		log.Printf("Select error: %v", err)
	} else {
		fmt.Println("Students with age > 18:")
		for _, student := range students {
			fmt.Printf("ID: %d, Name: %s, Age: %d, Grade: %s\n",
				student.ID, student.Name, student.Age, student.Grade)
		}
	}

	// 3. Update grade to "四年级" for student named "李四"
	updateQuery := `UPDATE students SET grade = ? WHERE name = ?`
	_, err = db.Exec(updateQuery, "四年级", "李四")
	if err != nil {
		log.Printf("Update error: %v", err)
	} else {
		fmt.Println("Successfully updated 李四's grade to 四年级")
	}

	// 4. Delete students with age < 15
	deleteQuery := `DELETE FROM students WHERE age < ?`
	result, err := db.Exec(deleteQuery, 15)
	if err != nil {
		log.Printf("Delete error: %v", err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("Successfully deleted %d students with age < 15\n", rowsAffected)
	}
}
