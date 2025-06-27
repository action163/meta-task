package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100" json:"username" required:"true"`
	Password string `gorm:"size:100" json:"password" required:"true"`
	Email    string `gorm:"size:100" json:"email" required:"true"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}
type Post struct {
	gorm.Model
	Title    string    `gorm:"size:100" json:"title" required:"true"`
	Content  string    `gorm:"size:100" json:"content" required:"true"`
	UserID   uint      `json:"user_id"  required:"true"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}
type Comment struct {
	gorm.Model
	Content string `gorm:"size:100" json:"content" required:"true"`
	UserID  uint   `json:"user_id" required:"true"`
	PostID  uint   `json:"post_id" required:"true"`
}

type SimplePost struct {
	PostID   int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	AuthorID uint   `gorm:"column:user_id"`
	Author   string `gorm:"column:username"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未提供Token"})
			return
		}

		// 解析Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名算法")
			}
			return []byte(secret), nil
		})

		// 验证Token有效性
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效Token"})
			return
		}

		// 存储Claims到上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user", claims)
		}
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"latency": latency,
		}).Info("request details")
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				log.WithFields(logrus.Fields{
					"error":   err,
					"stack":   string(debug.Stack()),
					"request": c.Request.URL.Path,
				}).Error("panic recovered")

				// 统一错误响应
				c.JSON(http.StatusInternalServerError, Response{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
					Data:    nil,
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}

func ConfigureLogger() {
	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)

	// 创建日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// 设置JSON格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}

var log = logrus.New()

func main() {

	ConfigureLogger()
	//Connect to database
	db, err := ConnectToDatabase()
	if err != nil {
		log.Error("Failed to connect to the database:", err)
		return
	}

	r := gin.Default()
	r.Use(LoggerMiddleware())
	r.Use(RecoveryMiddleware())
	r.Use(DatabaseMiddleware(db))
	r.Use(JWTAuthMiddleware("secret_key"))
	//用户注册
	r.POST("/register", Register)

	//用户登录
	r.POST("/login", Login)

	//文章创建
	r.POST("/post/create", PostCreate)
	//获取文章列表
	r.GET("/post/list", PostList)
	//获取文章详情
	r.GET("/post/detail/:id", PostDetail)
	//更新文章
	r.POST("/post/update/:id", PostUpdate)
	//删除文章
	r.DELETE("/post/delete/:id", PostDelete)
	//创建Comment
	r.POST("/comment/create", CommentCreate)
	//获取文章评论列表
	r.GET("/comment/list/:post_id", CommentListByPost)

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

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
			Data:    nil,
		})
		return
	}
	user.Password = string(bcryptPassword)
	//创建用户
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "User registered successfully",
		Data:    user,
	})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	var dbUser User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
			Data:    nil,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
			Data:    nil,
		})
		return
	}

	//生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       dbUser.ID,
		"username": dbUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}
	//将用户信息存储到上下文中
	c.Set("user", jwt.MapClaims{
		"id":       dbUser.ID,
		"username": dbUser.Username,
	})

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data: gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":       dbUser.ID,
				"username": dbUser.Username,
				"email":    dbUser.Email,
			},
		},
	})
}

func PostCreate(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create post",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "Post created successfully",
		Data:    post,
	})
}
func PostList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var posts []SimplePost
	//只获取文章的title 和 作者信息，其中作者信息要根据UserID从User表中获取
	if err := db.Debug().Table("posts").Where("posts.deleted_at IS NULL").Select("posts.id, posts.title, posts.user_id, users.username").Joins("join users on users.id = posts.user_id").Scan(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve posts",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Posts retrieved successfully",
		Data:    posts,
	})
}
func PostDetail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var post Post
	if err := db.Debug().Preload("Comments").Where("id = ? and deleted_at IS NULL", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "Post not found",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Post retrieved successfully",
		Data:    post,
	})
}
func PostUpdate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var post Post
	if err := db.Where("id = ? and deleted_at IS NULL", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "Post not found",
			Data:    nil,
		})
		return
	}

	userIDFloat := c.MustGet("user").(jwt.MapClaims)["id"].(float64)
	if post.UserID != uint(userIDFloat) {
		c.JSON(http.StatusForbidden, Response{
			Code:    http.StatusForbidden,
			Message: "You are not allowed to update this post",
			Data:    nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update post",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Post updated successfully",
		Data:    post,
	})
}
func PostDelete(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var post Post
	if err := db.Where("id = ? and deleted_at IS NULL", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "Post not found",
			Data:    nil,
		})
		return
	}
	// 检查当前用户是否是文章的作者
	userIDFloat := c.MustGet("user").(jwt.MapClaims)["id"].(float64)
	if post.UserID != uint(userIDFloat) {
		c.JSON(http.StatusForbidden, Response{
			Code:    http.StatusForbidden,
			Message: "You are not allowed to delete this post",
			Data:    nil,
		})
		return
	}

	if err := db.Where("id = ?", c.Param("id")).Delete(&Post{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete post",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Post deleted successfully",
		Data:    nil,
	})
}

func CommentCreate(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	userIDFloat := c.MustGet("user").(jwt.MapClaims)["id"].(float64)
	comment.UserID = uint(userIDFloat)

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create comment",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "Comment created successfully",
		Data:    comment,
	})
}

func CommentListByPost(c *gin.Context) {
	postID := c.Param("post_id")
	db := c.MustGet("db").(*gorm.DB)

	var comments []Comment
	if err := db.Where("post_id = ? and deleted_at IS NULL", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve comments",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Comments retrieved successfully",
		Data:    comments,
	})
}
