package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Account 账户结构体
type Account struct {
	ID      int     `db:"id"`
	Balance float64 `db:"balance"`
}

// Transaction 交易记录结构体
type Transaction struct {
	ID            int     `db:"id"`
	FromAccountID int     `db:"from_account_id"`
	ToAccountID   int     `db:"to_account_id"`
	Amount        float64 `db:"amount"`
}

func main() {
	// 题目2：事务语句 假设有两个表：
	// accounts 表（包含字段 id 主键， balance 账户余额
	// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	// 要求 ： 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
	// 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

	// 初始化数据库连接
	db, err := sqlx.Connect("mysql", "testuser:password@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表
	createTables(db)

	// 初始化测试数据
	initializeTestData(db)

	// 执行转账事务
	err = transferMoney(db, 1, 2, 100.0)
	if err != nil {
		log.Printf("转账失败: %v", err)
	} else {
		fmt.Println("转账成功!")
	}

	// 查询并显示账户余额
	showAccountBalances(db)
}

// 创建表结构
func createTables(db *sqlx.DB) {
	// 创建 accounts 表
	accountsTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		balance DECIMAL(10,2) NOT NULL DEFAULT 0.00
	)`
	_, err := db.Exec(accountsTable)
	if err != nil {
		log.Printf("创建 accounts 表失败: %v", err)
	}

	// 创建 transactions 表
	transactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		from_account_id INT NOT NULL,
		to_account_id INT NOT NULL,
		amount DECIMAL(10,2) NOT NULL
	)`
	_, err = db.Exec(transactionsTable)
	if err != nil {
		log.Printf("创建 transactions 表失败: %v", err)
	}
}

// 初始化测试数据
func initializeTestData(db *sqlx.DB) {
	// 清空表数据
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM accounts")

	// 插入初始账户数据
	accounts := []Account{
		{ID: 1, Balance: 500.00}, // 账户A，余额500元
		{ID: 2, Balance: 300.00}, // 账户B，余额300元
	}

	for _, account := range accounts {
		_, err := db.Exec("INSERT INTO accounts (id, balance) VALUES (?, ?)", account.ID, account.Balance)
		if err != nil {
			log.Printf("插入账户数据失败: %v", err)
		}
	}

	fmt.Println("初始化测试数据完成")
	fmt.Println("==================")
}

// 转账事务函数
func transferMoney(db *sqlx.DB, fromAccountID, toAccountID int, amount float64) error {
	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	// 使用 defer 确保在函数退出时正确处理事务
	defer func() {
		if err != nil {
			// 如果有错误，回滚事务
			tx.Rollback()
			fmt.Println("事务已回滚")
		}
	}()

	// 1. 查询转出账户余额
	var fromAccount Account
	err = tx.Get(&fromAccount, "SELECT id, balance FROM accounts WHERE id = ? FOR UPDATE", fromAccountID)
	if err != nil {
		return fmt.Errorf("查询转出账户失败: %v", err)
	}

	// 2. 检查余额是否足够
	if fromAccount.Balance < amount {
		return fmt.Errorf("账户余额不足，当前余额: %.2f, 需要: %.2f", fromAccount.Balance, amount)
	}

	// 3. 查询转入账户
	var toAccount Account
	err = tx.Get(&toAccount, "SELECT id, balance FROM accounts WHERE id = ? FOR UPDATE", toAccountID)
	if err != nil {
		return fmt.Errorf("查询转入账户失败: %v", err)
	}

	// 4. 更新转出账户余额
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("更新转出账户失败: %v", err)
	}

	// 5. 更新转入账户余额
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccountID)
	if err != nil {
		return fmt.Errorf("更新转入账户失败: %v", err)
	}

	// 6. 记录交易信息
	_, err = tx.Exec("INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)",
		fromAccountID, toAccountID, amount)
	if err != nil {
		return fmt.Errorf("记录交易信息失败: %v", err)
	}

	// 7. 提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Printf("转账成功: 从账户%d向账户%d转账%.2f元\n", fromAccountID, toAccountID, amount)
	return nil
}

// 显示账户余额
func showAccountBalances(db *sqlx.DB) {
	var accounts []Account
	err := db.Select(&accounts, "SELECT id, balance FROM accounts ORDER BY id")
	if err != nil {
		log.Printf("查询账户余额失败: %v", err)
		return
	}

	fmt.Println("\n当前账户余额:")
	fmt.Println("=============")
	for _, account := range accounts {
		fmt.Printf("账户%d: %.2f元\n", account.ID, account.Balance)
	}
}
