package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Accounts struct {
	Id      uint `gorm:"primarykey"`
	Balance float64
}

type Transactions struct {
	Id            uint `gorm:"primarykey"`
	FromAccountId uint
	ToAccountId   uint
	Amount        float64
}

func main() {
	dbConnStr := "root:p@ssw0rd@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//创建Accounts表和Transactions表 并添加两条记录
	db.AutoMigrate(&Accounts{})
	db.AutoMigrate(&Transactions{})
	db.Create(&Accounts{Balance: 300})
	db.Create(&Accounts{Balance: 50})

	accountA := Accounts{}
	accountB := Accounts{}
	db.Debug().Model(&Accounts{}).Where("id = ?", 1).Find(&accountA)
	db.Debug().Model(&Accounts{}).Where("id = ?", 2).Find(&accountB)

	fmt.Println("accountA balance is:", accountA.Balance, accountA.Id)
	fmt.Println("accountB balance is:", accountB.Balance, accountB.Id)
	db.Transaction(func(tx *gorm.DB) error {
		if accountA.Balance < 100 {
			return errors.New("余额不足")
		}
		if err := tx.Debug().Save(&Transactions{FromAccountId: accountA.Id, ToAccountId: accountB.Id, Amount: 100}).Error; err != nil {
			return err
		}

		accountA.Balance -= 100
		if err := tx.Debug().Save(&accountA).Error; err != nil {
			return err
		}
		accountB.Balance += 100
		if err := tx.Debug().Save(&accountB).Error; err != nil {
			return err
		}

		return nil
	})

}
