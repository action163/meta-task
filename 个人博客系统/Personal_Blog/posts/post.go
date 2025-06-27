package posts

import (
	"PersonalBlog/comments"
	"PersonalBlog/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Comment = comments.Comment
type Response = global.Response
type Post struct {
	gorm.Model
	Title    string    `gorm:"size:100" json:"title" required:"true"`
	Content  string    `gorm:"size:100" json:"content" required:"true"`
	UserID   uint      `json:"user_id"  required:"true"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

type SimplePost struct {
	PostID   int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	AuthorID uint   `gorm:"column:user_id"`
	Author   string `gorm:"column:username"`
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
