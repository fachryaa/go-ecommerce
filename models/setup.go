package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)
	fmt.Println(URL)

	db, err := gorm.Open(mysql.Open(URL))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}, &Cart{}, &Order{})
	if err != nil {
		panic(err)
	}

	seedProduct(db)

	DB = db
}

func seedProduct(db *gorm.DB) {
	// seed data
	if err := db.AutoMigrate(&Product{}); err == nil && db.Migrator().HasTable(&Product{}) {
		if err := db.First(&Product{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert seed data
			db.Create(SeedData)
		}
	}
}
