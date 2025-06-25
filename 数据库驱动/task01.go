package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	Id    uint `gorm:"primarykey"`
	Name  string
	Age   int
	Grade string
}

func main() {
	dbConnStr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//创建Students表
	db.AutoMigrate(&Student{})

	//插入一条记录
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	db.Create(&student)

	//查询年龄大于18岁的学生信息
	students := []Student{}
	db.Debug().Model(&Student{}).Where("age > ?", 18).Find(&students)
	fmt.Println("年龄大于18岁的学生信息：", students)

	//将姓名为张三的学生年级更新为四年级
	db.Debug().Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")

	//删除年龄小于15岁的学生记录
	db.Debug().Where("age < ?", 15).Delete(&Student{})
}
