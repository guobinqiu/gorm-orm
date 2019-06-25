package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type User struct {
	ID     uint `gorm:"primary_key"`
	Name   string
	Orders []Order `gorm:"many2many:user_orders"`
}

type Order struct {
	ID      uint `gorm:"primary_key"`
	OrderNo string
	Users []User `gorm:"many2many:user_orders"`
}

func main() {
	app := iris.Default()
	db, err := gorm.Open("mysql", "root:111111@tcp(localhost:3306)/go_hello?parseTime=true")
	db.LogMode(true) // show SQL logger
	if err != nil {
		app.Logger().Fatalf("connect to mysql failed")
		return
	}
	iris.RegisterOnInterrupt(func() {
		defer db.Close()
	})

	db.DropTableIfExists(&User{})  // drop table
	db.DropTableIfExists(&Order{}) // drop table

	db.AutoMigrate(&User{}) // create table: // AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
	db.AutoMigrate(&Order{})

	var user User
	var order Order

	// 造数据
	user = User{
		Name: "Guobin",
		Orders: []Order{
			Order{OrderNo: "123"},
			Order{OrderNo: "456"},
		},
	}
	db.Save(&user) //级联保存

	// 从user导航到orders
	var orders []Order
	db.First(&user).Related(&orders, "Orders")
	fmt.Println("orders:", orders)

	// 从order导航到users
	var users []User
	db.First(&order).Related(&users, "Users")
	fmt.Println("users:", users)

	app.Run(iris.Addr(":8082"), iris.WithoutServerError(iris.ErrServerClosed))
}
