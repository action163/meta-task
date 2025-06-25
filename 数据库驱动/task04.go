package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID     int     `db:"id" json:"id"`
	Title  string  `db:"title" json:"title"`
	Author string  `db:"author" json:"author"`
	Price  float64 `db:"price" json:"price"`
}

func main() {
	dbConnStr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	books := []Book{}
	db.Select(&books, "Select * from books where price > ?", 50)

	for _, book := range books {
		fmt.Println(book.ID, book.Title, book.Author, book.Price)
	}
}
