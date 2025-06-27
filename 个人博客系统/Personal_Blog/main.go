package main

import (
	"PersonalBlog/comments"
	"PersonalBlog/global"
	"PersonalBlog/middleware"
	"PersonalBlog/posts"
	"PersonalBlog/users"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	global.ConfigureLogger()
	//Connect to database
	db, err := ConnectToDatabase()
	if err != nil {
		global.Log.Error("Failed to connect to the database:", err)
		return
	}

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.DatabaseMiddleware(db))
	//用户注册
	u := r.Group("/user")
	{
		u.POST("/register", users.Register)
		u.POST("/login", users.Login)
	}

	p := r.Group("/post", middleware.JWTAuthMiddleware("secret_key"))
	{
		//文章创建
		p.POST("/create", posts.PostCreate)
		//获取文章列表
		p.GET("/list", posts.PostList)
		//获取文章详情
		p.GET("/detail/:id", posts.PostDetail)
		//更新文章
		p.POST("/update/:id", posts.PostUpdate)
		//删除文章
		p.DELETE("/delete/:id", posts.PostDelete)
	}

	c := r.Group("/comment", middleware.JWTAuthMiddleware("secret_key"))
	{
		//创建Comment
		c.POST("/create", comments.CommentCreate)
		//获取文章评论列表
		c.GET("/list/:post_id", comments.CommentListByPost)
	}

	r.Run(":8080") // Run on port 8080

}

func ConnectToDatabase() (*gorm.DB, error) {
	dbConnstr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/personal_blog?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbConnstr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
