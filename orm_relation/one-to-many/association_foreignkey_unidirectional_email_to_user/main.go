package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

type User struct {
	ID   uint `gorm:"primary_key"`
	Uid  uint
	Name string
}

type Email struct {
	ID      uint `gorm:"primary_key"`
	Email   string
	User    User `gorm:"foreignkey:UserRef;association_foreignkey:Uid"`
	UserRef uint
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
	db.DropTableIfExists(&Email{}) // drop table

	db.AutoMigrate(&User{}) // create table: // AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
	db.AutoMigrate(&Email{})

	// 造数据
	var email Email
	email = Email{
		Email: "qracle@126.com",
		User: User{
			Name: "Guobin",
			Uid:  1,
		},
	}
	db.Save(&email) //级联保存

	var user User

	// 从email导航到user
	db.First(&email).Related(&user, "User")
	fmt.Println("user:", user)
	//or
	//db.First(&email).Association("User").Find(&user)
	//fmt.Println("user:", user)

	app.Run(iris.Addr(":8082"), iris.WithoutServerError(iris.ErrServerClosed))
}
