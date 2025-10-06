# web3study


### golang_base homework

 - go run golang_base/task_1/main.go 梦的开始
 - go run golang_base/task_2/main.go 语言进阶
 - go run golang_base/task_3/t1_base_question1/main.go 数据库驱动-SQL语句练习
 - go run golang_base/task_3/t1_base_question2/main.go 数据库驱动-SQL语句练习
 - go run golang_base/task_3/t2_sqlx_question1/main.go 数据库驱动-Sqlx入门
 - go run golang_base/task_3/t2_sqlx_question2/main.go 数据库驱动-Sqlx入门
 - go run golang_base/task_3/t3_gorm_question1/main.go 数据库驱动-进阶gorm
 - go run golang_base/task_3/t3_gorm_question2/main.go 数据库驱动-进阶gorm
 - go run golang_base/task_3/t3_gorm_question3/main.go 数据库驱动-进阶gorm

### go base used

```
go mod github.com/web3study

go get -u gorm.io/driver/mysql
```

### solidity_base

#### task_2

```
# 需要的任务目录
cd solidity_base/task_2
# npm 初始化
npm init -y
# 安装 Hardhat（推荐）
npm install --save-dev hardhat

# 安装 OpenZeppelin 合约
npm install @openzeppelin/contracts

# 安装其他开发工具
npm install --save-dev @nomiclabs/hardhat-ethers ethers
# 依赖目录
node_modules/@openzeppelin/contracts
# 查看版本
npm openzeppelin --version
```

###### SimpleNFT
```
npm install
npx hardhat compile
```