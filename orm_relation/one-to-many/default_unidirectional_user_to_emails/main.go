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
	Emails []Email
}

type Email struct {
	ID     uint `gorm:"primary_key"`
	Email  string
	UserID uint
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
	var user User
	user = User{
		Name: "Guobin",
		Emails: []Email{
			Email{Email: "qracle@126.com"},
			Email{Email: "guobin@gmail.com"},
		},
	}
	db.Save(&user) //级联保存

	// 从user导航到emails
	var emails []Email
	db.First(&user).Related(&emails, "Emails")
	fmt.Println("emails:", emails)
	//or
	//db.First(&user).Association("Emails").Find(&emails)
	//fmt.Println("emails:", emails)

	app.Run(iris.Addr(":8082"), iris.WithoutServerError(iris.ErrServerClosed))
}
