package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Employee 员工结构体，对应 employees 表
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func main() {
	// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	// 要求 ：
	// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

	// 初始化数据库连接
	db, err := sqlx.Connect("mysql", "testuser:password@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建 employees 表
	createTable(db)

	// 插入测试数据
	insertTestData(db)

	// 1. 查询部门为"技术部"的所有员工
	fmt.Println("1. 查询技术部所有员工:")
	fmt.Println("=====================")
	techEmployees, err := queryEmployeesByDepartment(db, "技术部")
	if err != nil {
		log.Printf("查询技术部员工失败: %v", err)
	} else {
		for _, emp := range techEmployees {
			fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
				emp.ID, emp.Name, emp.Department, emp.Salary)
		}
		fmt.Printf("技术部共有 %d 名员工\n\n", len(techEmployees))
	}

	// 2. 查询工资最高的员工
	fmt.Println("2. 查询工资最高的员工:")
	fmt.Println("=====================")
	topEmployee, err := queryHighestPaidEmployee(db)
	if err != nil {
		log.Printf("查询最高工资员工失败: %v", err)
	} else {
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
			topEmployee.ID, topEmployee.Name, topEmployee.Department, topEmployee.Salary)
	}
}

// 创建 employees 表
func createTable(db *sqlx.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(50) NOT NULL,
		salary DECIMAL(10,2) NOT NULL
	)`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建表失败: %v", err)
	}
}

// 插入测试数据
func insertTestData(db *sqlx.DB) {
	// 清空表数据
	db.Exec("DELETE FROM employees")

	// 插入测试员工数据
	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 12000.00},
		{Name: "李四", Department: "技术部", Salary: 15000.00},
		{Name: "王五", Department: "技术部", Salary: 18000.00},
		{Name: "赵六", Department: "销售部", Salary: 8000.00},
		{Name: "钱七", Department: "销售部", Salary: 9500.00},
		{Name: "孙八", Department: "人事部", Salary: 10000.00},
		{Name: "周九", Department: "财务部", Salary: 11000.00},
		{Name: "吴十", Department: "技术部", Salary: 22000.00},
	}

	for _, emp := range employees {
		_, err := db.Exec("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)",
			emp.Name, emp.Department, emp.Salary)
		if err != nil {
			log.Printf("插入员工数据失败: %v", err)
		}
	}

	fmt.Println("初始化员工测试数据完成")
	fmt.Println("====================")
}

// 1. 查询指定部门的所有员工
func queryEmployeesByDepartment(db *sqlx.DB, department string) ([]Employee, error) {
	query := `SELECT id, name, department, salary FROM employees WHERE department = ?`
	var employees []Employee

	// 使用 sqlx.Select 将查询结果映射到 Employee 结构体切片
	err := db.Select(&employees, query, department)
	if err != nil {
		return nil, fmt.Errorf("查询员工失败: %v", err)
	}

	return employees, nil
}

// 2. 查询工资最高的员工
func queryHighestPaidEmployee(db *sqlx.DB) (*Employee, error) {
	query := `SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1`
	var employee Employee

	// 使用 sqlx.Get 将单行查询结果映射到 Employee 结构体
	err := db.Get(&employee, query)
	if err != nil {
		return nil, fmt.Errorf("查询最高工资员工失败: %v", err)
	}

	return &employee, nil
}
