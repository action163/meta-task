package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func main() {
	dbConnStr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	employees := []Employee{}
	db.Select(&employees, "Select * from employees where department =?", "技术部")

	for _, emp := range employees {
		fmt.Println(emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	expenseEmp := Employee{}
	db.Select(&expenseEmp, "Select * from employees where salary = (select max(Salary) from employees)")
	fmt.Println("工资最高的员工信息："，expenseEmp)
}
