package comments

import (
	"PersonalBlog/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Response = global.Response
type Comment struct {
	gorm.Model
	Content string `gorm:"size:100" json:"content" required:"true"`
	UserID  uint   `json:"user_id" required:"true"`
	PostID  uint   `json:"post_id" required:"true"`
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
