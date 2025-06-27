package middleware

import (
	"PersonalBlog/global"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Response = global.Response

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		global.Log.WithFields(logrus.Fields{
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
				global.Log.WithFields(logrus.Fields{
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
