package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint
	Name      string
	PostCount int
	Posts     []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	ID       uint
	Title    string
	Comments []Comment `gorm:"foreignKey:PostID"`
	UserID   uint
}

type Comment struct {
	ID      uint
	PostID  uint
	Content string
}

func main() {
	dbConnStr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// db.AutoMigrate(&Comment{})
	// db.AutoMigrate(&Post{})
	// db.AutoMigrate(&User{})
	user1 := User{
		Name:      "张三",
		PostCount: 0,
		Posts: []Post{
			{
				Title: "张三的第一篇博客",
				Comments: []Comment{
					{
						Content: "这是张三的第一篇博客的第一条评论",
					},
					{
						Content: "这是张三的第一篇博客的第二条评论",
					},
				},
			},
			{
				Title: "张三的第二篇博客",
				Comments: []Comment{
					{
						Content: "这是张三的第二篇博客的第一条评论",
					},
					{
						Content: "这是张三的第二篇博客的第二条评论",
					},
				},
			},
		},
	}
	user2 := User{
		Name:      "李四",
		PostCount: 0,
		Posts: []Post{
			{
				Title: "李四的第一篇博客",
				Comments: []Comment{
					{
						Content: "这是李四的第一篇博客的第一条评论",
					},
					{
						Content: "这是李四的第一篇博客的第二条评论",
					},
				},
			},
			{
				Title: "李四的第二篇博客",
				Comments: []Comment{
					{
						Content: "这是李四的第二篇博客的第一条评论",
					},
					{
						Content: "这是李四的第二篇博客的第二条评论",
					},
				},
			},
		},
	}

	db.Create(&user1)
	db.Create(&user2)

	posts := []Post{}
	userId := 1
	db.Model(&Post{}).Where("user_id = ?", userId).Preload("Comments").Find(&posts)
	for _, post := range posts {
		fmt.Println(post.ID, post.Title, post.Comments)
	}

}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).
		Error
}
