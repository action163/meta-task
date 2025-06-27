# PersonalBlog

PersonalBlog 是一个基于 Go 语言和 Gin 框架开发的简单博客系统，支持用户注册、登录、文章发布、评论等功能，并集成了 JWT 认证和日志记录。

## 目录结构

```
PersonalBlog/
├── app.log
├── go.mod
├── go.sum
├── main.go
├── comments/
│   └── comment.go
├── global/
│   └── global.go
├── middleware/
│   └── middleware.go
├── posts/
│   └── post.go
└── users/
    └── user.go
```

## 功能简介

- 用户注册与登录（密码加密，JWT 认证）
- 文章的增删改查
- 文章评论功能
- 全局日志记录（Logrus）
- 错误恢复中间件
- 数据库操作基于 GORM

## 快速开始

### 1. 环境准备

- Go 1.18+
- MySQL 数据库

### 2. 数据库配置

请确保本地 MySQL 已创建 `personal_blog` 数据库，并在 `main.go` 的 `ConnectToDatabase` 函数中配置好用户名和密码：

```go
dbConnstr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/personal_blog?charset=utf8&parseTime=True&loc=Local"
```

### 3. 安装依赖

在项目根目录下执行：

```sh
go mod tidy
```

### 4. 数据库迁移

首次运行前请确保已创建表结构，可在 main.go 中添加如下迁移代码后运行一次：

```go
db.AutoMigrate(&users.User{}, &posts.Post{}, &comments.Comment{})
```

### 5. 启动服务

```sh
go run main.go
```

服务默认监听在 `http://localhost:8080`

## API 说明

### 用户相关

- `POST /user/register` 用户注册
- `POST /user/login` 用户登录（返回 JWT）

### 文章相关（需携带 JWT）

- `POST /post/create` 创建文章
- `GET /post/list` 获取文章列表
- `GET /post/detail/:id` 获取文章详情
- `POST /post/update/:id` 更新文章
- `DELETE /post/delete/:id` 删除文章

### 评论相关（需携带 JWT）

- `POST /comment/create` 创建评论
- `GET /comment/list/:post_id` 获取指定文章的评论列表

## 认证说明

除注册和登录外，所有文章和评论相关接口均需在 Header 中携带 `Authorization: <token>`。

## 日志

日志采用 Logrus，日志文件输出到 `app.log`。

## 贡献

欢迎提交 issue 和 PR！

---

如需更多细节，请参考各模块源码：

- main.go
- users/user.go
- posts/post.go
- comments/comment.go
- middleware/middleware.go
- global/global.go
