# 博客系统

## 运行环境

 - Go 1.25.1
 - MySQL 9.4.0
 - Git (用于版本控制)


##  依赖安装步骤 （在 task_4_bloh_system 下执行）

```
# 初始化
go mod init blog-system

# 安装目录下依赖
go mod tidy
```

## 启动方式

 - go run main.go


 ## test

 ```
 # 注册用户
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 创建文章
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzU5Mjc0MjI0fQ.jgy3so6O2mVfvjM6KlRD635OyqSLn_PmM2SpSduw9ew" \
  -d '{"title":"我的博客","content":"博客内容..."}'
 ```