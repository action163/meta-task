package users

import (
	"PersonalBlog/global"
	"PersonalBlog/posts"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Response = global.Response
type Post = posts.Post
type User struct {
	gorm.Model
	Username string `gorm:"size:100" json:"username" required:"true"`
	Password string `gorm:"size:100" json:"password" required:"true"`
	Email    string `gorm:"size:100" json:"email" required:"true"`
	Posts    []Post `gorm:"foreignKey:UserID"`
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
